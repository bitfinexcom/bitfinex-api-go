package main

import (
	"log"
	"os"
	"time"

	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/common"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/pulse"
	"github.com/bitfinexcom/bitfinex-api-go/v2/rest"
	"github.com/davecgh/go-spew/spew"
)

// Set BFX_API_KEY and BFX_API_SECRET:
//
// export BFX_API_KEY=<your-api-key>
// export BFX_API_SECRET=<your-api-secret>
//
// you can obtain it from https://www.bitfinex.com/api

func main() {
	key := os.Getenv("BFX_API_KEY")
	secret := os.Getenv("BFX_API_SECRET")

	c := rest.
		NewClient().
		Credentials(key, secret)

	pulseProfile(c)
	publicPulseHistory(c)
	addPulse(c)
	addComment(c)
	pulseHistory(c)
	deletePulse(c)
}

func pulseProfile(c *rest.Client) {
	nn := rest.Nickname("Bitfinex")
	resp, err := c.Pulse.PublicPulseProfile(nn)
	if err != nil {
		log.Fatalf("PublicPulseProfile: %s", err)
	}

	spew.Dump(resp)
}

func publicPulseHistory(c *rest.Client) {
	now := time.Now()
	millis := now.UnixNano() / 1000000
	from := common.Mts(millis)

	pulseHist, err := c.Pulse.PublicPulseHistory(2, from)
	if err != nil {
		log.Fatalf("PublicPulseHistory: %s", err)
	}

	spew.Dump(pulseHist)
}

func addPulse(c *rest.Client) {
	payload := &pulse.Pulse{
		Title:   "GO GO GO GO GO GO TITLE",
		Content: "GO GO GO GO GO GO Content",
	}

	resp, err := c.Pulse.AddPulse(payload)
	if err != nil {
		log.Fatalf("AddPulse: %s", err)
	}

	spew.Dump(resp)
}

func addComment(c *rest.Client) {
	pulse1 := &pulse.Pulse{
		Title:   "TITLE TO BE COMMENTED ON",
		Content: "CONTENT TO BE COMMENTED ON",
	}

	p1, err := c.Pulse.AddPulse(pulse1)
	if err != nil {
		log.Fatalf("addComment:AddPulse: %s", err)
	}

	spew.Dump(p1)

	pulse2 := &pulse.Pulse{
		Title:   "TITLE THAT HOLDS COMMENT",
		Content: "CONTENT THAT HOLDS COMMENT",
		Parent:  p1.ID,
	}

	p2, err := c.Pulse.AddComment(pulse2)
	if err != nil {
		log.Fatalf("addComment:AddComment: %s", err)
	}

	spew.Dump(p2)
}

func pulseHistory(c *rest.Client) {
	pulseHist, err := c.Pulse.PulseHistory()
	if err != nil {
		log.Fatalf("PulseHistory: %s", err)
	}

	spew.Dump(pulseHist)
}

func deletePulse(c *rest.Client) {
	payload := &pulse.Pulse{
		Title:   "TITLE TO BE DELETED",
		Content: "CONTENT TO BE DELETED",
	}

	p, err := c.Pulse.AddPulse(payload)
	if err != nil {
		log.Fatalf("deletePulse:AddPulse: %s", err)
	}

	spew.Dump(p)

	d, err := c.Pulse.DeletePulse(p.ID)
	if err != nil {
		log.Fatalf("deletePulse:DeletePulse: %s", err)
	}

	spew.Dump(d)
}
