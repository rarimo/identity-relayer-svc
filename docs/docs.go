package docs

import (
	"embed"
	httptemplate "html/template"
	"net/http"
)

const (
	apiFile   = "/static/service.swagger.json"
	indexFile = "template/index.tpl"
)

//go:embed static
var Docs embed.FS

//go:embed template
var template embed.FS

func RegisterOpenAPIService(appName string, rtr *http.ServeMux) {
	rtr.Handle(apiFile, http.FileServer(http.FS(Docs)))
	rtr.HandleFunc("/swagger", handler(appName))
}

// handler returns an http handler that servers OpenAPI console for an OpenAPI spec at specURL.
func handler(title string) http.HandlerFunc {
	t, _ := httptemplate.ParseFS(template, indexFile)

	return func(w http.ResponseWriter, req *http.Request) {
		t.Execute(w, struct {
			Title string
			URL   string
		}{
			title,
			apiFile,
		})
	}
}
