package main

import (
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

type templateHandler struct {
	once     sync.Once
	fileName string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.fileName)))
	})
	if err := t.templ.Execute(w, nil); err != nil {
		log.Fatal("ServeHTTP:", err)
	}
}

func main() {
	r := newRoom()

	// ルート
	http.Handle("/", &templateHandler{fileName: "chat.html"})
	http.Handle("/room", r)

	// チャットルームを開始
	go r.run()

	// webサーバを起動
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
