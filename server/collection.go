package server

import (
	"encoding/json"
	"net/http"

	"github.com/stinkyfingers/shoppinglistapi/list"
)

func (s *Server) handleUpsertCollection(w http.ResponseWriter, r *http.Request) {
	var c list.Collection
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = c.Upsert(r.Context(), s.Storage)
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	JSON(w, c)
}

func (s *Server) handleGetCollection(w http.ResponseWriter, r *http.Request) {
	c := list.Collection{
		ID: r.PathValue("id"),
	}
	err := c.Get(r.Context(), s.Storage)
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	JSON(w, c)
}

func (s *Server) handleListCollections(w http.ResponseWriter, r *http.Request) {
	collections, err := list.GetCollections(r.Context(), s.Storage)
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	JSON(w, collections)
}

func (s *Server) DeleteCollection(w http.ResponseWriter, r *http.Request) {
	c := list.Collection{
		ID: r.PathValue("id"),
	}
	err := c.Delete(r.Context(), s.Storage)
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	JSON(w, c)
}
