package auth

import (
	"context"
	"encoding/json"
	"github.com/golang-jwt/jwt"
	"net/http"
	"serverless/config"
	"serverless/router/schema"
	"strings"
)

func Middleware(conf *config.Config) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			encoder := json.NewEncoder(w)

			tokenStr := r.Header.Get("Authorization")
			if tokenStr == "" {
				w.WriteHeader(http.StatusUnauthorized)
				encoder.Encode(schema.Response{Error: "Authorization header missing"})
				return
			}

			tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

			claims := &schema.Claims{}

			token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				return conf.Auth.JWTKey, nil
			})

			if err != nil || !token.Valid {
				w.WriteHeader(http.StatusUnauthorized)
				encoder.Encode(schema.Response{Error: "Invalid token"})
				return
			}

			ctx := context.WithValue(r.Context(), "username", claims.Username)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
