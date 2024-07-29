package banner

import (
	resp "BannerService/cmd/internal/lib/api/response"
	"BannerService/cmd/internal/lib/logger/sl"
	templates "BannerService/cmd/internal/templates"
	"log/slog"
	"net/http"

	//"os/exec"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Content struct {
	Title *string `json:"title"`
	Text  *string `json:"text"`
	URL   *string `json:"url" validate:"required,url"`
}

type CreateRequest struct {
	Tag_ids    []int64  `json:"tag_ids" validate:"required"`
	Feature_id *int64   `json:"feature_id" validate:"required"`
	Content    *Content `json:"content" validate:"required"`
	Is_active  *bool    `json:"is_active" validate:"required"`
}

type UpdateRequest struct {
	Banner_id  *int64   `json:"banner_id" validate:"required"`
	Tag_ids    []int64  `json:"tag_ids" validate:"required"`
	Feature_id *int64   `json:"feature_id,omitempty"`
	Content    *Content `json:"content,omitempty"`
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
	CreateBanner(req CreateRequest) (int64, error)
	UpdateBanner(req UpdateRequest) error
	GetBanner(req GetBannerRequest) (Content, error)
	GetBanners(req GetBannersRequest) ([]GetBannersResponce, error)
	DeleteBanner(id int64) (Response, error)
}

func SetBannerPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := templates.Tmpl.ExecuteTemplate(w, "setBanner.html", nil)
	if err != nil {
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
	}

	//JSON и редерект в SET BANNER () POST /banner
}

func SetBanner(log *slog.Logger, BannerImpl Banner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.banner.CreateBanner"

		log = log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))

		var req CreateRequest

		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decide request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {

			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.Error("invalid request"))
			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}
		_, err = Banner.CreateBanner(BannerImpl, req)
		if err != nil {
			log.Error("failed")
		}
	}
}

// / Поправить UpdateBanner исправить через URL параметр и sqlite соответственно!!!!
func UpdateBannerPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := templates.Tmpl.ExecuteTemplate(w, "update.html", nil)
	if err != nil {
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
	}

	//JSON и редерект в SET BANNER () POST /banner
}

func UpdateBanner(log *slog.Logger, BannerImpl Banner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.banner.UpdateBanner"

		log = log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))

		var req UpdateRequest

		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decide request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {

			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.Error("invalid request"))
			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		//Проверить
		err = Banner.UpdateBanner(BannerImpl, req) // хз
		if err != nil {
			log.Error("failed")
		}
	}
}

// / Поправить UpdateBanner исправить через URL параметр и sqlite соответственно!!!!
func GetBannerPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := templates.Tmpl.ExecuteTemplate(w, "getBanner.html", nil)
	if err != nil {
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
	}

	//JSON и редерект в SET BANNER () POST /banner
}

func GetBanner(log *slog.Logger, BannerImpl Banner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.banner.GetBanner"

		log = log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		var req GetBannerRequest

		tag_id := r.URL.Query().Get("tag_id")
		tag, err := strconv.Atoi(tag_id)
		if err != nil {
			log.Error("bad value")
		}
		tag64 := int64(tag)
		req.Tag_id = &tag64

		feature_id := r.URL.Query().Get("feature_id")
		feature, err := strconv.Atoi(feature_id)
		if err != nil {
			log.Error("bad value")
		}
		feature64 := int64(feature)
		req.Tag_id = &feature64

		res, err := Banner.GetBanner(BannerImpl, req)
		if err != nil {
			log.Error("db error")
		}
		_ = res //gg

	}
}

func GetBannersPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := templates.Tmpl.ExecuteTemplate(w, "getBanners.html", nil)
	if err != nil {
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
	}

	//JSON и редерект в SET BANNER () POST /banner
}

func GetBanners(log *slog.Logger, BannerImpl Banner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.banner.GetBanners"

		log = log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))

		var req GetBannersRequest

		tag_id := r.URL.Query().Get("tag_id")
		tag, err := strconv.Atoi(tag_id)
		if err != nil {
			log.Error("bad value")
		}
		tag64 := int64(tag)
		req.Tag_id = &tag64

		feature_id := r.URL.Query().Get("feature_id")
		feature, err := strconv.Atoi(feature_id)
		if err != nil {
			log.Error("bad value")
		}
		feature64 := int64(feature)
		req.Tag_id = &feature64

		limitt := r.URL.Query().Get("limit")
		limit, err := strconv.Atoi(limitt)
		if err != nil {
			log.Error("bad value")
		}
		limit64 := int64(limit)
		req.Limit = &limit64

		offsett := r.URL.Query().Get("offset")
		offset, err := strconv.Atoi(offsett)
		if err != nil {
			log.Error("bad value")
		}
		offset64 := int64(offset)
		req.Offset = &offset64

		res, err := Banner.GetBanners(BannerImpl, req)
		if err != nil {
			log.Error("db error")
		}
		_ = res

	}
}

func DeleteBannerPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := templates.Tmpl.ExecuteTemplate(w, "deleteBanner.html", nil)
	if err != nil {
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
	}

	//JSON и редерект в SET BANNER () POST /banner
}

func DeleteBanner(log *slog.Logger, BannerImpl Banner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.banner.DeleteBanner"
		log = log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		id := chi.URLParam(r, "id")
		limit, err := strconv.Atoi(id)
		if err != nil {
			log.Error("bad value")
		}
		id64 := int64(limit)

		res, err := Banner.DeleteBanner(BannerImpl, id64)
		if err != nil {
			log.Error("db error")
			_ = res

		}
	}
}
