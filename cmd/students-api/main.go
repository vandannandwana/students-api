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

	config "github.com/vandannandwana/students-api/internal"
)

func main() {

	//load config
	cfg := config.MustLoad()

	//database setup

	//server route setup
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(writer http.ResponseWriter, req *http.Request) {
		writer.Write([]byte("Welcome, to Student's API"))
	})

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
