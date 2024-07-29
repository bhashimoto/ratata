package routing

import (
	"net/http"

	"github.com/bhashimoto/ratata/handlers"
)

func SetRoutes(cfg *handlers.ApiConfig) (*http.ServeMux ){
	mux := http.NewServeMux()
	mux.HandleFunc("/", cfg.HandleIndex)


	return mux
}
