package handlers

import "net/http"

func (cfg *ApiConfig) HandleIndex(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, 200, "Hello")
}
