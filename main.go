package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	commonlogger "github.com/pocwithmehul/common-go-lib/pkg/logger"
	"github.com/pocwithmehul/stock-webhook-service/internal/config"
	"github.com/pocwithmehul/stock-webhook-service/internal/handler"
	"github.com/pocwithmehul/stock-webhook-service/internal/messaging"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	logger := commonlogger.NewLogger("stock-webhook-service", cfg.Datadog)
	writer := messaging.NewKafkaWriter(cfg)
	defer func() {
		if err := writer.Close(); err != nil {
			logger.Error("kafka writer close failed", map[string]interface{}{"error": err.Error()})
		}
	}()

	router := mux.NewRouter()
	router.HandleFunc("/v1/stock/callbacks", handler.CallbackHandler(writer, logger)).Methods(http.MethodPost)

	port := cfg.Server.Port
	if port == 0 {
		port = 8090
	}
	addr := fmt.Sprintf(":%d", port)
	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info("starting stock-webhook-service", map[string]interface{}{"addr": addr})
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("server failed", map[string]interface{}{"error": err.Error()})
		log.Fatalf("server failed: %v", err)
	}
}
