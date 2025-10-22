package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func createToken(w http.ResponseWriter, r *http.Request) {
	jwt_secret := os.Getenv("JWT_SECRET")
	user_id := r.URL.Query().Get("user_id")
	if user_id == "" {
		http.Error(w, "Missing user_id parameter", http.StatusBadRequest)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user_id,
		"exp": jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
	})

	secretKey := []byte(jwt_secret)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		http.Error(w, "Could not create token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwt_secret := os.Getenv("JWT_SECRET")
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := authHeader[len("Bearer "):]
		secretKey := []byte(jwt_secret)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}
		// fmt.Println("token", claims)

		userID, ok := claims["id"]
		// fmt.Println("id", userID)
		if !ok {
			http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
			fmt.Println("Invalid user ID in token")
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
