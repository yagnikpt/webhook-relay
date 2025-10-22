package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		githubToken := authHeader[7:]

		client := &http.Client{Timeout: 10 * time.Second}
		req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		req.Header.Set("Authorization", "Bearer "+githubToken)
		req.Header.Set("Accept", "application/vnd.github+json")

		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "Failed to validate token", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			http.Error(w, "Invalid GitHub token", http.StatusUnauthorized)
			return
		}

		var userInfo map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
			http.Error(w, "Failed to parse user info", http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", fmt.Sprintf("%d", int(userInfo["id"].(float64))))
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
