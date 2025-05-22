package main

import (
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/AntonChubarov/goit-cloud-fp/internal/config"
	"github.com/AntonChubarov/goit-cloud-fp/internal/handler"
	"github.com/AntonChubarov/goit-cloud-fp/internal/repository/postgres"
	"github.com/AntonChubarov/goit-cloud-fp/internal/service"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed dist/*
var webFS embed.FS

//go:embed migrations/*.sql
var migrationFiles embed.FS

func main() {
	cfg := config.FromEnv()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.LogLevel,
	}))
	slog.SetDefault(logger)
	slog.Info("configuration loaded")

	db := waitForDB(cfg.DB_DSN, 10, 2*time.Second)
	defer func() {
		if err := db.Close(); err != nil {
			slog.Warn("unable to close database", "error", err)
		}
	}()

	applyMigrations(db)

	repo := postgres.NewLinkRepo(db)
	slog.Debug("repository initialized")

	svc := service.NewShortener(repo)
	slog.Debug("service initialized")

	sub, err := fs.Sub(webFS, "dist")
	if err != nil {
		slog.Error("unable to access embedded filesystem", "error", err)
		os.Exit(1)
	}
	slog.Debug("embedded static filesystem mounted")

	h := handler.New(svc, sub)
	r := chi.NewRouter()
	r.Mount("/", h.Routes())
	slog.Info("routes mounted")

	slog.Info("starting server", "port", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		slog.Error("unable to start HTTP server", "error", err)
		os.Exit(1)
	}
}

// waitForDB attempts to connect to the database with retry logic.
func waitForDB(dsn string, maxAttempts int, delay time.Duration) *sqlx.DB {
	var db *sqlx.DB
	var err error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		db, err = sqlx.Connect("pgx", dsn)
		if err == nil {
			slog.Info("connected to database")
			return db
		}
		slog.Warn("unable to connect to database, will retry", "attempt", attempt, "error", err)
		time.Sleep(delay)
	}

	slog.Error("unable to connect to database after retries", "error", err)
	os.Exit(1)
	return nil
}

func applyMigrations(db *sqlx.DB) {
	driver, err := pgx.WithInstance(db.DB, &pgx.Config{})
	if err != nil {
		slog.Error("unable to create migration driver", "error", err)
		os.Exit(1)
	}

	d, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		slog.Error("unable to load embedded migrations", "error", err)
		os.Exit(1)
	}

	m, err := migrate.NewWithInstance("iofs", d, "pgx", driver)
	if err != nil {
		slog.Error("unable to create migrate instance", "error", err)
		os.Exit(1)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		slog.Error("unable to apply migrations", "error", err)
		os.Exit(1)
	}
	slog.Info("migrations applied successfully")
}
