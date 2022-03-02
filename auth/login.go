package auth

import (
	"net/http"

	"github.com/duo-labs/webauthn/webauthn"
	"github.com/intob/npchat/kv"
)

func HandleLoginStart(w http.ResponseWriter, r *http.Request, store *kv.Store, authn *webauthn.WebAuthn) {
	w.Write([]byte("login"))
}
