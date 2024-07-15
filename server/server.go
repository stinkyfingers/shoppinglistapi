package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/stinkyfingers/shoppinglistapi/storage"
	"github.com/stinkyfingers/shoppinglistapi/user"
)

type Server struct {
	Storage storage.Storage
}

func NewServer(profile string) (*Server, error) {
	storage, err := storage.NewDynamo(profile)
	if err != nil {
		return nil, err
	}

	return &Server{
		Storage: storage,
	}, nil
}

// NewMux returns the router
func NewMux(s *Server) (http.Handler, error) {
	gcp := &user.GCP{}
	mux := http.NewServeMux()
	mux.Handle("GET /test", cors(gcp.Middleware(s.handleTest)))
	mux.Handle("OPTIONS /", cors(gcp.Middleware(s.handleTest)))
	mux.Handle("POST /user/login", cors(gcp.Middleware(s.handleLoginUser)))
	mux.Handle("POST /user", cors(s.handleInsertUser))              // TODO middleware
	mux.Handle("GET /user", cors(s.handleGetUser))                  // TODO middleware
	mux.Handle("GET /users", cors(s.handleListUser))                // TODO middleware
	mux.Handle("GET /collection/{id}", cors(s.handleGetCollection)) // TODO middleware
	mux.Handle("POST /collection", cors(s.handleUpsertCollection))  // TODO middleware
	mux.Handle("GET /collections", cors(s.handleListCollections))   // TODO middleware
	mux.Handle("POST /list", cors(s.handleInsertList))              // TODO middleware
	mux.Handle("GET /lists", cors(s.handleListLists))               // TODO middleware
	mux.Handle("GET /list", cors(s.handleGetList))                  // TODO middleware
	mux.Handle("POST /item", cors(s.handleInsertItem))              // TODO middleware
	mux.Handle("GET /items", cors(s.handleListItems))               // TODO middleware
	mux.Handle("GET /item", cors(s.handleGetItem))                  // TODO middleware
	return mux, nil
}

func isPermittedOrigin(origin string) string {
	var permittedOrigins = []string{
		"https://shoppinglistapp.john-shenk.com",
		"https://www.shoppinglistapp.john-shenk.com",
		"http://localhost:3000",
		"http://localhost:3001",
		"http://localhost:3002",
		"http://localhost:3003",
	}
	for _, permittedOrigin := range permittedOrigins {
		if permittedOrigin == origin {
			return origin
		}
	}
	return ""
}

func cors(handler http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		permittedOrigin := isPermittedOrigin(r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Origin", permittedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if r.Method == "OPTIONS" {
			return
		}
		next := http.HandlerFunc(handler)
		next.ServeHTTP(w, r)
	})
}

func JSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type RequestError struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

func Error(w http.ResponseWriter, errStr string, code int) {
	w.Header().Set("Content-Type", "application/json")
	e := RequestError{
		Error: errStr,
		Code:  code,
	}
	j, err := json.Marshal(e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(code)
	w.Write(j)
}

func (s *Server) handleTest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("TEST")
	w.Write([]byte("hello"))
}
