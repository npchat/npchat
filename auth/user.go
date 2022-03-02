package auth

import (
	"bytes"
	"encoding/gob"
	"errors"

	"github.com/duo-labs/webauthn/webauthn"
	"github.com/intob/npchat/kv"
	"github.com/intob/rocketkv/protocol"
)

const USER = "user/"

// User data, implements webauthn.User
type User struct {
	DisplayName string
	Username    string
	Avatar      string // URL
	Credentials []webauthn.Credential
}

// User ID according to the Relying Party,
// using username for now
func (u *User) WebAuthnID() []byte {
	return []byte(u.Username)
}

// User Name according to the Relying Party
func (u *User) WebAuthnName() string {
	return u.Username
}

// Display Name of the user
func (u *User) WebAuthnDisplayName() string {
	return u.DisplayName
}

// User's icon url
func (u *User) WebAuthnIcon() string {
	return u.Avatar
}

// Credentials owned by the user
func (u *User) WebAuthnCredentials() []webauthn.Credential {
	return u.Credentials
}

func GetUserByUsername(st *kv.Store, username string) (*User, error) {
	resp, err := st.Get(USER + username)
	if err != nil {
		return nil, err
	}
	if resp.Status != protocol.StatusOk {
		return nil, errors.New("user not found")
	}
	var buf bytes.Buffer
	buf.Write(resp.Value)
	user := &User{}
	err = gob.NewDecoder(&buf).Decode(user)
	return user, err
}

func StoreUser(st *kv.Store, user *User) error {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(user)
	if err != nil {
		return err
	}
	return st.Set(USER+user.Username, buf.Bytes(), 0)
}
