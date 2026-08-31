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
	"sync"
	"time"

	"github.com/ktei/socrates/internal/assets"
	"github.com/ktei/socrates/internal/document"
	"github.com/ktei/socrates/internal/refusal"
	"github.com/ktei/socrates/internal/render"
)

const (
	TokenHeader   = "x-agora-token"
	SurfaceHeader = "x-agora-surface"

	IdlePeriod   = 45 * time.Minute
	PollDeadline = 6 * time.Hour

	SessionPrefix = "/sessions/"

	BodyCeiling = 1 << 20

	assetPrefix = "/assets/"
)

type Opened struct {
	Session string `json:"session"`
	Token   string `json:"token"`
}

type Server struct {
	listener net.Listener
	serving  *http.Server
	control  string
	address  string
	idle     time.Duration
	stopping chan struct{}
	stopped  chan struct{}
	stopOnce sync.Once

	sessions *Sessions
}

func Serve(address string, idle time.Duration) (*Server, error) {
	if err := CheckLoopback(address); err != nil {
		return nil, err
	}

	control, err := MintToken()
	if err != nil {
		return nil, err
	}

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf(bindingTheAddress, address, err)
	}

	server := &Server{
		listener: listener,
		control:  control,
		address:  listener.Addr().String(),
		idle:     idle,
		stopping: make(chan struct{}),
		stopped:  make(chan struct{}),
		sessions: NewSessions(time.Now()),
	}
	server.serving = &http.Server{Handler: server.routes(), ReadHeaderTimeout: 10 * time.Second}

	go server.hold()
	go server.watchIdle()
	return server, nil
}

func (s *Server) Address() string {
	return s.address
}

func (s *Server) Control() string {
	return s.control
}

func (s *Server) Open(read document.Document, raw []byte) (Opened, error) {
	name, err := MintSession()
	if err != nil {
		return Opened{}, err
	}
	token, err := MintToken()
	if err != nil {
		return Opened{}, err
	}

	if err := WritePage(Page{Session: name, Address: s.address, Token: token}, raw); err != nil {
		return Opened{}, err
	}

	s.sessions.Open(name, token, read, time.Now())
	return Opened{Session: name, Token: token}, nil
}

func (s *Server) Revise(session string, read document.Document, raw []byte) error {
	one, standing := s.sessions.Exchange(session)
	if !standing {
		return refusal.New(surfaceHoldsNoSuchSession, session)
	}
	if err := WriteDocument(session, raw); err != nil {
		return err
	}
	one.Revise(read)
	return nil
}

func (s *Server) Exchange(session string) (*Exchange, bool) {
	return s.sessions.Exchange(session)
}

func (s *Server) Token(session string) (string, bool) {
	return s.sessions.Token(session)
}

func (s *Server) Holds() []string {
	return s.sessions.Names()
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
			collected, empty := s.sessions.Collect(now, s.idle)
			for _, name := range collected {
				_ = ForgetPage(name)
			}
			if empty {
				s.Stop()
				return
			}
		}
	}
}

func (s *Server) Stop() {
	s.stopOnce.Do(func() {
		close(s.stopping)

		shutting, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = s.serving.Shutdown(shutting)
	})
}

func (s *Server) Wait() {
	<-s.stopped
}

func (s *Server) Stopping() <-chan struct{} {
	return s.stopping
}

func (s *Server) routes() http.Handler {
	routed := http.NewServeMux()
	routed.HandleFunc("GET /{$}", s.controlled(s.reachable))
	routed.HandleFunc("GET "+assetPrefix, s.asset)
	routed.HandleFunc("POST "+SessionPrefix+"{$}", s.controlled(s.opening))
	routed.HandleFunc("POST /stop", s.controlled(s.stop))
	routed.HandleFunc("GET "+SessionPrefix+"{session}", s.page)
	routed.HandleFunc("GET "+SessionPrefix+"{session}/{$}", s.page)
	routed.HandleFunc("GET "+SessionPrefix+"{session}/live", s.live)
	routed.HandleFunc("POST "+SessionPrefix+"{session}/answer", s.guarded(s.answer))
	routed.HandleFunc("GET "+SessionPrefix+"{session}/answer", s.guarded(s.deliver))
	routed.HandleFunc("POST "+SessionPrefix+"{session}/document", s.guarded(s.revise))
	routed.HandleFunc("POST "+SessionPrefix+"{session}/close", s.guarded(s.close))
	return s.checkingHost(routed)
}

func (s *Server) checkingHost(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if err := hostIsTheSurface(request.Host, s.address); err != nil {
			http.Error(writer, err.Error(), http.StatusMisdirectedRequest)
			return
		}
		next.ServeHTTP(writer, request)
	})
}

func (s *Server) routed(request *http.Request) (*held, bool) {
	return s.sessions.find(request.PathValue("session"))
}

func (s *Server) missing(writer http.ResponseWriter, request *http.Request) {
	http.Error(writer, fmt.Sprintf(surfaceHoldsNoSuchSession, request.PathValue("session")), http.StatusNotFound)
}

