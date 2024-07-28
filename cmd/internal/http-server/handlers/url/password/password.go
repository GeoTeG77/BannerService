package password

import (
	//"BannerService/cmd/internal/lib/logger/sl"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/golang-jwt/jwt/v4"

	"golang.org/x/crypto/argon2"
	//"github.com/go-chi/render"
	//"github.com/go-playground/validator"
)

type CustomClaim struct {
	Name            string `json:"name"`
	IsAdmin         string `json:"is_admin"`
	UseLastRevision string `json:"use_last_revision"`
	jwt.StandardClaims
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type User struct {
	UserName       string `json:"username"`
	Password       string `json:"password"`
	UseLastVersion string `json:"use_last_revision"`
	IsAdmin        string `json:"is_admin"`
}

type PasswordManager interface {
	CheckPassword(username string) (string, string, string, string, error)
	CreateUser(username string, password_hash string, salt string, isAdmin string, useLastRevision string) error
	ChangePassword(username string, password_hash string, salt string) error
}

func CheckPassword(log *slog.Logger, PasswordManager PasswordManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.PasswordManager.CheckPassword"
		log = log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))

		username := r.URL.Query().Get("username")
		password := r.URL.Query().Get("password")

		passwordHash, salt, isAdmin, useLastRevision, err := PasswordManager.CheckPassword(username)
		if err != nil {
			log.Error("wrong username or password")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}

		enterPasHash := hashPassword(password, salt)
		if enterPasHash != passwordHash {
			log.Error("wrong username or password")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}

		claim := &CustomClaim{
			Name:            username,
			IsAdmin:         isAdmin,
			UseLastRevision: useLastRevision,
		}

		tokenString, err := CreateToken(claim)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "banner_token",
			Value:    tokenString,
			Expires:  time.Now().Add(time.Minute * 15),
			HttpOnly: true,
		})
		
	switch isAdmin{
	case "0":
		http.Redirect(w,r, "/user", http.StatusSeeOther )
	case "1":
		http.Redirect(w,r, "/admin",  http.StatusSeeOther )
	}
}
}

func CreateUser(log *slog.Logger, PasswordManager PasswordManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.PasswordManager.CreateUser"
		log = log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		var data User
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println(data)

		salt, err := generateSalt()
		if err != nil {
			log.Info("internal error(500)")
		}
		fmt.Println(salt)
		passwordHash := hashPassword(data.Password, salt)
		err = PasswordManager.CreateUser(data.UserName, passwordHash, salt, data.IsAdmin, data.UseLastVersion)
	
		if err != nil {
			log.Error("wrong username or password")
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}


func generateSalt() (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(salt), nil
}

func hashPassword(password, salt string) string {
	saltBytes, _ := base64.RawStdEncoding.DecodeString(salt)
	hash := argon2.IDKey([]byte(password), saltBytes, 1, 64*1024, 4, 32)
	return base64.RawStdEncoding.EncodeToString(hash)
}

func CreateToken(claim *CustomClaim) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateToken(tokenStr string) (*CustomClaim, error) {
	claims := &CustomClaim{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	return claims, nil
}
