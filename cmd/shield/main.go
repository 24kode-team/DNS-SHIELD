package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/24kode-team/DNS-SHIELD/internal/api"
	"github.com/24kode-team/DNS-SHIELD/internal/blocklist"
	"github.com/24kode-team/DNS-SHIELD/internal/config"
	"github.com/24kode-team/DNS-SHIELD/internal/filter"
	"github.com/24kode-team/DNS-SHIELD/internal/metrics"
	"github.com/24kode-team/DNS-SHIELD/internal/resolver"
	"go.uber.org/zap"
)

func main() {
	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer func() { _ = log.Sync() }()

	cfg, err := config.Load("configs/shield.yaml")
	if err != nil {
		log.Fatal("failed to load config", zap.Error(err))
	}

	m := metrics.New()

	bl, err := blocklist.NewManager(cfg.Blocklists, log)
	if err != nil {
		log.Fatal("failed to init blocklist manager", zap.Error(err))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go bl.StartRefreshLoop(ctx)

	fe := filter.NewEngine(bl, cfg.Filter, log)

	res, err := resolver.New(cfg.Resolver, fe, m, log)
	if err != nil {
		log.Fatal("failed to init resolver", zap.Error(err))
	}

	go func() {
		if err := res.ListenAndServe(ctx); err != nil {
			log.Error("resolver stopped", zap.Error(err))
		}
	}()

	srv := api.New(cfg.API, bl, fe, m, log, cfg.Resolver.Upstreams)
	go func() {
		if err := srv.Start(); err != nil {
			log.Error("api stopped", zap.Error(err))
		}
	}()

	log.Info("DNS-SHIELD running",
		zap.String("dns", cfg.Resolver.ListenAddr),
		zap.String("dot", cfg.Resolver.DoTAddr),
		zap.String("doh", cfg.API.ListenAddr+cfg.Resolver.DoHPath),
		zap.String("dashboard", cfg.API.ListenAddr),
	)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down...")
	cancel()
	res.Shutdown()
	srv.Shutdown()
}
