package main

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hadv/go-charity-me/internal/handler"
	"github.com/hadv/go-charity-me/internal/repo"
	"github.com/hadv/go-charity-me/internal/service"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

func main() {
	viper.AutomaticEnv()
	dbURL := viper.GetString("DB_URL")
	db, err := sqlx.Connect("mysql", dbURL)
	if err != nil {
		log.Fatalf("Cannot connect to MySQL at %v: %v", dbURL, err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			log.Printf("Error while closing DB connection: %v", err)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var srv http.Server
	done := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		// interrupt signal sent from terminal
		signal.Notify(sigint, os.Interrupt)
		// sigterm signal sent from kubernetes
		signal.Notify(sigint, syscall.SIGTERM)

		<-sigint

		// We received an interrupt signal, shut down gracefully
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(done)
	}()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	accountHandler := handler.NewAccount(service.NewAccount(repo.NewUser(db)))
	r.Post("/signin", accountHandler.Login)
	r.Put("/register", accountHandler.Register)

	srv = http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-done
}
