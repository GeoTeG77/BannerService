package role

import (
	password "BannerService/cmd/internal/http-server/handlers/url/password"
	"net/http"
)

func IsAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("banner_token")
		if err != nil || cookie == nil {
			http.Redirect(w, r, "/", http.StatusUnauthorized)
			return
		}

		stringCookie := cookie.Value
		_,isAdmin, err := password.ValidateToken(stringCookie)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusUnauthorized)
			return
		}
		if isAdmin =="0"{
			http.Redirect(w, r, "/user", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
		
		if isAdmin =="1"{
			http.Redirect(w, r, "/admin", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
