package main

import (
	"fmt"
	"os"

	"github.com/achetronic/wizgo/pkg/wizgo"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: jojo <mode>")
		fmt.Println("Modes: on, bedtime, party, off")
		return
	}

	var scene int

	switch os.Args[1] {
	case "on":
		scene = 12
	case "bedtime":
		scene = 6
	case "party":
		scene = 4
	case "off":
		scene = 0
	default:
		fmt.Println("Usage: jojo <mode>")
		fmt.Println("Modes: on, bedtime, party, off")
		return
	}

	wizClient, err := wizgo.CreateWizClient("192.168.4.145", 38899)
	if err != nil {
		panic(err)
	}

	// Off = 0
	if scene == 0 {
		wizClient.TurnOff()
		return
	}

	// Party = 4
	// Cozy = 6
	// Daylight = 12
	wizClient.SetScene(scene)

	if scene == 4 {
		wizClient.SetSpeed(200)
	}

	if scene == 6 {
		wizClient.SetBrightness(10)
	} else {
		wizClient.SetBrightness(100)
	}
}
