package tag

import (
	resp "BannerService/cmd/internal/lib/api/response"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5/middleware"
)

type Tag interface {
	CreateTag(tag_name string) (int64, error)
	UpdateTagName(tag_id int64, tag_name string) error
	DeleteTag(tag_id int64) error
}

func CreateTag(log *slog.Logger, TagImpl Tag) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.banner.CreateTag"
		log = log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		name := r.URL.Query().Get("id")

		//fix this
		f:= resp.BadRequest()
		_=f
		res, err := TagImpl.CreateTag(name)
		if err != nil {
			log.Error("db error")
			_ = res
		}
	}
}

func UpdateTagName(log *slog.Logger, TagImpl Tag) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.banner.UpdateTagName"
		log = log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		tag_id := r.URL.Query().Get("id")
		name := r.URL.Query().Get("name")
		tag_id32, err := strconv.Atoi(tag_id)
		if err != nil {
			log.Error("db error")
		}
		tag_id64 := int64(tag_id32)
		err = Tag.UpdateTagName(TagImpl, tag_id64, name)
		if err != nil {
			log.Error("db error")
		}
	}
}

func DeleteTag(log *slog.Logger, TagImpl Tag) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.banner.CreateTag"
		log = log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		tag_id := r.URL.Query().Get("id")
		tag_id32, err := strconv.Atoi(tag_id)
		if err != nil {
			log.Error("db error")
		}
		tag_id64 := int64(tag_id32)
		err = Tag.DeleteTag(TagImpl, tag_id64)
		if err != nil {
			log.Error("db error")
		}
	}
}
