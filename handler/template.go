package handler

import (
	"sync"
	"html/template"
	"net/http"
	"path/filepath"
	"github.com/stretchr/objx"
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

	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.template.Execute(w, data)
}

func NewTemplate(filename string) *Template {
	return &Template{
		filename: filename,
	}
}
