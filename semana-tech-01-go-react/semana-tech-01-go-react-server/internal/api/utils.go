package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"odhs/semana-tech-01-go-react-server-main/internal/store/pgstore"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (h apiHandler) readRoom(
	w http.ResponseWriter,
	r *http.Request,
) (room pgstore.Room, rawRoomID string, roomID uuid.UUID, ok bool) {
	// Get and check for valid room ID
	rawRoomID = chi.URLParam(r, "room_id")
	roomID, err := uuid.Parse(rawRoomID)
	if err != nil {
		http.Error(w, "invalid room id", http.StatusBadRequest)
		return pgstore.Room{}, "", uuid.UUID{}, false
	}
	// Verify if room exists
	room, err = h.q.GetRoom(r.Context(), roomID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "room not found", http.StatusBadRequest)
			return pgstore.Room{}, "", uuid.UUID{}, false
		}

		slog.Error("failed to get room", "error", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return pgstore.Room{}, "", uuid.UUID{}, false
	}

	return room, rawRoomID, roomID, true
}

func sendJSON(w http.ResponseWriter, rawData any) {
	data, _ := json.Marshal(rawData)
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(data)

	if err != nil {
		slog.Error("failed to write json", "error", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
}
