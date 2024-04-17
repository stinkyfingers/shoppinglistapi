package server

import (
	"net/http"

	"github.com/stinkyfingers/shoppinglistapi/storage"
	"github.com/stinkyfingers/shoppinglistapi/user"
)

type Server struct {
	Storage storage.Storage
}

func NewServer(profile string) (*Server, error) {
	storage, err := storage.NewS3(profile)
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
	mux.Handle("/test", cors(gcp.Middleware(s.handleTest)))
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

func (s *Server) handleTest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}
