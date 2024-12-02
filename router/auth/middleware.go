package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"

	"serverless/config"
	"serverless/router/database/crud"
	"serverless/router/schema"
)

func Middleware(db *sql.DB, conf *config.Auth) func(handler http.Handler) http.Handler {
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
				return conf.JWTKey, nil
			})

			if err != nil || !token.Valid {
				w.WriteHeader(http.StatusUnauthorized)
				encoder.Encode(schema.Response{Error: "Invalid token"})
				return
			}

			user := claims.ToUser()
			_, err = crud.GetUser(db, schema.CredentialsFromUser(user))
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				encoder.Encode(schema.Response{Error: "Invalid token"})
				return
			}

			ctx := context.WithValue(r.Context(), "user", user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
