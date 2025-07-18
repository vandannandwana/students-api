package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	config "github.com/vandannandwana/students-api/internal/config"
	"github.com/vandannandwana/students-api/internal/http/handlers/student"
	"github.com/vandannandwana/students-api/internal/storage/sqlite"
)

func main() {

	//load config
	cfg := config.MustLoad()

	//database setup

	storage, err := sqlite.New(cfg)

	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Storage Initialized", slog.String("env", cfg.Env), slog.String("version:", "1.0.0"))

	//server route setup
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))

	//http server setup
	server := http.Server{
		Addr:    cfg.HttpServer.Address,
		Handler: router,
	}
	slog.Info("Server Started ", slog.String("on:", cfg.HttpServer.Address))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
	}()

	<-done

	slog.Info("Shutting Down Server")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown the server", slog.String("error:", err.Error()))
	}

	slog.Info("Server Shutdown Successfully")

}
