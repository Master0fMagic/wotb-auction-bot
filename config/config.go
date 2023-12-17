package config

import (
	"github.com/caarlos0/env/v10"
	"time"
)

type Config struct {
	LogLevel             string        `env:"LOG_LEVEL" envDefault:"DEBUG"`
	AuctionAPI           string        `env:"AUCTION_API"`
	BotToken             string        `env:"BOT_TOKEN"`
	DbPath               string        `env:"DB_PATH" envDefault:"/data/wotb-bot/wotb.db"`
	AuctionCacheLifetime time.Duration `env:"AUCTION_CACHE_LIFETIME"`
}

func Parse() (Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)
	return cfg, err
}
