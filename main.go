package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/bhashimoto/ratata/handlers"
	"github.com/bhashimoto/ratata/internal/database"
	"github.com/bhashimoto/ratata/routing"
	"github.com/joho/godotenv"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	log.Println("loading environment variables")
	godotenv.Load()

	port := os.Getenv("PORT")

	apiCfg := handlers.ApiConfig{}

	// https://github.com/libsql/libsql-client-go/#open-a-connection-to-sqld
	// libsql://[your-database].turso.io?authToken=[your-auth-token]
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Println("DATABASE_URL environment variable is not set")
		log.Println("Running without CRUD endpoints")
	} else {
		db, err := sql.Open("libsql", dbURL)
		if err != nil {
			log.Fatal(err)
		}
		dbQueries := database.New(db)
		apiCfg.DB = dbQueries
		log.Println("Connected to database!")
	}

	mux := routing.SetRoutes(&apiCfg)

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

