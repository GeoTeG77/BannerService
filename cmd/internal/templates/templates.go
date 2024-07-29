package templates

import (
	"html/template"
	"path/filepath"
)

var Tmpl *template.Template

func LoadTemplates() {
	Tmpl = template.Must(template.ParseGlob(filepath.Join("cmd", "internal", "templates", "*.html")))
}
