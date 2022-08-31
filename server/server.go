package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

type server struct {
	router chi.Router
	client hatchwayClientWithCache
}

// NewServer returns a pointer to a new server
func NewServer() *server {
	newServer := server{
		router: chi.NewRouter(),
		client: hatchwayClientWithCache{
			url:          "https://api.hatchways.io/assessment/blog",
			cache:        map[string]cachedData{},
			cacheTimeout: 5 * time.Minute,
		},
	}

	newServer.router.Mount("/api", router(newServer.client))

	return &newServer
}

// ServeHTTP implements the http.Handler interface.
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
