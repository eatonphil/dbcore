package server

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"

	"{{ api.extra.repo }}/go/pkg/dao"
)

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.logger.Infof("%s %s", r.Method, r.URL.RequestURI())
	w.Header().Set("Content-Type", "application/json")

	{{ if api.auth.enabled }}
	cookie, err := r.Cookie("au")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
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

		expInterface, ok := claims["exp"]
		if !ok {
			return false
		}

		expUnix, ok := expInterface.(int64)
		if !ok {
			return false
		}

		if time.Unix(expUnix, 0).Before(time.Now()) {
			// Session has expired
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

		pageInfo := dao.Pagination{Offset: 0, Limit: 1, Order: `"{{ api.auth.username }} DESC"`}
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
		return
	}
	{{ end }}

	s.router.ServeHTTP(w, r)
}
