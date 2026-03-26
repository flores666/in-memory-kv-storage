package storage

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type StorageApi struct {
	service StorageService
}

func NewStorageApi(s StorageService) *StorageApi {
	return &StorageApi{
		service: s,
	}
}

func (s *StorageApi) MapRoutes(r chi.Router) {
	r.Post("/storage/set", s.create)
	r.Post("/storage/get", s.get)
	r.Delete("/storage/delete/{key}", s.delete)
}

func (s *StorageApi) create(w http.ResponseWriter, r *http.Request) {
	request := KeyValue{}

	if err := render.DecodeJSON(r.Body, &request); err != nil {
		render.Status(r, http.StatusBadRequest)
		return
	}

	err := s.service.Set(request.Key, request.Value)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
}

func (s *StorageApi) get(w http.ResponseWriter, r *http.Request) {
	request := ValueRequest{}

	if err := render.DecodeJSON(r.Body, request); err != nil {
		render.Status(r, http.StatusBadRequest)
		return
	}

	value, err := s.service.Get(request.Key)

	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, value)
}

func (s *StorageApi) delete(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")

	if key == "" {
		render.Status(r, http.StatusBadRequest)
		return
	}

	err := s.service.Delete(key)

	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
}
