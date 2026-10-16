package surface

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/socrates-tools/socrates/agora/document"
	"github.com/socrates-tools/socrates/agora/render"
	"github.com/socrates-tools/socrates/internal/refusal"
)

const (
	BaseHeader = "x-agora-base"

	IllustrationCeiling = 16 << 20

	admittedHosts = "https://cdnjs.cloudflare.com https://cdn.jsdelivr.net https://unpkg.com"
)

var htmlExtensions = []string{".html", ".htm"}

var mediaTypes = map[string]string{
	".html":  "text/html; charset=utf-8",
	".htm":   "text/html; charset=utf-8",
	".css":   "text/css; charset=utf-8",
	".js":    "text/javascript; charset=utf-8",
	".json":  "application/json",
	".svg":   "image/svg+xml",
	".png":   "image/png",
	".jpg":   "image/jpeg",
	".jpeg":  "image/jpeg",
	".webp":  "image/webp",
	".gif":   "image/gif",
	".woff":  "font/woff",
	".woff2": "font/woff2",
}

type drawing struct {
	position int
	decision string
	source   string
	file     string
}

func EncodeBase(base string) string {
	return url.PathEscape(base)
}

func decodeBase(header string) string {
	decoded, err := url.PathUnescape(header)
	if err != nil {
		return header
	}
	return decoded
}

func CheckIllustrations(read document.Document, base string) error {
	_, err := drawingsOf(read, base)
	return err
}

func drawingsOf(read document.Document, base string) ([]drawing, error) {
	var drawings []drawing

	for position, decision := range read.Decisions {
		path := decision.Illustration.Path
		if path == "" {
			continue
		}

		source, err := whereItIs(decision.ID, path, base)
		if err != nil {
			return nil, err
		}
		if !slices.ContainsFunc(htmlExtensions, func(held string) bool { return strings.EqualFold(held, filepath.Ext(source)) }) {
			return nil, refusal.New(illustrationIsNotHTML, decision.ID, path)
		}

		found, err := os.Stat(source)
		if errors.Is(err, os.ErrNotExist) {
			return nil, refusal.New(illustrationDoesNotExist, decision.ID, path)
		}
		if err != nil {
			return nil, refusal.New(illustrationCannotBeRead, decision.ID, path, err)
		}
		if !found.Mode().IsRegular() {
			return nil, refusal.New(illustrationIsNotAFile, decision.ID, path)
		}
		if found.Size() > IllustrationCeiling {
			return nil, refusal.New(illustrationIsTooLarge, decision.ID, path, found.Size(), IllustrationCeiling)
		}

		drawings = append(drawings, drawing{
			position: position,
			decision: decision.ID,
			source:   source,
			file:     render.BaseName(path),
		})
	}
	return drawings, nil
}

func whereItIs(decision, path, base string) (string, error) {
	if filepath.IsAbs(path) {
		return filepath.Clean(path), nil
	}
	if base == "" {
		return "", refusal.New(illustrationHasNoFolder, decision, path)
	}
	return filepath.Join(base, filepath.FromSlash(path)), nil
}

func copyIllustrations(session string, read document.Document, base string) (string, error) {
	drawings, err := drawingsOf(read, base)
	if err != nil || len(drawings) == 0 {
		return "", err
	}

	copied, err := MintSession()
	if err != nil {
		return "", err
	}
	within, err := IllustrationsPath(session, copied)
	if err != nil {
		return "", err
	}

	for _, one := range drawings {
		contents, err := os.ReadFile(one.source)
		if err != nil {
			_ = forgetCopy(session, copied)
			return "", refusal.New(illustrationCannotBeRead, one.decision, one.source, err)
		}

		folder := filepath.Join(within, strconv.Itoa(one.position))
		if err := os.MkdirAll(folder, 0o700); err != nil {
			_ = forgetCopy(session, copied)
			return "", fmt.Errorf(copyingTheIllustration, err)
		}
		if err := os.WriteFile(filepath.Join(folder, one.file), contents, 0o600); err != nil {
			_ = forgetCopy(session, copied)
			return "", fmt.Errorf(copyingTheIllustration, err)
		}
	}
	return copied, nil
}

func forgetCopy(session, copied string) error {
	if copied == "" {
		return nil
	}
	within, err := IllustrationsPath(session, copied)
	if err != nil {
		return err
	}
	return os.RemoveAll(within)
}

func illustrationPolicy(origin string) string {
	own := "'self' " + origin
	return strings.Join([]string{
		"sandbox allow-scripts",
		"default-src 'none'",
		"script-src " + own + " 'unsafe-inline' " + admittedHosts,
		"style-src " + own + " 'unsafe-inline' " + admittedHosts,
		"font-src " + own + " " + admittedHosts,
		"img-src " + own + " data:",
		"media-src " + own + " data:",
		"connect-src 'none'",
		"frame-src 'none'",
		"form-action 'none'",
		"base-uri 'none'",
	}, "; ")
}

func (s *Server) illustration(writer http.ResponseWriter, request *http.Request) {
	header := writer.Header()
	header.Set("content-security-policy", illustrationPolicy("http://"+request.Host))
	header.Set("x-content-type-options", "nosniff")
	header.Set("cache-control", "no-store")

	session := request.PathValue("session")
	one, standing := s.routed(request)
	if !standing {
		s.missing(writer, request)
		return
	}
	one.exchange.Touch(time.Now())

	read, copied := one.exchange.Illustrated()
	position := -1
	for at, decision := range read.Decisions {
		if decision.ID == request.PathValue("decision") && decision.Illustration.Path != "" {
			position = at
		}
	}
	file := request.PathValue("file")
	media, served := mediaTypes[strings.ToLower(filepath.Ext(file))]
	if position < 0 || copied == "" || !served || !filepath.IsLocal(file) {
		http.Error(writer, illustrationIsNotHere, http.StatusNotFound)
		return
	}

	within, err := IllustrationsPath(session, copied)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	contents, err := os.ReadFile(filepath.Join(within, strconv.Itoa(position), filepath.FromSlash(file)))
	if err != nil {
		http.Error(writer, illustrationIsNotHere, http.StatusNotFound)
		return
	}

	header.Set("content-type", media)
	_, _ = writer.Write(contents)
}
