package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/pocwithmehul/common-go-lib"
	"github.com/segmentio/kafka-go"
)

type StockEvent struct {
	Ticker    string    `json:"ticker"`
	Price     float64   `json:"price"`
	Timestamp time.Time `json:"timestamp"`
	Volume    int64     `json:"volume"`
}

func CallbackHandler(writer *kafka.Writer, logger *commonlib.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var event StockEvent
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			logger.Error("invalid request body", map[string]interface{}{"error": err.Error()})
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}

		if event.Ticker == "" {
			http.Error(w, "missing ticker", http.StatusBadRequest)
			return
		}

		payload, err := json.Marshal(event)
		if err != nil {
			logger.Error("serialize event", map[string]interface{}{"error": err.Error()})
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		msg := kafka.Message{
			Key:   []byte(event.Ticker),
			Value: payload,
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		if err := writer.WriteMessages(ctx, msg); err != nil {
			logger.Error("kafka publish failed", map[string]interface{}{"error": err.Error(), "ticker": event.Ticker})
			http.Error(w, "failed to publish event", http.StatusInternalServerError)
			return
		}

		logger.Info("published stock event", map[string]interface{}{"ticker": event.Ticker, "price": event.Price})
		w.WriteHeader(http.StatusNoContent)
	}
}
