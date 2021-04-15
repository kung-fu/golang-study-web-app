package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
		// 未認証
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		// 何らかのエラー
		panic(err.Error())
	} else {
		// 成功
		h.next.ServeHTTP(w, r)
	}
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

// loginHandler はサードパーティへのログイン処理を受け持ちます
// パスの形式： /auth/{action}/{provider}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	segments := strings.Split(r.URL.Path, "/")
	if len(segments) < 3 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "URLが不正です")
	}

	action := segments[2]
	provider := segments[3]
	switch action {
	case "login":
		log.Println("TODO: ログイン処理", provider)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "アクション%sには非対応です", action)
	}
}
