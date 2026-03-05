package application

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/viktare/go-shortener/repository"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "local"
	schema   = "shortener"
)

type App struct {
	router http.Handler
	db     *sql.DB // db now lives on the struct so it stays open
}

func New() *App {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s search_path=%s sslmode=disable",
		host, port, user, password, dbname, schema,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	// Build the repo here, pass it down — this is your "service container"
	urlRepo := repository.NewUrlRepository(db)

	return &App{
		db:     db,
		router: loadRoutes(urlRepo),
	}
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":3000",
		Handler: a.router,
	}

	// Close the DB when Start() returns (server shuts down)
	defer a.db.Close()

	fmt.Println("Server running on :3000")

	if err := server.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}
