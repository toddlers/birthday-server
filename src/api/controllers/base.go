package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/toddlers/birthday-server/src/api/models"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

func (server *Server) Initialize(DbDriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {
	var err error
	if DbDriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
		server.DB, err = gorm.Open(DbDriver, DBURL)
		if err != nil {
			fmt.Printf("Can not connect to %s database", DbDriver)
			log.Fatal("This is the error: ", err)
		} else {
			fmt.Printf("We are connected to the %s database", DbDriver)
		}
	}
	server.DB.Debug().AutoMigrate(&models.User{})
	server.Router = mux.NewRouter()
	server.initializeRoutes()
}

func (server *Server) makeHTTPServer(ip, port string) *http.Server {
	srv := &http.Server{
		//the maximum duration for reading the entire request, including the body
		ReadTimeout: 1 * time.Second,
		//the maximum duration before timing out writes of the response
		WriteTimeout: 1 * time.Second,
		//the maximum amount of time to wait for the next request
		//when keep-alive is enabled
		IdleTimeout: 30 * time.Second,
		//the amount of time allowed to read request headers
		ReadHeaderTimeout: 2 * time.Second,
		Handler:           server.Router,
		Addr:              fmt.Sprintf("%s:%s",ip, port),
	}
	return srv
}

func (server *Server) Run(ip, port string) {
	srv := server.makeHTTPServer(ip, port)
	fmt.Println("Listening to port 8080")
	log.Fatal(srv.ListenAndServe())
}
