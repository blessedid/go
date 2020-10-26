package apiserver

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// APIServer ...
type APIServer struct {
	config *Config
	logger *zap.Logger
	router *mux.Router
}

// New ...
func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		logger: logger(),
		router: mux.NewRouter(),
	}
}

// Start ...
func (s *APIServer) Start() error {

	s.configureRouter()

	s.logger.Info("Server start")

	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func logger() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
	}

	return logger
}

func (s *APIServer) configureRouter() {
	s.router.HandleFunc("/calbackAPI", s.handleCalbackAPI())
}

func (s *APIServer) handleCalbackAPI() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "ok")
	}
}
