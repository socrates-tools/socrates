package surface

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rui-ktei/socrates/internal/assets"
	"github.com/rui-ktei/socrates/internal/document"
	"github.com/rui-ktei/socrates/internal/render"
)

const (
	TokenHeader = "x-agora-token"

	IdlePeriod   = 45 * time.Minute
	PollDeadline = 6 * time.Hour

	assetPrefix = "/assets/"
	bodyCeiling = 1 << 20
)

type Server struct {
	exchange *Exchange
	listener net.Listener
	serving  *http.Server
	token    string
	address  string
	idle     time.Duration
	stopping chan struct{}
	stopped  chan struct{}
}

func Serve(read document.Document, address, token string, idle time.Duration) (*Server, error) {
	if err := CheckLoopback(address); err != nil {
		return nil, err
	}

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf(bindingTheAddress, address, err)
	}

	server := &Server{
		exchange: NewExchange(read, time.Now()),
		listener: listener,
		token:    token,
		address:  listener.Addr().String(),
		idle:     idle,
		stopping: make(chan struct{}),
		stopped:  make(chan struct{}),
	}
	server.serving = &http.Server{Handler: server.routes(), ReadHeaderTimeout: 10 * time.Second}

	go server.hold()
	go server.watchIdle()
	return server, nil
}

func (s *Server) Address() string {
	return s.address
}

func (s *Server) Exchange() *Exchange {
	return s.exchange
}

func (s *Server) hold() {
	defer close(s.stopped)
	if err := s.serving.Serve(s.listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return
	}
}

func (s *Server) watchIdle() {
	if s.idle <= 0 {
		return
	}

	ticker := time.NewTicker(s.idle / 4)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopping:
			return
		case now := <-ticker.C:
			if s.exchange.IdleSince(now, s.idle) {
				s.Stop()
				return
			}
		}
	}
}

func (s *Server) Stop() {
	select {
	case <-s.stopping:
		return
	default:
		close(s.stopping)
	}

	shutting, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = s.serving.Shutdown(shutting)
}

func (s *Server) Wait() {
	<-s.stopped
}

func (s *Server) Stopping() <-chan struct{} {
	return s.stopping
}

func (s *Server) routes() http.Handler {
	routed := http.NewServeMux()
	routed.HandleFunc("GET /{$}", s.page)
	routed.HandleFunc("GET "+assetPrefix, s.asset)
	routed.HandleFunc("GET /live", s.live)
	routed.HandleFunc("POST /answer", s.guarded(s.answer))
	routed.HandleFunc("GET /answer", s.guarded(s.deliver))
	routed.HandleFunc("POST /document", s.guarded(s.revise))
	routed.HandleFunc("POST /close", s.guarded(s.close))
	routed.HandleFunc("POST /stop", s.guarded(s.stop))
	return s.checkingHost(routed)
}

func (s *Server) checkingHost(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if err := hostIsTheSurface(request.Host, s.address); err != nil {
			http.Error(writer, err.Error(), http.StatusMisdirectedRequest)
			return
		}
		s.exchange.Touch(time.Now())
		next.ServeHTTP(writer, request)
	})
}

func (s *Server) guarded(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Header.Get(TokenHeader) != s.token {
			http.Error(writer, answerCarriesNoToken, http.StatusForbidden)
			return
		}
		next(writer, request)
	}
}

func (s *Server) page(writer http.ResponseWriter, request *http.Request) {
	read, revision := s.exchange.Document()

	written, err := render.Page(read, render.Frame{Token: s.token, Revision: revision, Over: s.exchange.Over()})
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("content-type", "text/html; charset=utf-8")
	writer.Header().Set("cache-control", "no-store")
	_, _ = io.WriteString(writer, written)
}

func (s *Server) asset(writer http.ResponseWriter, request *http.Request) {
	http.StripPrefix(assetPrefix, http.FileServerFS(assets.Files())).ServeHTTP(writer, request)
}

func (s *Server) live(writer http.ResponseWriter, request *http.Request) {
	flusher, streaming := writer.(http.Flusher)
	if !streaming {
		http.Error(writer, cannotStream, http.StatusInternalServerError)
		return
	}

	s.exchange.Reading(1, time.Now())
	defer s.exchange.Reading(-1, time.Now())

	writer.Header().Set("content-type", "text/event-stream")
	writer.Header().Set("cache-control", "no-store")
	writer.WriteHeader(http.StatusOK)

	waiting := s.exchange.Changed()
	flusher.Flush()

	for {
		if s.exchange.Over() {
			_, _ = io.WriteString(writer, "event: closed\ndata: yes\n\n")
			flusher.Flush()
			return
		}

		select {
		case <-request.Context().Done():
			return
		case <-s.stopping:
			return
		case <-waiting:
		}
		waiting = s.exchange.Changed()

		if s.exchange.Over() {
			continue
		}
		_, revision := s.exchange.Document()
		_, _ = io.WriteString(writer, "event: revision\ndata: "+strconv.Itoa(revision)+"\n\n")
		flusher.Flush()
	}
}

func (s *Server) answer(writer http.ResponseWriter, request *http.Request) {
	var sent Sent
	if err := readBody(request, &sent); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := s.exchange.Receive(sent); err != nil {
		http.Error(writer, err.Error(), http.StatusConflict)
		return
	}
	writer.WriteHeader(http.StatusNoContent)
}

func (s *Server) deliver(writer http.ResponseWriter, request *http.Request) {
	within := PollDeadline
	if stated := request.URL.Query().Get("within"); stated != "" {
		parsed, err := time.ParseDuration(stated)
		if err != nil {
			http.Error(writer, fmt.Sprintf(waitIsNotADuration, stated), http.StatusBadRequest)
			return
		}
		within = parsed
	}

	answered, arrived := s.exchange.Await(request.Context(), within)
	writer.Header().Set("content-type", "application/json")
	if !arrived {
		writer.WriteHeader(http.StatusNoContent)
		return
	}
	_ = json.NewEncoder(writer).Encode(answered)
}

func (s *Server) revise(writer http.ResponseWriter, request *http.Request) {
	held, err := io.ReadAll(io.LimitReader(request.Body, bodyCeiling))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	read, err := document.Read(held)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	if err := WriteDocument(held); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	s.exchange.Revise(read)
	writer.WriteHeader(http.StatusNoContent)
}

func (s *Server) close(writer http.ResponseWriter, request *http.Request) {
	s.exchange.End()
	writer.WriteHeader(http.StatusNoContent)
}

func (s *Server) stop(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusNoContent)
	go s.Stop()
}

func readBody(request *http.Request, into any) error {
	held, err := io.ReadAll(io.LimitReader(request.Body, bodyCeiling))
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(strings.NewReader(string(held)))
	decoder.DisallowUnknownFields()
	return decoder.Decode(into)
}
