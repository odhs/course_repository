package api

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"sync"

	"odhs/semana-tech-01-go-react-server-main/internal/store/pgstore"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5"
)

type apiHandler struct {
	// Concrete mode with pgstore, because this app is simple,
	// not necessary an interface to attach another database if necessary.
	q        *pgstore.Queries
	r        *chi.Mux
	upgrader websocket.Upgrader
	// Pool de Conexões para clientes conectados via WebSocket
	// map[roomId] por map[clientes conectados] retornando uma função Cancel
	subscribers map[string]map[*websocket.Conn]context.CancelFunc
	// Blocks data race
	mu *sync.Mutex
}

// Call the handler ServeHTTP inside the Chi framework
func (h apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.r.ServeHTTP(w, r)
}

func NewHandler(q *pgstore.Queries) http.Handler {
	a := apiHandler{
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
	r.Get("/subscribe/{room_id}", a.handleSubscribe)

	// Routes
	r.Route("/api", func(r chi.Router) {
		r.Route("/rooms", func(r chi.Router) {
			r.Post("/", a.handleCreateRoom)
			r.Get("/", a.handleGetRooms)

			r.Route("/{room_id}", func(r chi.Router) {
				r.Get("/", a.handleGetRoom)

				r.Route("/messages", func(r chi.Router) {
					r.Post("/", a.handleCreateRoomMessage)
					r.Get("/", a.handleGetRoomMessages)

					r.Route("/{message_id}", func(r chi.Router) {
						r.Get("/", a.handleGetRoomMessage)
						r.Patch("/react", a.handleReactToMessage)
						r.Delete("/react", a.handleRemoveReactFromMessage)
						r.Patch("/answer", a.handleMarkMessageAsAnswered)
					})
				})
			})
		})
	})

	a.r = r
	return a
}

const (
	MessageKindMessageCreated          = "message_created"
	MessageKindMessageRactionIncreased = "message_reaction_increased"
	MessageKindMessageRactionDecreased = "message_reaction_decreased"
	MessageKindMessageAnswered         = "message_answered"
)

type MessageMessageReactionIncreased struct {
	ID    string `json:"id"`
	Count int64  `json:"count"`
}

type MessageMessageReactionDecreased struct {
	ID    string `json:"id"`
	Count int64  `json:"count"`
}

type MessageMessageAnswered struct {
	ID string `json:"id"`
}

type MessageMessageCreated struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

type Message struct {
	Kind string `json:"kind"`
	/* Value it will be structs above  */
	Value any `json:"value"`
	/* "-" = Don't encode RoomID */
	RoomID string `json:"-"`
}

func (h apiHandler) notifyClients(msg Message) {
	h.mu.Lock()
	defer h.mu.Unlock()

	subscribers, ok := h.subscribers[msg.RoomID]
	if !ok || len(subscribers) == 0 {
		return
	}

	for conn, cancel := range subscribers {
		if err := conn.WriteJSON(msg); err != nil {
			slog.Error("failed to send message to client", "error", err)
			// o cliente não existe mais então vamos cancelar o contexto,
			// que por sua vez irá cancelar a conexão e retirar o subscriber do pool
			cancel()
		}
	}
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
	_, rawRoomID, _, ok := h.readRoom(w, r)
	if !ok {
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

	// Context with the Cancel function, to be able to release the function if
	// the connection be interrupted.
	ctx, cancel := context.WithCancel(r.Context())

	// Updating the customer map with the new subscriber that has just arrived
	// If room not exists webSocket it will be created
	// If room exists already only attributes the "CANCEL" function.
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

func (h apiHandler) handleCreateRoom(w http.ResponseWriter, r *http.Request) {

	// Read JSON
	type _body struct {
		Theme string `json:"theme"`
	}
	var body _body
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
	}

	// Writes to the database
	roomID, err := h.q.InsertRoom(r.Context(), body.Theme)
	if err != nil {
		slog.Error("failed to insert room", "error", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
	}

	// Response with ID
	type response struct {
		ID string `json:"id"`
	}
	sendJSON(w, response{ID: roomID.String()})
}

func (h apiHandler) handleGetRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := h.q.GetRooms(r.Context())
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		slog.Error("failed to get rooms", "error", err)
		return
	}

	if rooms == nil {
		rooms = []pgstore.Room{}
	}

	sendJSON(w, rooms)
}

func (h apiHandler) handleGetRoom(w http.ResponseWriter, r *http.Request) {
	room, _, _, ok := h.readRoom(w, r)
	if !ok {
		return
	}

	sendJSON(w, room)
}

func (h apiHandler) handleCreateRoomMessage(w http.ResponseWriter, r *http.Request) {
	// Check if the room exists
	_, rawRoomID, roomID, ok := h.readRoom(w, r)
	if !ok {
		return
	}

	// Read JSON
	type _body struct {
		Message string `json:"message"`
	}
	var body _body
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
	}

	// Writes to the database
	messageID, err := h.q.InsertMessage(
		r.Context(),
		pgstore.InsertMessageParams{
			RoomID:  roomID,
			Message: body.Message,
		},
	)
	if err != nil {
		slog.Error("failed to insert message", "error", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
	}

	type response struct {
		ID string `json:"id"`
	}

	sendJSON(w, response{ID: messageID.String()})

	// Notify the clients who are connected to the room
	go h.notifyClients(
		Message{
			Kind:   MessageKindMessageCreated,
			RoomID: rawRoomID,
			Value: MessageMessageCreated{
				ID:      messageID.String(),
				Message: body.Message,
			},
		},
	)

}
func (h apiHandler) handleGetRoomMessages(w http.ResponseWriter, r *http.Request) {
	// Check if the room exists
	_, _, roomID, ok := h.readRoom(w, r)
	if !ok {
		return
	}

	// Get from database
	messages, err := h.q.GetRoomMessages(r.Context(), roomID)
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		slog.Error("failed to get room messages", "error", err)
		return
	}

	if messages == nil {
		messages = []pgstore.Message{}
	}

	sendJSON(w, messages)
}
func (h apiHandler) handleGetRoomMessage(w http.ResponseWriter, r *http.Request) {
	// Check if the room exists
	_, _, _, ok := h.readRoom(w, r)
	if !ok {
		return
	}

	rawMessageID := chi.URLParam(r, "message_id")
	messageID, err := uuid.Parse(rawMessageID)
	if err != nil {
		http.Error(w, "invalid message id", http.StatusBadRequest)
		return
	}

	// Get from database
	messages, err := h.q.GetMessage(r.Context(), messageID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "message not found", http.StatusBadRequest)
			return
		}

		slog.Error("failed to get message", "error", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	sendJSON(w, messages)
}

