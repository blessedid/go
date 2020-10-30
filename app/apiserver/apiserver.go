package apiserver

import (
	"io"
	"net/http"
	"vk-go/app/store"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// APIServer ...
type APIServer struct {
	config *Config
	logger *zap.Logger
	router *mux.Router
	store  *store.Store
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

	if err := s.configureStore(); err != nil {
		return err
	}

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

func (s *APIServer) configureStore() error {
	st := store.New(s.config.Store)
	if err := st.Open(); err != nil {
		return err
	}

	s.store = st

	return nil
}

func (s *APIServer) handleCalbackAPI() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "ok")
	}
}
