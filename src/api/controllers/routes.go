package controllers

import (
	"github.com/sureshk/birthday-server/src/api/logger"
	"github.com/sureshk/birthday-server/src/api/middlewares"
)

func (s *Server) initializeRoutes() {
	s.Router.Use(logger.Logger)
	s.Router.Use(middlewares.AddContentTypeMiddleware)
	// Home Route
	s.Router.HandleFunc("/", s.Home).Methods("GET")

	//Users Routes
	s.Router.HandleFunc("/users", s.GetUsers).Methods("GET")
	s.Router.HandleFunc("/users", s.CreateUser).Methods("POST")
	s.Router.HandleFunc("/users/{id}", s.GetUser).Methods("GET")
	s.Router.HandleFunc("/users/{id}", s.UpdateUser).Methods("PUT")
	s.Router.HandleFunc("/users/{id}", s.DeleteUser).Methods("DELETE")
}
