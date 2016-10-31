package v1

import (
	"social_network/api/api"
)

func API() {
	api.Init([]string{"http://www.example.com", "http://www.mastergoco.com"})
	api.StartServer()
}
