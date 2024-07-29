package auth

import (
	password "BannerService/cmd/internal/http-server/handlers/url/password"
	templates "BannerService/cmd/internal/templates"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type CreateUser struct {
	UserName       string `json:"username"`
	Password       string `json:"password"`
	UseLastVersion string `json:"use_last_revision"`
	IsAdmin        string `json:"is_admin"`
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("banner_token")
		if err != nil || cookie == nil {
			http.Redirect(w, r, "/", http.StatusUnauthorized)
			return
		}

		stringCookie := cookie.Value
		_, _, err = password.ValidateToken(stringCookie)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := templates.Tmpl.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")

	params := url.Values{}
	params.Add("username", username)
	params.Add("password", password)
	fmt.Println(username, password)
	redirectURL := fmt.Sprintf("/auth?%s", params.Encode())
	fmt.Println(redirectURL)
	http.Redirect(w, r, redirectURL, http.StatusFound)

}

func RegisterPageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := templates.Tmpl.ExecuteTemplate(w, "register.html", nil)
	if err != nil {
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	username := r.FormValue("user_name")
	password := r.FormValue("password")
	use_last_revision := r.FormValue("use_last_revision")
	is_admin := r.FormValue("is_admin")
	redirectURL := "http://localhost:8080/auth"
	data := CreateUser{
		username,
		password,
		use_last_revision,
		is_admin,
	}

	jsonData, err := json.Marshal(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest("POST", redirectURL, bytes.NewBuffer(jsonData))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Error(w, "Failed to create user", resp.StatusCode)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {

	http.SetCookie(w, &http.Cookie{
		Name:   "banner_token",
		Value:  "",
		MaxAge: -1,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
