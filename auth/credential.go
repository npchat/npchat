package auth

import (
	"bytes"
	"encoding/gob"

	"github.com/duo-labs/webauthn/webauthn"
	"github.com/intob/npchat/kv"
)

const CREDENTIAL = "credential/"

// TODO: support multiple credentials per user
func StoreCredential(credential *webauthn.Credential, username string, st *kv.Store) error {

	var buf bytes.Buffer

	err := gob.NewEncoder(&buf).Encode(credential)
	if err != nil {
		return err
	}

	return st.Set(CREDENTIAL+username, buf.Bytes(), 0)
}
