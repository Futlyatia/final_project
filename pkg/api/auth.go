package api

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

const jwtSecret = "todo-secret-key"

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		password := os.Getenv("TODO_PASSWORD")

		if password == "" {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(
			cookie.Value,
			func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrTokenSignatureInvalid
				}

				return []byte(jwtSecret), nil
			},
		)

		if err != nil || !token.Valid {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}

		passwordHash, ok := claims["password_hash"].(string)
		if !ok {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}

		currentHash := hashPassword(password)

		if passwordHash != currentHash {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}
