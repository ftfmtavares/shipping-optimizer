package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/ftfmtavares/shipping-optimizer/internal/instrumentation"
	"github.com/gorilla/mux"
)

type HTTPServer struct {
	server *http.Server
	router *mux.Router
	logger instrumentation.Logger
	env    string
}

type HTTPServerConfig struct {
	Address      string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	Logger       instrumentation.Logger
	Env          string
}

func NewHTTPServer(sc HTTPServerConfig) HTTPServer {
	router := mux.NewRouter()
	return HTTPServer{
		server: &http.Server{
			Addr:         sc.Address + ":" + strconv.Itoa(sc.Port),
			Handler:      router,
			ReadTimeout:  sc.ReadTimeout,
			WriteTimeout: sc.WriteTimeout,
			IdleTimeout:  sc.IdleTimeout,
		},
		router: router,
		logger: sc.Logger,
		env:    sc.Env,
	}
}

func (s *HTTPServer) WithHealthCheck() {
	healthHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "healthy",
			"time":   time.Now().Format(time.DateTime),
		})
		s.logger.Info("Checking server health")
	}

	s.router.HandleFunc("/health", healthHandler).Methods("GET")
}

func (s *HTTPServer) WithRoute(path string, handler http.HandlerFunc, methods ...string) {
	handler = s.wrapPanicRecovery(handler)
	handler = s.wrapJsonContentType(handler)
	handler = s.wrapLogging(handler)

	s.router.HandleFunc(path, handler).Methods(methods...)
}

func (s *HTTPServer) StartHTTPServerAsync() {
	s.logger.Info("starting api on " + s.server.Addr)

	go func() {
		err := s.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			s.logger.Error("Failed to start server")
		}
	}()
}

func (s *HTTPServer) WithShutdownGracefully() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	s.logger.Info("Stopping server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := s.server.Shutdown(ctx)
	if err != nil {
		s.logger.Error("Server shutdown failed")
	}

	s.logger.Info("Server stopped gracefully")
}

func (s *HTTPServer) wrapLogging(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info(r.RemoteAddr + r.Method + r.URL.String())

		handler.ServeHTTP(w, r)
	})
}

func (s *HTTPServer) wrapJsonContentType(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		handler.ServeHTTP(w, r)
	})
}

func (s *HTTPServer) wrapPanicRecovery(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				s.logger.Error(fmt.Sprintf("panic recovered: %v", err))
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		handler.ServeHTTP(w, r)
	})
}
