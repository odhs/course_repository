package api

import (
	"net/http"

	"odhs/semana-tech-01-go-react-server-main/internal/store/pgstore"

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

	apiH.r = r
	return apiH
}
