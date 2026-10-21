//go:build agoradev

package surface

import (
	"io"
	"net/http"

	"github.com/socrates-tools/socrates/agora/render"
)

const CatalogRoute = "/catalog"

func (s *Server) mountTheCatalog(routed *http.ServeMux) {
	routed.HandleFunc("GET "+CatalogRoute, s.catalog)
}

func (s *Server) catalog(writer http.ResponseWriter, _ *http.Request) {
	written, err := render.Catalog()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("content-type", "text/html; charset=utf-8")
	writer.Header().Set("cache-control", "no-store")
	_, _ = io.WriteString(writer, written)
}
