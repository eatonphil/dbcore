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
		sendValidationErrorResponse(w, fmt.Sprintf(`Expected "username" and "password" in body, got: %s`, err))
		return
	}

	filter, err := dao.ParseFilter(fmt.Sprintf("{{ api.auth.username }} = '%s'", userPass.Username))
	if err != nil {
		s.logger.Debugf("Error while parsing {{ api.auth.username }}: %s", err)
		sendValidationErrorResponse(w, `Expected valid "username"`)
		return
	}

	pageInfo := dao.Pagination{Offset: 0, Limit: 1, Order: `"{{ api.auth.username }} DESC"`}
	result, err := s.dao.{{ api.auth.table|string.capitalize }}GetMany(filter, pageInfo)
	if err != nil {
		sendErrorResponse(w, err)
		return
	}

	if result.Total == 0 {
		sendValidationErrorResponse(w, `Invalid "username" or "password"`)
		return
	}

	user := result.Data[0]
	err = bcrypt.CompareHashAndPassword([]byte(user.C_{{ api.auth.password }}), []byte(userPass.Password))
	if err != nil {
		sendValidationErrorResponse(w, `Invalid "username" or "password"`)
		return
	}

	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.C_{{ api.auth.username }},
		"exp": time.Now().Add(s.sessionDuration).Unix(),
	})
	token, err := unsignedToken.SignedString(s.secret)
	if err != nil {
		sendErrorResponse(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name: "au",
		Value: token,
		Expires: time.Now().Add(s.sessionDuration),
	})
			
	sendResponse(w, struct{
		Token string `json:"token"`
	}{token})
}
{{ else }}
// Auth not enabled.
{{ end }}
