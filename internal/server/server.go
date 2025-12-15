package server

import (
	"2025/internal/handlers"
	"2025/internal/service"
	"2025/internal/storage"
	"sync"

	"github.com/gin-gonic/gin"
)

// Server is the HTTP server structure.
// It holds the router, startup address, data storage, and the task channel for the worker pool.
type Server struct {
	router  *gin.Engine
	address string
	storage *storage.Storage
	tasks   chan service.Task
	wg      *sync.WaitGroup
}

// NewServer is the server constructor.
func NewServer(
	addr string,
	strg *storage.Storage,
	tsks chan service.Task,
	wg *sync.WaitGroup,
) *Server {
	r := gin.Default()
	return &Server{
		router:  r,
		address: addr,
		storage: strg,
		tasks:   tsks,
		wg:      wg,
	}
}

// registerRoutes registers all HTTP routes of the service.
func (s *Server) registerRoutes() {
	s.router.POST(
		"/check-links",
		handlers.CheckURLs(s.storage, s.tasks, s.wg),
	)
	s.router.POST(
		"/make-pdf",
		handlers.MakePDF(s.storage),
	)
}

// Start starts the HTTP server on the specified address.
func (s *Server) Start() error {
	s.registerRoutes()
	return s.router.Run(s.address)
}
