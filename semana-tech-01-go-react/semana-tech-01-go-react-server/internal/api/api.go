package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"sync"

	"odhs/semana-tech-01-go-react-server-main/internal/store/pgstore"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/gorilla/websocket"
)

type apiHandler struct {
	// Concrete mode with pgstore, because this app is simple,
	// not necessary an interface to attach another database if necessary.
	q        *pgstore.Queries
	r        *chi.Mux
	upgrader websocket.Upgrader
	// Pool de Conexões para clientes conectados via WebSocket
	// map[roomId] por map[clientes conectados]
	subscribers map[string]map[*websocket.Conn]context.CancelFunc
	// Blocks data race
	mu *sync.Mutex
}

// Call the handler ServeHTTP inside the Chi framework
func (h apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.r.ServeHTTP(w, r)
}

func NewHandler(q *pgstore.Queries) http.Handler {
	apiH := apiHandler{
		q: q,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		subscribers: make(map[string]map[*websocket.Conn]context.CancelFunc),
		mu:          &sync.Mutex{},
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

	// Websocket
	r.Get("/subscribe/{room_id}", apiH.handleSubscribe)

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

/*
handleSubscribe
- It is a function that blocks execution because it maintains the connection
with the user open, waiting the user or the server cancel it. So we use the
customer flag to cancel the context, canceling the context,
we canceled the connection.
*/
func (h apiHandler) handleSubscribe(w http.ResponseWriter, r *http.Request) {
	// Get and check for valid room ID
	rawRoomID := chi.URLParam(r, "room_id")
	roomID, err := uuid.Parse(rawRoomID)
	if err != nil {
		http.Error(w, "invalid room id", http.StatusBadRequest)
	}

	// Verify if room exists
	_, err = h.q.GetRoom(r.Context(), roomID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "room not found", http.StatusBadRequest)
			return
		}
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	// Until now the connection is a normal HTTP Request,
	// Let's transform this connection in a websocket
	c, err := h.upgrader.Upgrade(w, r, nil)
	// Let's verify why the cliente is not able to upgrade
	if err != nil {
		slog.Warn("failed to upgrade connection", "error", err)
		http.Error(w, "failed to upgrade to ws connection", http.StatusBadRequest)
		return
	}
	defer c.Close()

	// Contexto com a função cancel, para poder liberar a função caso a conexão
	// seja interrompida.
	ctx, cancel := context.WithCancel(r.Context())

	// Atualizando o mapa de clientes com o novo subscriber que acabou de chegar
	// Se ainda não existir o websocket da sala cria
	// Se já existir apenas atribui a função "cancel".
	h.mu.Lock()
	if _, ok := h.subscribers[rawRoomID]; !ok {
		h.subscribers[rawRoomID] = make(map[*websocket.Conn]context.CancelFunc)
	}
	slog.Info(
		"new cliente connected",
		"room_id", rawRoomID,
		"cliente_ip",
		r.RemoteAddr,
	)
	h.subscribers[rawRoomID][c] = cancel
	h.mu.Unlock()

	// O programa ficará aguardando aqui a conexão ser cancelada
	// ou pelo cliente ou pelo servidor.
	<-ctx.Done()

	// Quando o programa chegar aqui quer dizer que a conexão foi cancelada
	// Então vamos retirar a conexão do pool de conexões
	h.mu.Lock()
	delete(h.subscribers[rawRoomID], c)
	h.mu.Unlock()
}

func (h apiHandler) handleCreateRoom(w http.ResponseWriter, r *http.Request)             {}
func (h apiHandler) handleGetRooms(w http.ResponseWriter, r *http.Request)               {}
func (h apiHandler) handleCreateRoomMessage(w http.ResponseWriter, r *http.Request)      {}
func (h apiHandler) handleGetRoomMessages(w http.ResponseWriter, r *http.Request)        {}
func (h apiHandler) handleGetRoomMessage(w http.ResponseWriter, r *http.Request)         {}
func (h apiHandler) handleReactToMessage(w http.ResponseWriter, r *http.Request)         {}
func (h apiHandler) handleRemoveReactFromMessage(w http.ResponseWriter, r *http.Request) {}
func (h apiHandler) handleMarkMessageAnswered(w http.ResponseWriter, r *http.Request)    {}
