package api

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sureshk/birthday-server/src/api/controllers"
	"github.com/sureshk/birthday-server/src/api/seed"
)

var server = controllers.Server{}

func init() {
	// load values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Printf("sad .env file found")
	}
}

func Run() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}
	server.Initialize(
		os.Getenv("DB_DRIVER"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"))
	seed.Load(server.DB)
	server.Run(":8080")

}
