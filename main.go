package main

import (
	"net/http"
	"log"
	"chat/handler"
	"flag"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		panic(err.Error())
	} else {
		h.next.ServeHTTP(w, r)
	}
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{
		next: handler,
	}
}
func main() {
	var addr = flag.String("addr", ":8080", "The address of the application")
	flag.Parse()
	r := handler.NewRoom()

	http.Handle("/", handler.NewTemplate("chat.html"))

	http.Handle("/room", r)

	log.Println("Starting web server on ", *addr)

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

	println("Started!")
}

