package handlers

import (
	"net/http"

	"pitch.ideas/internal/views"
)

func setSessionCookies(w http.ResponseWriter, sessionId string, username string) {
	maxAge := 60 * 60 * 24 * 7

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionId,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   maxAge,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "logged_in",
		Value:    "True",
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   maxAge,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "username",
		Value:    username,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   maxAge,
	})
}

func deleteSessionCookies(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "logged_in",
		Value:    "",
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "username",
		Value:    "",
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}

func LoginPage(renderer *views.Renderer) http.HandlerFunc {
	return  func(w http.ResponseWriter, r *http.Request) {
		renderer.Render(w, "login.html", "")
	}
}

func Login(w http.ResponseWriter, r *http.Request) {

}


func RegisterPage(renderer *views.Renderer) http.HandlerFunc {
	return  func(w http.ResponseWriter, r *http.Request) {
		renderer.Render(w, "register.html", "")
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
}

func LogoutPage(renderer *views.Renderer) http.HandlerFunc {
	return  func(w http.ResponseWriter, r *http.Request) {
		renderer.Render(w, "logout.html", "")
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
}

func AuthStatus(w http.ResponseWriter, r *http.Request) {
}