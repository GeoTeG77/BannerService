package banner

import (
	resp "BannerService/cmd/internal/lib/api/response"
	"fmt"

	//"BannerService/cmd/internal/lib/logger/sl"
	templates "BannerService/cmd/internal/templates"
	"log/slog"
	"net/http"
	"strings"

	//"os/exec"
	"strconv"

	"github.com/go-chi/chi/v5/middleware"
	//"github.com/go-chi/render"
	//"github.com/go-playground/validator/v10"
)

type Content struct {
	Title *string `json:"title"`
	Text  *string `json:"text"`
	URL   *string `json:"url" validate:"required,url"`
}

type CreateRequest struct {
	Tag_ids    []int64  `json:"tag_ids" validate:"required"`
	Feature_id *int64   `json:"feature_id" validate:"required"`
	Content    *Content `json:"content"`
	Is_active  *bool    `json:"is_active" validate:"required"`
}

type GetBannerRequest struct {
	Tag_id           *int64
	Feature_id       *int64
	Use_last_version *bool
}

type GetBannersRequest struct {
	Tag_id     *int64
	Feature_id *int64
	Limit      *int64
	Offset     *int64
}

type GetBannersResponce struct {
	Banner_id  *int64   `json:"banner_id" validate:"required"`
	Tag_ids    []int64  `json:"tag_ids" validate:"required"`
	Feature_id *int64   `json:"feature_id,omitempty"`
	Content    *Content `json:"content,omitempty"`
	Is_active  *bool    `json:"is_active" validate:"required"`
	Created_at *string  `json:"created_at"`
	Updated_at *string  `json:"upadated_at"`
}

