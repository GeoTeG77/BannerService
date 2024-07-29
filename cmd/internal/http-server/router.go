package httpserver

/*import (
	"BannerService/cmd/internal/http-server/handlers/url/banner"
	"BannerService/cmd/internal/http-server/handlers/url/feature"
	"BannerService/cmd/internal/http-server/handlers/url/password"
	"BannerService/cmd/internal/http-server/handlers/url/tag"
	"BannerService/cmd/internal/http-server/middleware/auth"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	//v1 "BannerService/internal/http-server/handlers/v1"
	//v2 "BannerService/internal/http-server/handlers/v2"
)

func Router() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Get("/", auth.LoginPageHandler)
	router.Post("/login", auth.LoginHandler)
	router.Get("/register", auth.RegisterPageHandler)
	router.Post("/register", auth.RegisterHandler)
	router.Post("/auth", password.CreateUser(log, storage))
	router.Get("/auth", password.CheckPassword(log, storage))
	router.Get("/logout", auth.Logout)

	router.Route("/", func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Use(middleware.SetHeader("Content-Type", "application/json"))
		r.Use(middleware.Recoverer)
		r.Use(middleware.URLFormat)

		r.Post("/banner", banner.SetBanner(log, storage))
		r.Post("/tag", tag.CreateTag(log, storage))
		r.Post("/feature", feature.CreateFeature(log, storage))

		r.Patch("/banner/{id}", banner.UpdateBanner(log, storage))
		r.Patch("/feature/{id}/{name}", feature.UpdateFeatureName(log, storage))
		r.Patch("/tag/{id}/{name}", tag.UpdateTagName(log, storage))

		r.Get("/user_banner", banner.GetBanner(log, storage))
		r.Get("/banner", banner.GetBanners(log, storage))

		r.Delete("/banner/{id}", banner.DeleteBanner(log, storage))
		r.Delete("/tag/{id}", tag.DeleteTag(log, storage))
		r.Delete("/feature/{id}", feature.DeleteFeature(log, storage))
	})

	return router
}
*/