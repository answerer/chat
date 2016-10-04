package main

import (
	"net/http"
	"log"
	"chat/handler"
	"flag"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/signature"
)

func main() {
	var addr = flag.String("addr", ":8080", "The address of the application")
	flag.Parse()
	gomniauth.SetSecurityKey(signature.RandomKey(64))
	gomniauth.WithProviders(
		facebook.New("1626414917656337", "0fc55c87caa6b4b76b5a61c4dab2e850",
			"http://localhost:8080/auth/callback/facebook"),
		github.New("key", "secret",
			"http://localhost:8080/auth/callback/github"),
		google.New("54b18c5e912d687efe15", "dad519b463c2f7131cb4f552afe2351cae659406",
			"http://localhost:8080/auth/callback/google"),
	)
	r := handler.NewRoom()

	http.HandleFunc("/auth/", handler.LoginHandle)
	http.Handle("/chat", handler.MustAuth(handler.NewTemplate("chat.html")))
	http.Handle("/login", handler.NewTemplate("login.html"))

	http.Handle("/room", r)

	log.Println("Starting web server on ", *addr)

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

	println("Started!")
}

