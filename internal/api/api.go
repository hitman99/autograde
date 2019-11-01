package api

import (
	"encoding/base64"
	"github.com/hitman99/autograde/internal/config"
	"net/http"
	"strings"
)

type authMiddleware struct {
	adminToken string
}

func NewAuthMiddleware(adminToken string) *authMiddleware {
	return &authMiddleware{
		adminToken: adminToken,
	}
}

func (a *authMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		parts := strings.Split(strings.TrimSpace(token), " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				http.Error(w, "Access Denied", http.StatusUnauthorized)
				return
			} else {
				if string(decoded) != config.GetConfig().AdminToken {
					http.Error(w, "Access Denied", http.StatusUnauthorized)
					return
				}
				next.ServeHTTP(w, r)
			}
		} else {
			http.Error(w, "Access Denied", http.StatusUnauthorized)
		}
	})
}
