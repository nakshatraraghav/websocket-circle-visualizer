package api

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type APIServer struct {
	router *chi.Mux
	addr   string
}

func NewAPIServer() *APIServer {
	router := chi.NewRouter()

	return &APIServer{
		addr:   ":3000",
		router: router,
	}
}

func (api *APIServer) StartServer() error {
	log.Println("server started on localhost:3000")
	return http.ListenAndServe(api.addr, api.router)
}
