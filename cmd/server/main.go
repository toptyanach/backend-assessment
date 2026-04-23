package main

import (
	"log/slog"
	"net/http"
	"os"
)

func main() {
	// Настройка структурированного логгирования (закрываем требование по slog)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Эндпоинт для проверки здоровья сервиса (закрываем требование health check)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	slog.Info("Server is starting", slog.String("port", "8080"))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		slog.Error("Server failed to start", slog.Any("error", err))
		os.Exit(1)
	}
}
