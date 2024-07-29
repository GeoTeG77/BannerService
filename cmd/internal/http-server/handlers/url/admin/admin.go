package admin

import (
	//resp "BannerService/cmd/internal/lib/api/response"
	templates "BannerService/cmd/internal/templates"
	//"log/slog"
	"net/http"
)

func AdminPage(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err := templates.Tmpl.ExecuteTemplate(w, "adminPannel.html", nil)
		if err != nil {
			http.Error(w, "Unable to load template", http.StatusInternalServerError)
		}
}
