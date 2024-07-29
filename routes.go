package main

import "net/http"

func SetRoutes() *http.ServeMux{
	mux := http.NewServeMux()
	mux.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})
	return mux
}
