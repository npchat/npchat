package cfg

import (
	"fmt"
	"runtime"

	"github.com/spf13/viper"
)

const TLS_CERTFILE = "TLSCertFile"
const TLS_KEYFILE = "TLSKeyFile"
const PORT = "Port"

const ROCKET_NET = "RocketNetwork"
const ROCKET_ADDRESS = "RocketAddress"
const ROCKET_TLS_CERTFILE = "RocketTLSCertFile"
const ROCKET_TLS_KEYFILE = "RocketTLSKeyFile"
const ROCKET_AUTH = "RocketAuthSecret"
const ROCKET_WORKERS_MIN = "RocketWorkersMin"

func InitViper() {
	initViperDefaults()

	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/npchat/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}

func initViperDefaults() {
	viper.SetDefault(PORT, 8000)
	viper.SetDefault(ROCKET_NET, "tcp")
	viper.SetDefault(ROCKET_ADDRESS, ":8100")
	viper.SetDefault(ROCKET_WORKERS_MIN, uint32(runtime.NumCPU()))
}
