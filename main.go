package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("loading environment variables")
	godotenv.Load()

	port := os.Getenv("PORT")


	mux := SetRoutes()

	server := http.Server{
		Addr: ":" + port,
		Handler: mux,
	}

	
	log.Println("Listening on port", port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

