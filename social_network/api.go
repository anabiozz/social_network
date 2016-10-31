package main

import (
	"social_network/api/v1"
	"social_network/api/v2"
)

var API struct{}

func main() {

	v := 1

	if v == 1 {
		v1.API()
		// do stuff with API v1
	} else {
		v2.API()
		// do stuff with API v2
	}

}
