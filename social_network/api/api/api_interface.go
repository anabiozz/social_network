package api

import (
	"net/http"
)

func APIInterface(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "src/social_network/api/api/interface.html")
}
