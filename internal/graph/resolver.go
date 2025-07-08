package graph

import (
	"database/sql"
	"social-network/internal/config"
)

type Resolver struct {
	DB  *sql.DB
	Cfg *config.Config
}
