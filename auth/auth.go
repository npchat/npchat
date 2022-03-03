package auth

import (
	"bytes"
	"net/http"

	"github.com/intob/npchat/kv"
	"github.com/intob/rocketkv/protocol"
)

// Verifies that the session key is valid
// and that the associated IP addr matches the request
func VerifyAuthSessionKey(r *http.Request, sessionKey string, st *kv.Store) bool {
	resp, err := st.Get(AUTHED + sessionKey)
	if err != nil {
		return false
	}
	if resp.Status != protocol.StatusOk {
		return false
	}
	// check that stored IP matches that of request
	if !bytes.Equal([]byte(r.RemoteAddr), resp.Value) {
		return false
	}
	return true
}
