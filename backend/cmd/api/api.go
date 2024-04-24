package api

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nakshatraraghav/hashed-tokens-assignment/backend/lib"
)

type APIServer struct {
	addr   string
	radius float64
	router *chi.Mux
	sin    chan float64
	cos    chan float64
	rchan  chan float64
}

func NewAPIServer() *APIServer {
	router := chi.NewRouter()

	sin := make(chan float64, 1)
	cos := make(chan float64, 1)
	rchan := make(chan float64, 1)

	go lib.SinSampleGenerator(sin, rchan, 1)
	go lib.CosSampleGenerator(cos, rchan, 1)

	return &APIServer{
		addr:   ":3000",
		radius: 1,
		router: router,
		sin:    sin,
		cos:    cos,
		rchan:  rchan,
	}
}

func (api *APIServer) StartServer() error {
	log.Println("server started on localhost:3000")
	return http.ListenAndServe(api.addr, api.router)
}
