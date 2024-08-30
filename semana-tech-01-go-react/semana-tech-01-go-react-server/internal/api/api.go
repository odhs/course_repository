package api

import (
	"net/http"

	"odhs/semana-tech-01-go-react-server-main/internal/store/pgstore"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type apiHandler struct {
	// Concrete mode with pgstore, because this app is simple,
	// not necessary an interface to attach another database if necessary.
	q *pgstore.Queries
	r *chi.Mux
}

// Call the handler ServeHTTP inside the Chi framework
func (h apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.r.ServeHTTP(w, r)
}

func NewHandler(q *pgstore.Queries) http.Handler {
	apiH := apiHandler{
		q: q,
	}

	// It could also be NewMux, NewRouter is the same as NewMux
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.RequestID, middleware.Recoverer, middleware.Logger)

	// Basic CORS
	r.Use(cors.Handler(cors.Options{
		// On the production site, the following line can specify the allowed sites
		AllowedOrigins: []string{"https://*", "http://*"},
		// OPTIONS is important because is used by the CORS, the others are to the API
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Routes
	r.Route("/api", func(r chi.Router) {
		r.Route("/rooms", func(r chi.Router) {
			r.Post("/", apiH.handleCreateRoom)
			r.Get("/", apiH.handleGetRooms)

			r.Route("/{room_id}/messages", func(r chi.Router) {
				r.Post("/", apiH.handleCreateRoomMessage)
				r.Get("/", apiH.handleGetRoomMessages)

				r.Route("/{message_id}", func(r chi.Router) {
					r.Get("/", apiH.handleGetRoomMessage)
					r.Patch("/react", apiH.handleReactToMessage)
					r.Delete("/react", apiH.handleRemoveReactFromMessage)
					r.Patch("/answer", apiH.handleMarkMessageAnswered)
				})
			})
		})
	})

	apiH.r = r
	return apiH
}

func (h apiHandler) handleSubscribe(w http.ResponseWriter, r *http.Request)              {}
func (h apiHandler) handleCreateRoom(w http.ResponseWriter, r *http.Request)             {}
func (h apiHandler) handleGetRooms(w http.ResponseWriter, r *http.Request)               {}
func (h apiHandler) handleCreateRoomMessage(w http.ResponseWriter, r *http.Request)      {}
func (h apiHandler) handleGetRoomMessages(w http.ResponseWriter, r *http.Request)        {}
func (h apiHandler) handleGetRoomMessage(w http.ResponseWriter, r *http.Request)         {}
func (h apiHandler) handleReactToMessage(w http.ResponseWriter, r *http.Request)         {}
func (h apiHandler) handleRemoveReactFromMessage(w http.ResponseWriter, r *http.Request) {}
func (h apiHandler) handleMarkMessageAnswered(w http.ResponseWriter, r *http.Request)    {}
