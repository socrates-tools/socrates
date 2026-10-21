//go:build !agoradev

package surface

import "net/http"

func (s *Server) mountTheCatalog(*http.ServeMux) {}
