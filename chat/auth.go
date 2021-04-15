package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
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
		p, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln("認証プロバイダーの取得に失敗しました：", provider, "-", err)
		}
		loginUrl, err := p.GetBeginAuthURL(nil, nil)
		if err != nil {
			log.Fatalln("GetBeginAuthURLの呼び出し中にエラーが発生しました：", provider, "-", err)
		}
		w.Header().Set("Location", loginUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)
	case "callback":
		p, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln("認証プロバイダーの取得に失敗しました：", provider, "-", err)
		}

		creds, err := p.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			log.Fatalln("認証を完了できませんでした：", provider, "-", err)
		}

		user, err := p.GetUser(creds)
		if err != nil {
			log.Fatalln("ユーザーの取得に失敗しました：", provider, "-", err)
		}

		authCookieValue := objx.New(map[string]interface{}{
			"name": user.Name(),
		}).MustBase64()
		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: authCookieValue,
			Path:  "/",
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "アクション%sには非対応です", action)
	}
}
