package main

import (
	"context"
	"database/sql"
	"log"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/example/iching-fiber-app/internal/db"
	"github.com/example/iching-fiber-app/internal/config"
	"github.com/example/iching-fiber-app/internal/httpfiber"
	"github.com/example/iching-fiber-app/internal/service"
	mem "github.com/example/iching-fiber-app/internal/storage/memory"
	pgrepo "github.com/example/iching-fiber-app/internal/storage/postgres"
	sqliterepo "github.com/example/iching-fiber-app/internal/storage/sqlite"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()
	repo, cleanup := mustRepository(cfg)
	defer cleanup()

	svc := service.NewReadingService(repo)
	app := httpfiber.NewApp(cfg, svc)

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = app.ShutdownWithContext(shutdownCtx)
	}()

	log.Printf("listening on %s", cfg.Addr)
	if err := app.Listen(cfg.Addr); err != nil {
		log.Fatal(err)
	}
}

func mustRepository(cfg config.Config) (service.ReadingRepository, func()) {
	switch cfg.Storage {
	case "postgres":
		db, err := sql.Open("pgx", cfg.PostgresDSN)
		if err != nil {
			log.Fatal(err)
		}
		if err := db.Ping(); err != nil {
			log.Fatal(err)
		}
		return pgrepo.NewReadingRepository(db), func() { _ = db.Close() }
	case "sqlite":
		db, err := sql.Open("sqlite", cfg.SQLitePath)
		if err != nil {
			log.Fatal(err)
		}
		if err := db.Ping(); err != nil {
			log.Fatal(err)
		}
		return sqliterepo.NewReadingRepository(db), func() { _ = db.Close() }
	default:
		return mem.NewReadingRepository(), func() {}
	}
}
