package server

import (
	"net/http"
	"time"
)

const HTTP_PORT = "8080"

// Структура HTTP сервера приложения
type Server struct {
	httpServer *http.Server
}

// Функция запуска HTTP сервера
func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return s.httpServer.ListenAndServe()
}
