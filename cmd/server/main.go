package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/AntonChubarov/goit-cloud-fp/internal/config"
	"github.com/AntonChubarov/goit-cloud-fp/internal/handler"
	"github.com/AntonChubarov/goit-cloud-fp/internal/repository/postgres"
	"github.com/AntonChubarov/goit-cloud-fp/internal/service"
)

//go:embed dist/*
var webFS embed.FS

func main() {
	cfg := config.FromEnv()

	db, err := sqlx.Connect("pgx", cfg.DB_DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := postgres.NewLinkRepo(db)
	svc := service.NewShortener(repo)

	sub, err := fs.Sub(webFS, "dist")
	if err != nil {
		log.Fatal("failed to sub embed FS:", err)
	}

	h := handler.New(svc, sub)

	r := chi.NewRouter()
	r.Mount("/", h.Routes())

	log.Printf("listening on :%s …", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
