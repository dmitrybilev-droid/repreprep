package http

import (
	"encoding/json"
	"net/http"

	"golang-forum/internal/user/transport/grpc"
	gen "golang-forum/internal/user/transport/grpc/gen/api"
)

func RegisterAuthRoutes(mux *http.ServeMux, s *grpc.AuthServer) {
	mux.HandleFunc("/api/login", enableCORS(loginHandler(s)))
	mux.HandleFunc("/api/register", enableCORS(registerHandler(s)))
}

func loginHandler(s *grpc.AuthServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req gen.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		resp, err := s.Login(r.Context(), &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		json.NewEncoder(w).Encode(resp)
	}
}

func registerHandler(s *grpc.AuthServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req gen.RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		resp, err := s.Register(r.Context(), &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(resp)
	}
}

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			return
		}
		next(w, r)
	}
}
