package server

import (
	"2025/internal/handlers"
	"2025/internal/service"
	"2025/internal/storage"

	"github.com/gin-gonic/gin"
)

// Server — структура HTTP-сервера.
// Хранит роутер, адрес запуска, хранилище данных и канал задач для worker pool.
type Server struct {
	router  *gin.Engine
	address string
	storage *storage.Storage
	tasks   chan service.Task
}

// NewServer — конструктор сервера.
func NewServer(addr string, strg *storage.Storage, tsks chan service.Task) *Server {
	r := gin.Default()
	return &Server{
		router:  r,
		address: addr,
		storage: strg,
		tasks:   tsks,
	}
}

// registerRoutes — регистрирует все HTTP-маршруты сервиса.
func (s *Server) registerRoutes() {
	s.router.POST("/check-links", handlers.CheckURLs(s.storage, s.tasks))
	s.router.POST("/make-pdf", handlers.MakePDF(s.storage))
}

// Start — запускает HTTP-сервер на указанном адресе.
func (s *Server) Start() error {
	s.registerRoutes()
	return s.router.Run(s.address)
}