func (h apiHandler) handleReactToMessage(w http.ResponseWriter, r *http.Request) {
	_, rawRoomID, _, ok := h.readRoom(w, r)
	if !ok {
		return
	}

	rawID := chi.URLParam(r, "message_id")
	id, err := uuid.Parse(rawID)
	if err != nil {
		http.Error(w, "invalid message id", http.StatusBadRequest)
		return
	}

	// Set in database
	count, err := h.q.ReactToMessage(r.Context(), id)
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		slog.Error("failed to react to message", "error", err)
		return
	}

	type response struct {
		Count int64 `json:"count"`
	}

	sendJSON(w, response{Count: count})

	go h.notifyClients(Message{
		Kind:   MessageKindMessageRactionIncreased,
		RoomID: rawRoomID,
		Value: MessageMessageReactionIncreased{
			ID:    rawID,
			Count: count,
		},
	})
}

func (h apiHandler) handleRemoveReactFromMessage(w http.ResponseWriter, r *http.Request) {
	_, rawRoomID, _, ok := h.readRoom(w, r)
	if !ok {
		return
	}

	rawID := chi.URLParam(r, "message_id")
	id, err := uuid.Parse(rawID)
	if err != nil {
		http.Error(w, "invalid message id", http.StatusBadRequest)
		return
	}

	// Remove from database
	count, err := h.q.RemoveReactionFromMessage(r.Context(), id)
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		slog.Error("failed to react to message", "error", err)
		return
	}

	type response struct {
		Count int64 `json:"count"`
	}

	sendJSON(w, response{Count: count})

	go h.notifyClients(Message{
		Kind:   MessageKindMessageRactionDecreased,
		RoomID: rawRoomID,
		Value: MessageMessageReactionDecreased{
			ID:    rawID,
			Count: count,
		},
	})
}

func (h apiHandler) handleMarkMessageAsAnswered(w http.ResponseWriter, r *http.Request) {
	_, rawRoomID, _, ok := h.readRoom(w, r)
	if !ok {
		return
	}

	rawID := chi.URLParam(r, "message_id")
	id, err := uuid.Parse(rawID)
	if err != nil {
		http.Error(w, "invalid message id", http.StatusBadRequest)
		return
	}

	// Mark in the database
	err = h.q.MarkMessageAsAnswered(r.Context(), id)
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		slog.Error("failed to react to message", "error", err)
		return
	}

	w.WriteHeader(http.StatusOK)

	go h.notifyClients(Message{
		Kind:   MessageKindMessageAnswered,
		RoomID: rawRoomID,
		Value: MessageMessageAnswered{
			ID: rawID,
		},
	})
}
