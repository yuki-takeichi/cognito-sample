package main

import (
	"net/http"

	"github.com/go-chi/chi"
	// "github.com/rs/cors"
	"github.com/zitadel/oidc/pkg/client/rp"
	"github.com/zitadel/oidc/pkg/oidc"
)

func main() {
	router := chi.NewRouter()

	clientID := "1ald2ons4fa6remlmcbvescb5b"
	clientSecret := ""
	redirectUrl := "http://localhost:18888/auth/callback"
	issuer := "https://cognito-idp.ap-northeast-1.amazonaws.com/ap-northeast-1_FZjE5Kfvl"
	scopes := []string{"openid", "email", "profile"}

	relyingParty, err := rp.NewRelyingPartyOIDC(issuer, clientID, clientSecret, redirectUrl, scopes)
	if err != nil {
		panic(err)
	}

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<a href=\"/auth/login\">login</a>\n\n"))
	})

	router.Route("/auth", func(r chi.Router) {
		r.Handle("/login", rp.AuthURLHandler(func() string { return "hoge" }, relyingParty))
		r.Handle("/callback", rp.CodeExchangeHandler(func(w http.ResponseWriter, r *http.Request, tokens *oidc.Tokens, state string, rp rp.RelyingParty) {
			// w.Write([]byte("done"))
			http.Redirect(w, r, "/home", http.StatusFound)
		}, relyingParty))
	})

	router.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("logined!"))
	})

	http.ListenAndServe("localhost:18888", router)
}
