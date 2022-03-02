package auth

import (
	"fmt"
	"net/http"
	"path"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/intob/npchat/kv"
	"github.com/intob/npchat/response"
)

const REGISTRATION = "registration/"
const REGISTRATION_TTL int64 = 120 // 2 minutes

// Tries to return a CredentialOptions object to the client
func HandleRegistrationStart(w http.ResponseWriter, r *http.Request, st *kv.Store, authn *webauthn.WebAuthn) {
	_, username := path.Split(r.URL.Path)
	if username == "" {
		http.Error(w, "username is missing: /register/{username}", http.StatusBadRequest)
		return
	}

	user, err := ensureUser(st, username)
	if err != nil {
		http.Error(w, "failed to lookup or create user", http.StatusInternalServerError)
		return
	}

	credentialOptions, sessionData, err := authn.BeginRegistration(user,
		webauthn.WithAuthenticatorSelection(
			protocol.AuthenticatorSelection{
				UserVerification: protocol.UserVerificationRequirement("required"),
			}),
	)
	if err != nil {
		http.Error(w, "failed to create credential options", http.StatusInternalServerError)
		return
	}

	sk, err := GenerateSessionKey()
	if err != nil {
		http.Error(w, "failed to generate session key", http.StatusInternalServerError)
	}

	// the client must set this as the authorization header
	w.Header().Add("session", sk)

	err = SetSessionData(st, REGISTRATION+sk, sessionData, REGISTRATION_TTL)
	if err != nil {
		http.Error(w, "failed to store session data", http.StatusInternalServerError)
		return
	}

	response.Json(w, credentialOptions, http.StatusOK)
}

func HandleRegistrationFinish(w http.ResponseWriter, r *http.Request, st *kv.Store, authn *webauthn.WebAuthn) {
	sessionKey := r.Header.Get("authorization")
	if sessionKey == "" {
		http.Error(w, "missing session key", http.StatusUnauthorized)
		return
	}

	sessionData, err := GetSessionData(st, REGISTRATION+sessionKey)
	if err != nil {
		http.Error(w, "no session found", http.StatusUnauthorized)
		return
	}

	parsedResponse, err := protocol.ParseCredentialCreationResponseBody(r.Body)
	if err != nil {
		http.Error(w, "failed to parse credential creation response", http.StatusBadRequest)
		return
	}

	user, err := GetUserByUsername(st, string(sessionData.UserID))
	if err != nil {
		http.Error(w, "failed to load user", http.StatusInternalServerError)
		return
	}

	cred, err := authn.CreateCredential(user, sessionData, parsedResponse)
	if err != nil {
		http.Error(w, "failed to create credential", http.StatusInternalServerError)
		return
	}

	err = StoreCredential(cred, user.Username, st)
	if err != nil {
		http.Error(w, "failed to store credential", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func ensureUser(st *kv.Store, username string) (*User, error) {
	// lookup user
	user, err := GetUserByUsername(st, username)
	if err != nil {
		// not found, create one

		fmt.Println("no user found for", username)

		user = &User{
			DisplayName: username, // default to username for now
			Username:    username,
		}
		err := StoreUser(st, user)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}
