package main

import (
	"net/http"

	"github.com/intob/npchat/auth"
	"github.com/intob/npchat/cfg"
	"github.com/intob/npchat/kv"
	"github.com/intob/npchat/shareable"
	"github.com/intob/npchat/status"
)

func main() {
	cfg.InitViper()

	pool := kv.NewPool()

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		status.HandleGetHealth(w, r, pool)
	})

	http.HandleFunc("/register", auth.HandleRegister)
	http.HandleFunc("/login", auth.HandleLogin)
	http.HandleFunc("/shareable", shareable.HandleShareable)

	listenAndServe()
}
