package main

import (
	"flag"
	"github.com/kung-fu/golang-study-web-app/trace"
	"log"
	"net/http"
	"os"
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
	if err := t.templ.Execute(w, r); err != nil {
		log.Fatal("ServeHTTP:", err)
	}
}

func main() {
	log.Print("start chat server..")

	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse()

	r := newRoom()
	r.tracer = trace.New(os.Stdout)

	// ルート
	http.Handle("/", &templateHandler{fileName: "chat.html"})
	http.Handle("/room", r)

	// チャットルームを開始
	go r.run()

	// webサーバを起動
	log.Printf("Webサーバを開始します（ポート: %s）", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
