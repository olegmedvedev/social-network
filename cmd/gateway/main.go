package main

import (
	"github.com/joho/godotenv"

	"social-network/internal/config"
	"social-network/internal/server"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()
	server.StartServer(cfg)
}
