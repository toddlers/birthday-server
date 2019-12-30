package controllers

import (
	"net/http"

	"github.com/toddlers/birthday-server/src/api/responses"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome To The Birthday Server")
}
