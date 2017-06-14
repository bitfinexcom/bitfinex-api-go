package tests

import (
	"log"
	"os"
)

var (
	auth   = false
	key    = os.Getenv("BFX_API_KEY")
	secret = os.Getenv("BFX_API_SECRET")
)

func init() {

	if key != "" && secret != "" {
		auth = true
	} else {
		log.Println("No authentication credentials provided so running only public tests.")
	}
}
