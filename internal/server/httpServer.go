package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"pets/internal/config"
	"pets/internal/server/handlers"
	"pets/internal/service"
	"pets/pkg/logger"
)

// HttpServer is an app server layer struct
type HttpServer struct {
	Router   *chi.Mux
	handlers *handlers.Handlers
	conf     *config.Http
}

// NewServer is used to get new HttpServer instance
func NewServer(conf *config.Http, srv service.IService) *HttpServer {
	if conf == nil {
		logger.Log().WithField("layer", "Server").Fatalf("config is nil")
	}

	s := &HttpServer{}

	s.conf = conf

	s.Router = chi.NewRouter()
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)

	s.handlers = handlers.NewHandlers(srv)
	s.registerRoutes()

	logger.Log().WithField("layer", "Server").Infof("server created")

	return s
}

// ListenAndServ is used to run server
func (s *HttpServer) ListenAndServ() {
	logger.Log().WithField("layer", "Server").Infof("starting server at %v", s.conf.TCP)
	if err := http.ListenAndServe(s.conf.TCP, s.Router); err != nil {
		logger.Log().WithField("layer", "Server").Fatalf("error listen and serv :%v", err.Error())
	}
}

// registerRoutes is used to register routs in router
func (s *HttpServer) registerRoutes() {
	s.Router.Route("/api/v1", func(r chi.Router) {
		r.Get("/pet", s.handlers.GetPets())
		r.Post("/pet", s.handlers.CreatePet())
		r.Put("/pet", s.handlers.UpdatePet())
		r.Delete("/pet", s.handlers.DeletePet())
	})
}
