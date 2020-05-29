package server

{{ if api.auth.enabled }}
import (
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"

	"{{ api.extra.repo }}/go/pkg/dao"
)

func (s Server) SessionStartController(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var userPass struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := getBody(r, &userPass)
	if err != nil {
		sendValidationResponse(w, fmt.Sprintf(`Expected "username" and "password" in body, got: %s`, err))
		return
	}

	filter, err := dao.ParseFilter(fmt.Sprintf("{{ api.auth.username }} = '%s'", userPass.Username))
	if err != nil {
		s.logger.Debugf("Error while parsing {{ api.auth.username }}: %s", err)
		sendValidationResponse(w, `Expected valid "username"`)
		return
	}

	pageInfo := dao.PageInfo{Offset: 0, Limit: 1, Order: `"{{ api.auth.username }} DESC"`}
	result, err := s.dao.{{ api.auth.table|string.capitalize }}GetMany(filter, pageInfo)
	if err != nil {
		sendErrorResponse(w, err)
		return
	}

	if result.Total == 0 {
		sendValidationResponse(w, `Invalid "username" or "password"`)
		return
	}

	user := result.Data[0]
	err := bcrypt.CompareHashAndPassword([]byte(user.C_{{ api.auth.password }}), []byte(userPass.Password))
	if err != nil {
		sendValidationResponse(w, `Invalid "username" or "password"`)
		return
	}

	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.C_{{ api.auth.username }},
		"exp": time.Now().Add(s.sessionLength).Unix(),
	})
	token, err := unsignedToken.SignedString(s.sessionSigningKey)
	if err != nil {
		sendErrorResponse(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name: "au",
		Value: token,
		Expires: time.Now().Add(s.sessionLength)
	})
			
	sendResponse(w, struct{
		Token string `json:"token"`
	}{token})
}
{{ else }}
// Auth not enabled.
{{ end }}
