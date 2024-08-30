package api

import (
	"net/http"

	"odhs/semana-tech-01-go-react-server-main/internal/store/pgstore"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type apiHandler struct {
	// Concrete mode, because this app is simple, not necessary a interface
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

func (h apiHandler) handleCreateRoom(w http.ResponseWriter, r *http.Request)             {}
func (h apiHandler) handleGetRooms(w http.ResponseWriter, r *http.Request)               {}
func (h apiHandler) handleCreateRoomMessage(w http.ResponseWriter, r *http.Request)      {}
func (h apiHandler) handleGetRoomMessages(w http.ResponseWriter, r *http.Request)        {}
func (h apiHandler) handleGetRoomMessage(w http.ResponseWriter, r *http.Request)         {}
func (h apiHandler) handleReactToMessage(w http.ResponseWriter, r *http.Request)         {}
func (h apiHandler) handleRemoveReactFromMessage(w http.ResponseWriter, r *http.Request) {}
func (h apiHandler) handleMarkMessageAnswered(w http.ResponseWriter, r *http.Request)    {}
