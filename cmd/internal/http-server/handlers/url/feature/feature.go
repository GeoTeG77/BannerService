package feature

import (
	resp "BannerService/cmd/internal/lib/api/response"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5/middleware"
)

type Feature interface {
	CreateFeature(feature_name string) (int64, error)
	UpdateFeatureName(feature_id int64, feature_name string) error
	DeleteFeature(feature_id int64) error
}

func CreateFeature(log *slog.Logger, FeatureImpl Feature) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.banner.CreateFeature"
		log = log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		name := r.URL.Query().Get("id")
		//fix this
		f := resp.BadRequest()
		_ = f
		res, err := Feature.CreateFeature(FeatureImpl, name)
		if err != nil {
			log.Error("db error")
			_ = res
		}
	}
}

func UpdateFeatureName(log *slog.Logger, FeatureImpl Feature) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.banner.UpdateFeatureName"
		log = log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		feature_id := r.URL.Query().Get("id")
		name := r.URL.Query().Get("name")
		feature_id32, err := strconv.Atoi(feature_id)
		if err != nil {
			log.Error("db error")
		}
		feature_id64 := int64(feature_id32)
		err = Feature.UpdateFeatureName(FeatureImpl, feature_id64, name)
		if err != nil {
			log.Error("db error")
		}
	}
}

func DeleteFeature(log *slog.Logger, FeatureImpl Feature) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.banner.CreateFeature"
		log = log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))
		feature_id := r.URL.Query().Get("id")
		feature_id32, err := strconv.Atoi(feature_id)
		if err != nil {
			log.Error("db error")
		}
		feature_id64 := int64(feature_id32)
		err = Feature.DeleteFeature(FeatureImpl, feature_id64)
		if err != nil {
			log.Error("db error")
		}
	}
}
