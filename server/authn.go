package server

import (
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/intob/npchat/cfg"
	"github.com/spf13/viper"
)

func getWebAuthn() (*webauthn.WebAuthn, error) {
	webAuthnConfig := &webauthn.Config{
		RPDisplayName: viper.GetString(cfg.APP_DISPLAY_NAME),
		RPID:          viper.GetString(cfg.APP_ID),
		RPOrigin:      viper.GetString(cfg.APP_ORIGIN),
		RPIcon:        viper.GetString(cfg.APP_ICON),
	}
	return webauthn.New(webAuthnConfig)
}
