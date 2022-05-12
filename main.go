package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"

	// "github.com/rs/cors"
	"github.com/zitadel/oidc/pkg/client/rp"
	httphelper "github.com/zitadel/oidc/pkg/http"
	"github.com/zitadel/oidc/pkg/oidc"
)

func main() {
	router := chi.NewRouter()

	clientID := "1ald2ons4fa6remlmcbvescb5b"
	clientSecret := ""
	redirectUrl := "http://localhost:18888/auth/callback"
	issuer := "https://cognito-idp.ap-northeast-1.amazonaws.com/ap-northeast-1_FZjE5Kfvl"
	scopes := []string{"openid", "email", "profile"}

	hashKey := "test1234test1234"
	encryptKey := "test1234test1234"
	cookieHandler := httphelper.NewCookieHandler([]byte(hashKey), []byte(encryptKey), httphelper.WithUnsecure())
	options := []rp.Option{
		rp.WithCookieHandler(cookieHandler),
	}
	relyingParty, err := rp.NewRelyingPartyOIDC(issuer, clientID, clientSecret, redirectUrl, scopes, options...)
	if err != nil {
		panic(err)
	}

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<a href=\"/auth/login\">login</a>\n\n"))
	})

	router.Route("/auth", func(r chi.Router) {
		r.Handle("/login", rp.AuthURLHandler(uuid.NewString, relyingParty))
		r.Handle("/callback", rp.CodeExchangeHandler(func(w http.ResponseWriter, r *http.Request, tokens *oidc.Tokens, state string, rp rp.RelyingParty) {
			// 必要な処理
			http.Redirect(w, r, "/home", http.StatusFound)
		}, relyingParty))
	})

	router.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("logined!"))
	})

	http.ListenAndServe("localhost:18888", router)
}
