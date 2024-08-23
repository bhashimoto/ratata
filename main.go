package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/bhashimoto/ratata/api"
	"github.com/bhashimoto/ratata/front"
	"github.com/bhashimoto/ratata/internal/database"
	"github.com/bhashimoto/ratata/routing"
	"github.com/joho/godotenv"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	log.Println("loading environment variables")
	godotenv.Load()

	port := os.Getenv("PORT")

	log.Println("creating config structs")
	apiCfg := api.ApiConfig{
		AccountCache: make(map[string]*api.AccountData),
	}
	webCfg := front.WebAppConfig{
		BaseURL:   "http://localhost:8080/api/",
	}
	baseUrl := os.Getenv("BACKEND_BASE_URL")
	staticRoot := "./static/"
	
	log.Println("initializing webCfg")
	webCfg.Init(staticRoot, baseUrl)

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

	mux := http.NewServeMux()
	log.Println("Setting API routes")
	err := routing.SetApiRoutes(&apiCfg, mux)
	if err != nil {
		log.Fatal("Could not set API routes:", err)
	}

	log.Println("Setting front-end routes")
	err = routing.SetFrontEndRoutes(&webCfg, mux)
	if err != nil {
		log.Fatal("Could not set front-end routes:", err)
	}

	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Println("Listening on port", port)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
