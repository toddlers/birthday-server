package controllers

import "github.com/sureshk/birthday-server/src/api/middlewares"

func (s *Server) initializeRoutes() {
	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	//Users Routes
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/user", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/user/{username}", middlewares.SetMiddlewareJSON(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/user/{username}", middlewares.SetMiddlewareJSON(s.UpdateUser)).Methods("PUT")
	s.Router.HandleFunc("/user/{username}", middlewares.SetMiddlewareJSON(s.DeleteUser)).Methods("DELETE")
}
