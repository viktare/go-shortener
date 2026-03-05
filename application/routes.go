package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/viktare/go-shortener/handler"
	"github.com/viktare/go-shortener/repository"
)

func loadRoutes(urlRepo *repository.UrlRepository) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/urls", func(r chi.Router) {
		loadUrlRoutes(r, urlRepo)
	})

	return router
}

func loadUrlRoutes(router chi.Router, urlRepo *repository.UrlRepository) {
	urlHandler := &handler.Url{Repo: urlRepo} // inject the repo into the handler

	router.Get("/{shortUrl}", urlHandler.Redirect)
	router.Post("/", urlHandler.Shorten)
}
