package views

import (
	"html/template"
	"net/http"
)

type Views struct {
	*template.Template
}

var tpl = New("templates")

func New(path string) *Views {
	t := template.New(path)
	template.Must(t.ParseGlob(path + "/*.html"))
	return &Views{t}
}

func (v *Views) Render(w http.ResponseWriter, name string, data interface{}) {
	err := v.ExecuteTemplate(w, name+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Render(w http.ResponseWriter, name string, data interface{}) {
	tpl.Render(w, name, data)
}

func NotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	tpl.Render(w, "404", nil)
}

func Error(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	tpl.Render(w, "500", err)
}
