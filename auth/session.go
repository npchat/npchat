package auth

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"time"

	"github.com/duo-labs/webauthn/webauthn"
	"github.com/intob/npchat/kv"
)

const SESSION = "session/"
const SESSION_KEY_LEN = 32

var SESSION_TTL = int64(time.Hour.Seconds())

func GenerateSessionKey() (string, error) {
	buf := make([]byte, SESSION_KEY_LEN)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(buf), nil
}

func SetSessionData(st *kv.Store, sessionKey string, sessionData *webauthn.SessionData) error {
	var buf bytes.Buffer
	gob.NewEncoder(&buf).Encode(sessionData)
	return st.Set(SESSION+sessionKey, buf.Bytes(), SESSION_TTL)
}

func GetSessionData(st *kv.Store, sessionKey string) (webauthn.SessionData, error) {
	var buf bytes.Buffer
	resp, err := st.Get(SESSION + sessionKey)
	if err != nil {
		return webauthn.SessionData{}, err
	}
	buf.Write(resp.Value)
	sessionData := webauthn.SessionData{}
	err = gob.NewDecoder(&buf).Decode(&sessionData)
	return sessionData, err
}
