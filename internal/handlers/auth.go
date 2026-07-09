package handlers

import (
	"net/http"
	"encoding/json"

	"pitch.ideas/internal/database"
	"pitch.ideas/internal/views"
	"pitch.ideas/internal/auth"
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


type AuthRequest struct {
    Username       string `json:"username" validate:"required,min=1,max=100"`
    Password       string `json:"password" validate:"required,min=1,max=100"`
}


func LoginPage(renderer *views.Renderer) http.HandlerFunc {
	return  func(w http.ResponseWriter, r *http.Request) {
		renderer.Render(w, "login.html", "")
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	validUsername := auth.IsValidUsername(req.Username)
	if !validUsername {
		http.Error(w, "Invalid username", http.StatusBadRequest)
		return
	}

	user := database.GetUserByUsername(req.Username)
	if user == nil || !auth.VerifyPassword(req.Password, user.PasswordHash) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}


	sessionId, err := database.CreateSession(user.ID, 7)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	setSessionCookies(w, sessionId, user.Username)
	w.WriteHeader(http.StatusOK)
}

func RegisterPage(renderer *views.Renderer) http.HandlerFunc {
	return  func(w http.ResponseWriter, r *http.Request) {
		renderer.Render(w, "register.html", "")
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	validUsername := auth.IsValidUsername(req.Username)
	if !validUsername {
		http.Error(w, "Invalid username", http.StatusBadRequest)
		return
	}

	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	user, err := database.CreateUser(req.Username, passwordHash)
	if err != nil {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	sessionId, err := database.CreateSession(user.ID, 7)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	setSessionCookies(w, sessionId, user.Username)
	w.WriteHeader(http.StatusOK)
}

func LogoutPage(renderer *views.Renderer) http.HandlerFunc {
	return  func(w http.ResponseWriter, r *http.Request) {
		renderer.Render(w, "logout.html", "")
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	database.DeleteSession(cookie.Value)
	deleteSessionCookies(w)

	w.WriteHeader(http.StatusOK)
}

type UserResponse struct {
	Id uint `json:"id"`
	Username string `json:"username"`
}

type AuthStatusResponse struct {
    LoggedIn bool `json:"logged_in"`
	User UserResponse `json:"user"`
}

func AuthStatus(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user := database.GetUserBySession(cookie.Value)
	
	resp := AuthStatusResponse{
		LoggedIn: user != nil,
	}

	if user != nil {
		resp.User = UserResponse{
			Id: user.ID,
			Username: user.Username,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}


    // session_id = request.cookies.get(SESSION_COOKIE_NAME)
    // if not session_id:
    //     return jsonify({"logged_in": False}), 200

    // user = db.get_user_by_session(session_id=session_id)
    // if not user:
    //     return jsonify({"logged_in": False}), 200
    

    // return jsonify({
    //     "logged_in": True,
    //     "user": {
    //         "id": user.id,
    //         "username": user.username
    //     }
    // })