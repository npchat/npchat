package shareable

import "net/http"

func HandleShareable(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("shareable"))
}
