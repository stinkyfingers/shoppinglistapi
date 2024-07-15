package server

import (
	"encoding/json"
	"net/http"

	"github.com/stinkyfingers/shoppinglistapi/list"
)

func (s *Server) handleInsertItem(w http.ResponseWriter, r *http.Request) {
	var i list.Item
	err := json.NewDecoder(r.Body).Decode(&i)
	if err != nil {
		Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = i.Insert(r.Context(), s.Storage)
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	JSON(w, i)
}

func (s *Server) handleGetItem(w http.ResponseWriter, r *http.Request) {
	i := list.Item{
		ID: r.URL.Query().Get("id"),
	}
	err := i.Get(r.Context(), s.Storage)
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	JSON(w, i)
}

func (s *Server) handleListItems(w http.ResponseWriter, r *http.Request) {
	items, err := list.GetItems(r.Context(), s.Storage)
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	JSON(w, items)
}

func (s *Server) DeleteItem(w http.ResponseWriter, r *http.Request) {
	i := list.Item{
		ID: r.URL.Query().Get("id"),
	}
	err := i.Delete(r.Context(), s.Storage)
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	JSON(w, i)
}
