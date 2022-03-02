package auth

import (
	"github.com/intob/npchat/kv"
	"github.com/intob/rocketkv/protocol"
)

func VerifyAuthSessionKey(sessionKey string, st *kv.Store) bool {
	resp, err := st.Get(AUTHED + sessionKey)
	if err != nil {
		return false
	}
	if resp.Status != protocol.StatusOk {
		return false
	}
	return true
}
