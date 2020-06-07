package server

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"

	"{{ api.extra.repo }}/{{ out_dir }}/pkg/dao"
)

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.logger.Infof("%s %s", r.Method, r.URL.RequestURI())
	w.Header().Set("Content-Type", "application/json")

	for _, allowed := range s.allowedOrigins {
		origin := strings.ToLower(r.Header.Get("origin"))
		if strings.ToLower(allowed) == origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, Origin")
		}
	}

	{{ if api.auth.enabled }}
	if r.URL.Path == "/{{ api.router_prefix }}session/start" {
		s.router.ServeHTTP(w, r)
		return
	}

	cookie, err := r.Cookie("au")
	if err != nil {
		// Fall back to header check
		cookie = &http.Cookie{}
	}

	token := cookie.Value
	if token == "" {
		authHeader := r.Header.Get("Authorization")
		if len(authHeader) > len("bearer ") &&
			strings.ToLower(authHeader[:len("bearer ")]) == "bearer " {
			token = authHeader[len("bearer "):]
		}
	}

	authorized := func() bool {
		t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
			}

			return []byte(s.secret), nil
		})
		if err != nil {
			s.logger.Debugf("Error parsing JWT:", err)
			return false
		}

		claims, ok := t.Claims.(jwt.MapClaims)
		if !(ok && t.Valid) {
			return false
		}

		if err := claims.Valid(); err != nil {
			s.logger.Debugf("Invalid claims: %s", err)
			return false
		}

		usernameInterface, ok := claims["username"]
		if !ok {
			return false
		}

		username, ok := usernameInterface.(string)
		if !ok {
			return false
		}

		filter, err := dao.ParseFilter(fmt.Sprintf("{{ api.auth.username }} = '%s'", username))
		if err != nil {
			s.logger.Debugf("Error while parsing {{ api.auth.username }}: %s", err)
			return false
		}

		pageInfo := dao.Pagination{Offset: 0, Limit: 1, Order: `"{{ api.auth.username }}" DESC`}
		result, err := s.dao.{{ api.auth.table|string.capitalize }}GetMany(filter, pageInfo)
		if err != nil {
			s.logger.Debugf("Error retrieving user: %s", err)
			return false
		}

		if result.Total == 0 {
			return false
		}

		return true
	}()
	if !authorized {
		w.WriteHeader(http.StatusUnauthorized)
		sendResponse(w, struct{
			Error string `json:"error"`
		}{"Unauthorized"})
		return
	}
	{{ end }}

	s.router.ServeHTTP(w, r)
}
