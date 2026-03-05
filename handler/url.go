package handler

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/viktare/go-shortener/model"
	"github.com/viktare/go-shortener/repository"
)

type Url struct {
	Repo *repository.UrlRepository
}

type ShortenRequest struct {
	OriginalUrl string `json:"original_url"`
}

func (o *Url) Shorten(w http.ResponseWriter, r *http.Request) {
	var body ShortenRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if body.OriginalUrl == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("original_url is required"))
		return
	}

	hash := fmt.Sprintf("%x", md5.Sum([]byte(body.OriginalUrl)))
	shortUrl := hash[:6]

	url, err := o.Repo.Create(r.Context(), model.Url{
		OriginalUrl: body.OriginalUrl,
		ShortUrl:    shortUrl,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"short_url": fmt.Sprintf("http://localhost:3000/urls/%s", url.ShortUrl),
	})
}

func (o *Url) Redirect(w http.ResponseWriter, r *http.Request) {
	shortUrl := chi.URLParam(r, "shortUrl")
	if shortUrl == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	url, err := o.Repo.FindByShortUrl(r.Context(), shortUrl)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url.OriginalUrl, http.StatusMovedPermanently)
}
