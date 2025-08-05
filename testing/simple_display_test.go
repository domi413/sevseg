package main

// This script iterates over all fields of every 7-segment display. This can be
// used to verify that the display is working and connected correctly.
//
// This configures a 2-digit 7-segment Common Cathode display for a arduino-nano

import (
	"machine"
	"time"
)

func main() {
	digitPins := []machine.Pin{
		machine.D3,
		machine.D2,
	}
	segmentPins := []machine.Pin{
		machine.D4,  // A
		machine.D5,  // B
		machine.D6,  // C
		machine.D7,  // D
		machine.D8,  // E
		machine.D9,  // F
		machine.D10, // G
		machine.D11, // DP
	}

	for _, pin := range digitPins {
		pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
		pin.Low()
	}
	for _, pin := range segmentPins {
		pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
		pin.High()
	}

	for {
		for _, digitPin := range digitPins {
			digitPin.High()
			for _, segmentPin := range segmentPins {
				segmentPin.Low()
				time.Sleep(500 * time.Millisecond)
				segmentPin.High()
			}
			digitPin.Low()
		}
	}
}
