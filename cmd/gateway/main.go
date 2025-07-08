package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"

	"social-network/graph"
	"social-network/internal/config"
	"social-network/internal/db"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()
	db.Init(cfg)

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	if os.Getenv("ENV") == "production" {
	} else {
		srv.Use(extension.Introspection{})
		srv.AddTransport(transport.POST{})
		http.Handle("/playground", playground.Handler("GraphQL playground", "/graphql"))
		log.Printf("connect to http://%s/playground for GraphQL playground", addr)
	}

	http.Handle("/graphql", srv)

	server := &http.Server{Addr: addr}

	// Start server in a separate goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exiting")
}
