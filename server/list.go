package server

import (
	"encoding/json"
	"net/http"

	"github.com/stinkyfingers/shoppinglistapi/list"
)

func (s *Server) handleInsertList(w http.ResponseWriter, r *http.Request) {
	var l list.List
	err := json.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = l.Insert(r.Context(), s.Storage)
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	JSON(w, l)
}

func (s *Server) handleGetList(w http.ResponseWriter, r *http.Request) {
	l := list.List{
		ID: r.URL.Query().Get("id"),
	}
	err := l.Get(r.Context(), s.Storage)
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	JSON(w, l)
}

func (s *Server) handleListLists(w http.ResponseWriter, r *http.Request) {
	lists, err := list.GetLists(r.Context(), s.Storage)
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	JSON(w, lists)
}

func (s *Server) DeleteList(w http.ResponseWriter, r *http.Request) {
	l := list.List{
		ID: r.URL.Query().Get("id"),
	}
	err := l.Delete(r.Context(), s.Storage)
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	JSON(w, l)
}
