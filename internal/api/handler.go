package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rammtw/tgposter/internal/converter"
	"github.com/rammtw/tgposter/internal/poster"
	"github.com/rammtw/tgposter/internal/scheduler"
)

type PostRequest struct {
	Channel  string `json:"channel"`
	Markdown string `json:"markdown"`
	PostAt   string `json:"post_at,omitempty"`
	Timezone string `json:"timezone,omitempty"`
}

type PostResponse struct {
	OK        bool   `json:"ok"`
	MessageID int    `json:"message_id,omitempty"`
	Scheduled string `json:"scheduled_at,omitempty"`
	Error     string `json:"error,omitempty"`
}

func ListenAndServe(addr, token string) error {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/post", postHandler(token))
	mux.HandleFunc("GET /api/v1/health", healthHandler)

	return http.ListenAndServe(addr, mux)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func postHandler(token string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req PostRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, PostResponse{Error: "invalid JSON body"})
			return
		}

		if req.Channel == "" || req.Markdown == "" {
			writeJSON(w, http.StatusBadRequest, PostResponse{Error: "channel and markdown are required"})
			return
		}

		tgText := converter.MarkdownToTelegram(req.Markdown)

		if req.PostAt != "" {
			tz := req.Timezone
			if tz == "" {
				tz = "Europe/Moscow"
			}
			loc, err := time.LoadLocation(tz)
			if err != nil {
				writeJSON(w, http.StatusBadRequest, PostResponse{Error: fmt.Sprintf("invalid timezone: %s", tz)})
				return
			}

			postTime, err := time.ParseInLocation("2006-01-02 15:04", req.PostAt, loc)
			if err != nil {
				writeJSON(w, http.StatusBadRequest, PostResponse{Error: "invalid time format, expected: YYYY-MM-DD HH:MM"})
				return
			}

			if postTime.Before(time.Now()) {
				writeJSON(w, http.StatusBadRequest, PostResponse{Error: "post_at is in the past"})
				return
			}

			s := scheduler.New(token)
			s.Schedule(context.Background(), req.Channel, tgText, postTime)

			writeJSON(w, http.StatusAccepted, PostResponse{
				OK:        true,
				Scheduled: postTime.Format(time.RFC3339),
			})
			return
		}

		p, err := poster.New(token)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, PostResponse{Error: "failed to init bot"})
			return
		}

		msgID, err := p.Send(r.Context(), req.Channel, tgText)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, PostResponse{Error: err.Error()})
			return
		}

		writeJSON(w, http.StatusOK, PostResponse{OK: true, MessageID: msgID})
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