func (s *Server) controlled(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Header.Get(TokenHeader) != s.control {
			http.Error(writer, requestCarriesNoControl, http.StatusForbidden)
			return
		}
		next(writer, request)
	}
}

func (s *Server) guarded(next func(http.ResponseWriter, *http.Request, *held)) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		one, standing := s.routed(request)
		if !standing {
			s.missing(writer, request)
			return
		}
		if carried := request.Header.Get(TokenHeader); carried != one.token && carried != s.control {
			http.Error(writer, answerCarriesNoToken, http.StatusForbidden)
			return
		}
		one.exchange.Touch(time.Now())
		next(writer, request, one)
	}
}

func (s *Server) reachable(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set(SurfaceHeader, s.address)
	writer.WriteHeader(http.StatusNoContent)
}

func (s *Server) opening(writer http.ResponseWriter, request *http.Request) {
	raw, within := bodyWithin(writer, request)
	if !within {
		return
	}

	read, err := document.Read(raw)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	opened, err := s.Open(read, raw)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(writer).Encode(opened)
}

func (s *Server) page(writer http.ResponseWriter, request *http.Request) {
	one, standing := s.routed(request)
	if !standing {
		s.missing(writer, request)
		return
	}
	one.exchange.Touch(time.Now())
	read, revision := one.exchange.Document()

	written, err := render.Page(read, render.Frame{
		Base:     SessionPrefix + request.PathValue("session"),
		Token:    one.token,
		Revision: revision,
		Over:     one.exchange.Over(),
	})
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
	one, standing := s.routed(request)
	if !standing {
		s.missing(writer, request)
		return
	}
	flusher, streaming := writer.(http.Flusher)
	if !streaming {
		http.Error(writer, cannotStream, http.StatusInternalServerError)
		return
	}

	one.exchange.Reading(1, time.Now())
	defer one.exchange.Reading(-1, time.Now())

	writer.Header().Set("content-type", "text/event-stream")
	writer.Header().Set("cache-control", "no-store")
	writer.WriteHeader(http.StatusOK)

	waiting := one.exchange.Changed()
	flusher.Flush()

	for {
		if one.exchange.Over() {
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
		waiting = one.exchange.Changed()

		if one.exchange.Over() {
			continue
		}
		_, revision := one.exchange.Document()
		_, _ = io.WriteString(writer, "event: revision\ndata: "+strconv.Itoa(revision)+"\n\n")
		flusher.Flush()
	}
}

func (s *Server) answer(writer http.ResponseWriter, request *http.Request, one *held) {
	var sent Sent
	if !readSent(writer, request, &sent) {
		return
	}

	if _, err := one.exchange.Receive(sent); err != nil {
		http.Error(writer, err.Error(), http.StatusConflict)
		return
	}
	writer.WriteHeader(http.StatusNoContent)
}

func (s *Server) deliver(writer http.ResponseWriter, request *http.Request, one *held) {
	within := PollDeadline
	if stated := request.URL.Query().Get("within"); stated != "" {
		parsed, err := time.ParseDuration(stated)
		if err != nil {
			http.Error(writer, fmt.Sprintf(waitIsNotADuration, stated), http.StatusBadRequest)
			return
		}
		within = parsed
	}

	one.exchange.Reading(1, time.Now())
	defer one.exchange.Reading(-1, time.Now())

	answered, arrived := one.exchange.Await(request.Context(), within)
	writer.Header().Set("content-type", "application/json")
	if !arrived {
		writer.WriteHeader(http.StatusNoContent)
		return
	}
	_ = json.NewEncoder(writer).Encode(answered)
}

func (s *Server) revise(writer http.ResponseWriter, request *http.Request, one *held) {
	raw, within := bodyWithin(writer, request)
	if !within {
		return
	}

	read, err := document.Read(raw)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	if err := WriteDocument(request.PathValue("session"), raw); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	one.exchange.Revise(read)
	writer.WriteHeader(http.StatusNoContent)
}

func (s *Server) close(writer http.ResponseWriter, request *http.Request, one *held) {
	one.exchange.End()
	writer.WriteHeader(http.StatusNoContent)
}

func (s *Server) stop(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusNoContent)
	go s.Stop()
}

func bodyWithin(writer http.ResponseWriter, request *http.Request) ([]byte, bool) {
	raw, err := io.ReadAll(http.MaxBytesReader(writer, request.Body, BodyCeiling))
	if err == nil {
		return raw, true
	}

	var overTheCeiling *http.MaxBytesError
	if errors.As(err, &overTheCeiling) {
		http.Error(writer, fmt.Sprintf(bodyIsOverTheCeiling, BodyCeiling), http.StatusRequestEntityTooLarge)
		return nil, false
	}
	http.Error(writer, err.Error(), http.StatusBadRequest)
	return nil, false
}

func readSent(writer http.ResponseWriter, request *http.Request, into any) bool {
	raw, within := bodyWithin(writer, request)
	if !within {
		return false
	}

	decoder := json.NewDecoder(strings.NewReader(string(raw)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(into); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return false
	}
	return true
}
