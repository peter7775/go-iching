package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/url"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/example/iching-fiber-app/internal/config"
	_ "github.com/example/iching-fiber-app/internal/db"
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

	appURL := serverURL(cfg.Addr)
	log.Printf("listening on %s", cfg.Addr)
	log.Printf("open UI at %s", appURL)

	go func() {
		time.Sleep(700 * time.Millisecond)
		if err := openBrowser(appURL); err != nil {
			log.Printf("browser auto-open failed: %v", err)
		}
	}()

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

func serverURL(addr string) string {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		host = addr
		port = ""
	}
	if host == "" || host == "0.0.0.0" || host == "::" {
		host = "127.0.0.1"
	}
	u := url.URL{Scheme: "http", Host: host}
	if port != "" {
		u.Host = net.JoinHostPort(host, port)
	}
	return u.String()
}

func openBrowser(target string) error {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", target).Start()
	case "darwin":
		return exec.Command("open", target).Start()
	case "linux":
		return exec.Command("xdg-open", target).Start()
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}
