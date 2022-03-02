package auth

import "net/http"

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("register"))
}
