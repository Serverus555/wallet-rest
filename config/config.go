package config

import (
	"wallet-rest/internal/repository/cache"
	"wallet-rest/pkg/httpserver"
	"wallet-rest/pkg/postgres"
)

type Config struct {
	Postgres postgres.Config
	Http     httpserver.HttpConfig
	Cache    cache.Config
}
