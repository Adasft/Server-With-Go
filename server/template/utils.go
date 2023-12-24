package template

import (
	"html/template"
	"net/http"
	"path/filepath"
)

const (
	viewsDirectoryName      = "views/"
	layoutsDirectoryName    = viewsDirectoryName + "layouts/"
	layoutFileExtensionType = ".html"
)

var funcMap = template.FuncMap{
	"safeHTML": func(s string) template.HTML {
		return template.HTML(s)
	},
}

func Render(w http.ResponseWriter, data interface{}, path ...string) (*template.Template, error) {
	tmpl := template.New(filepath.Base(path[0])).Funcs(funcMap)
	tmpl, err := tmpl.ParseFiles(path...)
	if err != nil {
		return nil, err
	}

	err = tmpl.ExecuteTemplate(w, filepath.Base(path[0]), data)
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}
func GetLayout(filename string) string {
	return layoutsDirectoryName + filename + layoutFileExtensionType
}

func GetView(filename string) string {
	return viewsDirectoryName + filename + layoutFileExtensionType
}
