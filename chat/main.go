package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "application address")
	flag.Parse()

	// Gomniauth
	// random value for secret
	gomniauth.SetSecurityKey("vshimoda")
	// oauth2
	gomniauth.WithProviders(
		facebook.New(
			"318772611789665",
			"620c0b2ad1f21c7a1ac60540e02f8cbc",
			"http://localhost:8080/auth/callback/facebook",
		),
		github.New(
			"a64d23d1fb23faebd34a",
			"64249bfa218997b205076c9021b7bda9aca33770",
			"http://localhost:8080/auth/callback/github",
		),
		google.New(
			"577006250786-e2ina0bireqdu9g9v4qhdaa7hqqassu9.apps.googleusercontent.com",
			"Xj09MUsqQ4Rl2D8L3xR3_qos",
			"http://localhost:8080/auth/callback/google",
		),
	)

	r := newRoom()
	// room
	//http.Handle("/", &templateHandler{
	//	filename: "chat.html",
	//})
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	// start room
	go r.run()
	// Run Web Server
	log.Println("Web server running. Port:", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
