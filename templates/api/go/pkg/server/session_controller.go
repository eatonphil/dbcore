package server

{{ if api.auth.enabled }}
import (
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"

	"{{ api.extra.repo }}/{{ out_dir }}/pkg/dao"
)

func (s Server) getSessionUsername(r *http.Request) string {
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

	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(s.auth.secret), nil
	})
	if err != nil {
		s.logger.Debugf("Error parsing JWT: %s", err)
		return ""
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !(ok && t.Valid) {
		return ""
	}

	if err := claims.Valid(); err != nil {
		s.logger.Debugf("Invalid claims: %s", err)
		return ""
	}

	usernameInterface, ok := claims["username"]
	if !ok {
		return ""
	}

	username, ok := usernameInterface.(string)
	if !ok {
		return ""
	}

	return username
}

func (s Server) SessionStartController(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var userPass struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := getBody(r, &userPass)
	if err != nil {
		sendValidationErrorResponse(w, fmt.Sprintf("Expected username and password in body, got: %s", err))
		return
	}

	if userPass.Username == "" || userPass.Password == "" {
		sendValidationErrorResponse(w, "Expected username and password in body")
		return
	}

	q := fmt.Sprintf("{{ api.auth.username }} = '%s'", userPass.Username)
	filter, err := dao.ParseFilter(q)
	if err != nil {
		s.logger.Debugf("Error while parsing {{ api.auth.username }}: %s", err)
		sendValidationErrorResponse(w, "Expected valid username")
		return
	}

	pageInfo := dao.Pagination{Offset: 0, Limit: 1, Order: `"{{ api.auth.username }}" DESC`}
	result, err := s.dao.{{ api.auth.table|dbcore_capitalize }}GetMany(filter, pageInfo, "", nil)
	if err != nil {
		sendErrorResponse(w, err)
		return
	}

	if result.Total == 0 {
		sendValidationErrorResponse(w, "Invalid username or password")
		return
	}

	user := result.Data[0]
	err = bcrypt.CompareHashAndPassword([]byte(user.C_{{ api.auth.password }}), []byte(userPass.Password))
	if err != nil {
		sendValidationErrorResponse(w, "Invalid username or password")
		return
	}

	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.C_{{ api.auth.username }},
		"user_id": user.Id(),
		"exp": time.Now().Add(s.sessionDuration).Unix(),
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
	})
	token, err := unsignedToken.SignedString([]byte(s.auth.secret))
	if err != nil {
		s.logger.Debugf("Error signing string: %s", err)
		sendErrorResponse(w, fmt.Errorf("Internal server error"))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name: "au",
		Value: token,
		Expires: time.Now().Add(s.sessionDuration),
		Path: "/",
	})
			
	sendResponse(w, struct{
		Token string `json:"token"`
	}{token})
}

func (s Server) SessionStopController(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.SetCookie(w, &http.Cookie{
		Name: "au",
		Value: "",
		Expires: time.Now().Add(-1 * s.sessionDuration),
		Path: "/",
	})

	sendResponse(w, struct{
		Token string `json:"token"`
	}{""})
}
{{ else }}
// Auth not enabled.
{{ end }}
