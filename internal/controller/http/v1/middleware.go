package v1

import (
	"net/http"
	"strings"
	"vk-film-library/internal/service"
	"vk-film-library/pkg/logger"
)

const (
	userRoleHeader = "userRole"
)

type AuthMiddleware struct {
	authService service.Auth
	log         *logger.Logger
}

func (m *AuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		token, ok := getToken(req)
		if !ok {
			m.log.Errorf("AuthMiddleware RequireAuth: getToken %v", ErrInvalidAuthHeader)
			http.Error(w, ErrInvalidAuthHeader.Error(), http.StatusUnauthorized)
			return
		}

		userRole, err := m.authService.ParseToken(token)
		if err != nil {
			m.log.Errorf("AuthMiddleware RequireAuth: authService.ParseToken %v", err)
			http.Error(w, ErrCannotParseToken.Error(), http.StatusUnauthorized)
			return
		}

		w.Header().Set(userRoleHeader, userRole)
		next.ServeHTTP(w, req)
	})
}

func getToken(req *http.Request) (string, bool) {
	const prefix = "Bearer "

	header := req.Header.Get("Authorization")
	if header == "" {
		return "", false
	}

	if len(header) > len(prefix) && strings.EqualFold(header[:len(prefix)], prefix) {
		return header[len(prefix):], true
	}

	return "", false
}
