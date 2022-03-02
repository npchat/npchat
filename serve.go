package main

import (
	"fmt"
	"net/http"

	"github.com/intob/npchat/cfg"
	"github.com/spf13/viper"
)

func listenAndServe() {
	port := viper.GetInt(cfg.PORT)
	addr := fmt.Sprintf(":%v", port)

	certFile := viper.GetString(cfg.TLS_CERTFILE)
	keyFile := viper.GetString(cfg.TLS_KEYFILE)

	var err error
	if certFile != "" && keyFile != "" {
		fmt.Println("expecting HTTPS connections")
		err = http.ListenAndServeTLS(addr, certFile, keyFile, nil)
	} else {
		err = http.ListenAndServe(addr, nil)
	}
	if err != nil {
		panic(fmt.Errorf("failed to start server: %w", err))
	}
	fmt.Printf("listening on %v\r\n", addr)
}
