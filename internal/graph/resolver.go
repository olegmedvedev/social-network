package graph

import (
	"social-network/internal/db"
	"social-network/internal/config"
)

type Resolver struct {
	DB    *db.DB
	Cfg   *config.Config
} 