package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/duo-labs/webauthn/webauthn"
	"github.com/intob/npchat/auth"
	"github.com/intob/npchat/cfg"
	"github.com/intob/npchat/kv"
	"github.com/intob/npchat/shareable"
	"github.com/intob/npchat/status"
	"github.com/spf13/viper"
)

const timeout = 5 * time.Second

type WebServer struct {
	server   *http.Server
	store    kv.Store
	webAuthn *webauthn.WebAuthn
}

func NewWebServer() *WebServer {
	authn, err := getWebAuthn()
	if err != nil {
		panic(err)
	}

	port := viper.GetInt(cfg.PORT)
	addr := fmt.Sprintf(":%v", port)
	server := &http.Server{
		Addr:         addr,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	}

	ws := &WebServer{
		server:   server,
		store:    *kv.NewStore(),
		webAuthn: authn,
	}

	ws.registerRoutes()
	go ws.start()

	fmt.Printf("listening on %v\r\n", addr)

	return ws
}

func (ws *WebServer) registerRoutes() {
	http.HandleFunc("/health/", func(w http.ResponseWriter, r *http.Request) {
		status.HandleGetHealth(w, r, &ws.store)
	})

	http.HandleFunc("/register/", func(w http.ResponseWriter, r *http.Request) {
		auth.HandleRegistrationStart(w, r, &ws.store, ws.webAuthn)
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		auth.HandleLoginStart(w, r, &ws.store, ws.webAuthn)
	})

	http.HandleFunc("/shareable", shareable.HandleShareable)
}

func (ws *WebServer) start() {
	certFile := viper.GetString(cfg.TLS_CERTFILE)
	keyFile := viper.GetString(cfg.TLS_KEYFILE)

	var err error
	if certFile != "" && keyFile != "" {
		fmt.Println("expecting HTTPS connections")
		err = ws.server.ListenAndServeTLS(certFile, keyFile)
	} else {
		err = ws.server.ListenAndServe()
	}
	if err != nil {
		fmt.Println(err)
	}
}

func (ws *WebServer) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return ws.server.Shutdown(ctx)
}
