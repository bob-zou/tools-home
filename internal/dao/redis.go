package dao

import (
	"context"
	"time"
	"tools-home/internal/conf"

	"github.com/go-redis/redis/v8"
)

func fixTimeDuration(cfg *redis.Options) {
	if cfg == nil {
		return
	}
	if cfg.MinRetryBackoff > 0 {
		cfg.MinRetryBackoff = cfg.MinRetryBackoff * time.Second
	}
	if cfg.MaxRetryBackoff > 0 {
		cfg.MaxRetryBackoff = cfg.MaxRetryBackoff * time.Second
	}
	if cfg.DialTimeout > 0 {
		cfg.DialTimeout = cfg.DialTimeout * time.Second
	}
	if cfg.ReadTimeout > 0 {
		cfg.ReadTimeout = cfg.ReadTimeout * time.Second
	}
	if cfg.WriteTimeout > 0 {
		cfg.WriteTimeout = cfg.WriteTimeout * time.Second
	}
	if cfg.MaxConnAge > 0 {
		cfg.MaxConnAge = cfg.MaxConnAge * time.Second
	}
	if cfg.PoolTimeout > 0 {
		cfg.PoolTimeout = cfg.PoolTimeout * time.Second
	}
	if cfg.IdleTimeout > 0 {
		cfg.IdleTimeout = cfg.IdleTimeout * time.Second
	}
	if cfg.IdleCheckFrequency > 0 {
		cfg.IdleCheckFrequency = cfg.IdleCheckFrequency * time.Second
	}
}

func NewRedis() (r *redis.Client, cf func(), err error) {
	var cfg redis.Options
	if err = conf.Load("redis.json", &cfg); err != nil {
		return
	}

	fixTimeDuration(&cfg)
	r = redis.NewClient(&cfg)
	if err = r.Ping(context.Background()).Err(); err != nil {
		return
	}

	cf = func() { _ = r.Close() }
	return
}

func (d *dao) PingRedis(ctx context.Context) (err error) {
	if _, err = d.redis.Do(ctx, "SET", "ping", "pong").Result(); err != nil {
		return
	}
	return
}
