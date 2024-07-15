package server

import (
	"encoding/json"
	"net/http"

	"github.com/stinkyfingers/shoppinglistapi/user"
)

func (s *Server) handleInsertUser(w http.ResponseWriter, r *http.Request) {
	var u user.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = u.Insert(r.Context(), s.Storage)
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	JSON(w, u)
}

func (s *Server) handleGetUser(w http.ResponseWriter, r *http.Request) {
	u := user.User{
		ID: r.URL.Query().Get("id"),
	}
	err := u.Get(r.Context(), s.Storage)
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	JSON(w, u)
}

func (s *Server) handleListUser(w http.ResponseWriter, r *http.Request) {
	users, err := user.ListUsers(r.Context(), s.Storage)
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	JSON(w, users)
}

func (s *Server) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	var u user.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = u.FindByEmail(r.Context(), s.Storage)
	if err != nil {
		if err == user.ErrUserNotFound {
			err = u.Insert(r.Context(), s.Storage)
			if err != nil {
				Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	JSON(w, u)
}
