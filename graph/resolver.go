package graph

import (
	"social-network/internal/config"
	"social-network/internal/db"
)

type Resolver struct {
	Repo *db.SocialRepository
	Cfg  *config.Config
}