type Response struct {
	resp.Response
	Status    *string `json:"status"`
	Error     *string `json:"error,omitempty"`
	Banner_id *string `json:"alias,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name==ContentSaver

type Banner interface {
	CreateBanner(int64, []int64, string, string, string, string) (int64, error)
	UpdateBanner(int64, []int64, int64, string, string, string, string) error
	GetBanner(int64, int64) (string, string, string, string, error)
	GetBanners(int64, int64, int64, int64) ([]GetBannersResponce, error)
	DeleteBanner(id int64) (error)
}

func CreateBanner(log *slog.Logger, BannerImpl Banner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		const op = "handlers.url.banner.CreateBanner"
		log = log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))

		switch r.Method {
		case http.MethodGet:

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			err := templates.Tmpl.ExecuteTemplate(w, "setBanner.html", nil)
			if err != nil {
				http.Error(w, "Unable to load template", http.StatusInternalServerError)
			}

		case http.MethodPost:

			if err := r.ParseForm(); err != nil {
				http.Error(w, "Invalid form data", http.StatusBadRequest)
				return
			}

			feature := r.FormValue("feature_id")
			feature_id, err := strconv.ParseInt(feature, 10, 64)
			if err != nil {
				http.Error(w, "Invalid tag_id", http.StatusBadRequest)
			}

			tags := r.FormValue("tag_ids")
			tag_ids := make([]int64, 0, 8)
			tags = strings.ReplaceAll(tags, " ", "")
			tagsString := strings.Split(tags, ",")

			for _, tag := range tagsString {
				tagID, err := strconv.ParseInt(tag, 10, 64)
				if err != nil {
					http.Error(w, "Invalid tag_id", http.StatusBadRequest)
				}
				tag_ids = append(tag_ids, tagID)
			}

			title := r.FormValue("title")
			text := r.FormValue("text")
			url := r.FormValue("url")
			is_active := r.FormValue("is_active")

			_, err = Banner.CreateBanner(BannerImpl, feature_id, tag_ids, title, text, url, is_active)
			if err != nil {
				log.Error("failed")
			}
		}
	}
}

// поправить TAGBannerUpdate

func UpdateBanner(log *slog.Logger, BannerImpl Banner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.banner.UpdateBanner"

		switch r.Method {
		case http.MethodGet:

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			err := templates.Tmpl.ExecuteTemplate(w, "update.html", nil)
			if err != nil {
				http.Error(w, "Unable to load template", http.StatusInternalServerError)
			}

		case http.MethodPost:

			log = log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))

			if err := r.ParseForm(); err != nil {
				http.Error(w, "Invalid form data", http.StatusBadRequest)
				return
			}

			banner := r.FormValue("banner_id")
			banner_id, err := strconv.ParseInt(banner, 10, 64)
			if err != nil {
				http.Error(w, "Invalid tag_id", http.StatusBadRequest)
			}

			tags := r.FormValue("tag_ids")
			tag_ids := make([]int64, 0, 8)
			if tags != "" {
				tags = strings.ReplaceAll(tags, " ", "")
				tagsString := strings.Split(tags, ",")
				for _, tag := range tagsString {
					tagID, err := strconv.ParseInt(tag, 10, 64)
					if err != nil {
						http.Error(w, "Invalid tag_id", http.StatusBadRequest)
					}
					tag_ids = append(tag_ids, tagID)
				}
			}

			feature := r.FormValue("feature_id")
			var feature_id int64
			if feature != "" {
				feature_id, err = strconv.ParseInt(feature, 10, 64)
				if err != nil {
					http.Error(w, "Invalid tag_id", http.StatusBadRequest)
				}
			}

			title := r.FormValue("title")
			text := r.FormValue("text")
			url := r.FormValue("url")
			is_active := r.FormValue("is_active")

			//Проверить
			err = Banner.UpdateBanner(BannerImpl, banner_id, tag_ids, feature_id, title, text, url, is_active) // хз
			if err != nil {
				log.Error("failed")
			}
		}
	}
}

// доделать redirect на форму ответа
func GetBanner(log *slog.Logger, BannerImpl Banner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.banner.GetBanner"
		log = log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))

		switch r.Method {
		case http.MethodGet:

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			err := templates.Tmpl.ExecuteTemplate(w, "getBanner.html", nil)
			if err != nil {
				http.Error(w, "Unable to load template", http.StatusInternalServerError)
			}

		case http.MethodPost:

			if err := r.ParseForm(); err != nil {
				http.Error(w, "Invalid form data", http.StatusBadRequest)
				return
			}

			feature := r.FormValue("feature_id")
			feature_id, err := strconv.ParseInt(feature, 10, 64)
			if err != nil {
				http.Error(w, "Invalid feature_id", http.StatusBadRequest)
			}

			tag := r.FormValue("tag_id")
			tag_id, err := strconv.ParseInt(tag, 10, 64)
			if err != nil {
				http.Error(w, "Invalid tag_id", http.StatusBadRequest)
			}

			text, title, url, is_active, err := Banner.GetBanner(BannerImpl, feature_id, tag_id)
			if err != nil {
				log.Error("db error")
			}
			switch {
			case is_active == "1":
				//redirect for all
				fmt.Println(text, title, url) //Cделать выдачу в форму выдачи!
			case is_active == "0":
				//redirect to admin
			}
		}
	}

}

func GetBanners(log *slog.Logger, BannerImpl Banner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.banner.GetBanners"

		log = log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))

		switch r.Method {
		case http.MethodGet:

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			err := templates.Tmpl.ExecuteTemplate(w, "getBanners.html", nil)
			if err != nil {
				http.Error(w, "Unable to load template", http.StatusInternalServerError)
			}

		case http.MethodPost:

			if err := r.ParseForm(); err != nil {
				http.Error(w, "Invalid form data", http.StatusBadRequest)
				return
			}

			feature := r.FormValue("feature_id")
			tag := r.FormValue("tag_id")

			if feature == "" && tag == "" {
				http.Error(w, "Unable to load template", http.StatusBadRequest)
				return
			}

			var feature_id, tag_id, limit, offset int64
			var err error

			if feature != "" {
				feature_id, err = strconv.ParseInt(feature, 10, 64)
				if err != nil {
					http.Error(w, "Invalid feature_id", http.StatusBadRequest)
					return
				}
			} else {
				feature_id = 0
			}

			if tag != "" {
				tag_id, err = strconv.ParseInt(tag, 10, 64)
				if err != nil {
					http.Error(w, "Invalid tag_id", http.StatusBadRequest)
				}
			} else {
				tag_id = 0
			}

			limitStr := r.FormValue("limit")
			if limitStr == "" {
				limit = 10
			}

			offsetStr := r.FormValue("offset")
			if offsetStr == "" {
				offset = 0
			}

			res := make([]GetBannersResponce, 0, limit)
			_= res
			res, err = Banner.GetBanners(BannerImpl, feature_id, tag_id, limit, offset)
			if err != nil {
				log.Error("db error")
			}
			
			fmt.Println(res) //форма для вывода кучи баннеров
		}

	}
}

func DeleteBanner(log *slog.Logger, BannerImpl Banner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.banner.DeleteBanner"
		log = log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		switch r.Method {
		case http.MethodGet:

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			err := templates.Tmpl.ExecuteTemplate(w, "deleteBanner.html", nil)
			if err != nil {
				http.Error(w, "Unable to load template", http.StatusInternalServerError)
			}

		case http.MethodPost:

			if err := r.ParseForm(); err != nil {
				http.Error(w, "Invalid form data", http.StatusBadRequest)
				return
			}

			banner := r.FormValue("banner_id")

			var banner_id int64
			var err error

			if banner != "" {
				banner_id, err = strconv.ParseInt(banner, 10, 64)
				if err != nil {
					http.Error(w, "Invalid feature_id", http.StatusBadRequest)
					return
				}
			} else {
				banner_id = 0
			}

			err = Banner.DeleteBanner(BannerImpl, banner_id)
			if err != nil {
				log.Error("db error")
			}
		}
	}
}
