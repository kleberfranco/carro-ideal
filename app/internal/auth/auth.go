package auth

import (
	"net/http"
	"strconv"
)

const cookieName = "carro_session"

func SetUserSession(w http.ResponseWriter, userID int64) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    strconv.FormatInt(userID, 10),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func ClearUserSession(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func GetUserID(r *http.Request) (int64, bool) {
	c, err := r.Cookie(cookieName)
	if err != nil || c.Value == "" {
		return 0, false
	}

	id, err := strconv.ParseInt(c.Value, 10, 64)
	if err != nil {
		return 0, false
	}

	return id, true
}
