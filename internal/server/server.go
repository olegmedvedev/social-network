package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

	"social-network/graph"
	"social-network/internal/config"
	"social-network/internal/db"

	_ "github.com/lib/pq"
)

func StartServer(cfg *config.Config) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)
	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	if err = dbConn.Ping(); err != nil {
		log.Fatal("Failed to ping DB:", err)
	}
	repo := db.NewSocialRepository(dbConn)

	resolvers := &graph.Resolver{Repo: repo, Cfg: cfg}
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolvers}))

	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	if os.Getenv("ENV") != "production" {
		srv.Use(extension.Introspection{})
		srv.AddTransport(transport.POST{})
		http.Handle("/playground", playground.Handler("GraphQL playground", "/graphql"))
		log.Printf("connect to http://%s/playground for GraphQL playground", addr)
	}

	http.Handle("/graphql", AuthMiddleware(srv))

	server := &http.Server{Addr: addr}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

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
