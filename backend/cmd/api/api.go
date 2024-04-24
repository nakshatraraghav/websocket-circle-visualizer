package api

import (
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/nakshatraraghav/hashed-tokens-assignment/backend/lib"
	"github.com/nakshatraraghav/hashed-tokens-assignment/backend/ws"
)

type APIServer struct {
	addr   string
	radius float64
	router *chi.Mux
	hub    *ws.Hub
	sin    chan float64
	cos    chan float64
	rchan  chan float64
}

func NewAPIServer() *APIServer {
	router := chi.NewRouter()

	hub := ws.NewHub()

	sin := make(chan float64, 1)
	cos := make(chan float64, 1)
	rchan := make(chan float64, 1)

	go hub.RunHub(sin, cos)

	go lib.SinSampleGenerator(sin, rchan, 1)
	go lib.CosSampleGenerator(cos, rchan, 1)

	return &APIServer{
		addr:   ":3000",
		radius: 1,
		router: router,
		hub:    hub,
		sin:    sin,
		cos:    cos,
		rchan:  rchan,
	}
}

func (api *APIServer) StartServer() error {

	api.registerMiddlewares()
	api.registerRoutes()

	log.Println("server started on localhost:3000")
	return http.ListenAndServe(api.addr, api.router)
}

func (api *APIServer) registerRoutes() {
	router := api.router

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	router.Get("/sin", func(w http.ResponseWriter, r *http.Request) {
		sin := lib.FloatToString(<-api.sin)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(sin))
	})

	router.Get("/cos", func(w http.ResponseWriter, r *http.Request) {
		cos := lib.FloatToString(<-api.cos)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(cos))
	})

	router.Post("/update_radius", func(w http.ResponseWriter, r *http.Request) {
		radius, _ := io.ReadAll(r.Body)
		r.Body.Close()

		api.rchan <- lib.StringToFloat(string(radius))
		log.Println("radius updated samples are scaled up now")

	})

	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWebSocket(api.hub, w, r)
	})
}

func (api *APIServer) registerMiddlewares() {
	router := api.router

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	router.Use(middleware.Logger)
}
