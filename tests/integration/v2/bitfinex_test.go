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
	log.Println("Authenticated tests disabled.")
	auth = false
}
