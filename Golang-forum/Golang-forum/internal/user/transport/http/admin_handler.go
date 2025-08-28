package http

import (
	"context"
	"encoding/json"
	"net/http"

	"golang-forum/internal/user/transport/grpc"
	gen "golang-forum/internal/user/transport/grpc/gen/api"

	"google.golang.org/grpc/metadata"
)

func RegisterAdminRoutes(mux *http.ServeMux, s *grpc.AdminServer) {
	mux.HandleFunc("/api/admin/users", enableCORS(getUserListHandler(s)))
	mux.HandleFunc("/api/admin/ban", enableCORS(banHandler(s)))
	mux.HandleFunc("/api/admin/unban", enableCORS(unbanHandler(s)))
	mux.HandleFunc("/api/admin/delete-message", enableCORS(deleteMessageHandler(s)))
}

func getUserListHandler(s *grpc.AdminServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := authContextFromRequest(r)
		resp, err := s.GetUserList(ctx, &gen.GetUserListRequest{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		json.NewEncoder(w).Encode(resp.Users)
	}
}

func banHandler(s *grpc.AdminServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req gen.BanRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		ctx := authContextFromRequest(r)
		resp, err := s.BanUser(ctx, &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		json.NewEncoder(w).Encode(resp)
	}
}

func unbanHandler(s *grpc.AdminServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req gen.UnBanRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		ctx := authContextFromRequest(r)
		resp, err := s.UnBanUser(ctx, &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		json.NewEncoder(w).Encode(resp)
	}
}

func deleteMessageHandler(s *grpc.AdminServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req gen.DeleteMessageRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		ctx := authContextFromRequest(r)
		resp, err := s.DeleteMessage(ctx, &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		json.NewEncoder(w).Encode(resp)
	}
}

func authContextFromRequest(r *http.Request) context.Context {
	token := r.Header.Get("Authorization")
	md := metadata.New(map[string]string{"authorization": token})
	return metadata.NewIncomingContext(context.Background(), md)
}
