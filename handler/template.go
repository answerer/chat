package handler

import (
	"sync"
	"html/template"
	"net/http"
	"path/filepath"
)

type Template struct {
	once     sync.Once
	filename string
	template *template.Template
}

func (t *Template) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.template = template.Must(template.ParseFiles(filepath.Join("/app/src/chat/templates", t.filename)))
	})

	t.template.Execute(w, r)
}

func NewTemplate(filename string) *Template {
	return &Template{
		filename: filename,
	}
}
