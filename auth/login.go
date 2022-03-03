package auth

import (
	"net/http"
	"path"
	"time"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/intob/npchat/kv"
	"github.com/intob/npchat/response"
)

const LOGIN = "login/"
const LOGIN_TTL = time.Minute * 2

const AUTHED = "authed/"
const AUTHED_TTL = time.Hour

func HandleLoginStart(w http.ResponseWriter, r *http.Request, st *kv.Store, authn *webauthn.WebAuthn) {
	_, username := path.Split(r.URL.Path)
	if username == "" {
		http.Error(w, "username is missing: /login/{username}", http.StatusBadRequest)
		return
	}

	user, err := GetUserByUsername(st, username)
	if err != nil {
		http.Error(w, "no user found", http.StatusBadRequest)
		return
	}

	options, sessionData, err := authn.BeginLogin(user)
	if err != nil {
		http.Error(w, "failed to begin login process", http.StatusInternalServerError)
		return
	}

	sk, err := GenerateSessionKey()
	if err != nil {
		http.Error(w, "failed to generate session key", http.StatusInternalServerError)
		return
	}
	w.Header().Add("session", sk)

	err = SetSessionData(st, LOGIN+sk, sessionData, int64(LOGIN_TTL.Seconds()))
	if err != nil {
		http.Error(w, "failed to store login session", http.StatusInternalServerError)
		return
	}

	response.Json(w, options, http.StatusOK)
}

func HandleLoginFinish(w http.ResponseWriter, r *http.Request, st *kv.Store, authn *webauthn.WebAuthn) {
	sessionKey := r.Header.Get("authorization")
	if sessionKey == "" {
		http.Error(w, "missing session key", http.StatusUnauthorized)
		return
	}

	parsedResponse, err := protocol.ParseCredentialRequestResponseBody(r.Body)
	if err != nil {
		http.Error(w, "failed to parse credential creation response", http.StatusBadRequest)
		return
	}

	sessionData, err := GetSessionData(st, LOGIN+sessionKey)
	if err != nil {
		http.Error(w, "no session found", http.StatusUnauthorized)
		return
	}

	user, err := GetUserByUsername(st, string(sessionData.UserID))
	if err != nil {
		http.Error(w, "failed to load user", http.StatusInternalServerError)
		return
	}

	_, err = authn.ValidateLogin(user, sessionData, parsedResponse)
	if err != nil {
		http.Error(w, "failed to validate login", http.StatusUnauthorized)
		return
	}

	sk, err := GenerateSessionKey()
	if err != nil {
		http.Error(w, "failed to generate session key", http.StatusInternalServerError)
		return
	}
	w.Header().Add("session", sk)

	st.Set(AUTHED+sk, []byte(r.RemoteAddr), int64(AUTHED_TTL.Seconds()))

	w.WriteHeader(http.StatusOK)
}
