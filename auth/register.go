package auth

import (
	"net/http"

	"github.com/duo-labs/webauthn/webauthn"
	"github.com/intob/npchat/kv"
)

func HandleRegistrationStart(w http.ResponseWriter, r *http.Request, st *kv.Store, authn *webauthn.WebAuthn) {
	w.Write([]byte("register"))

}
