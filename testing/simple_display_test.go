// This script iterates over all fields of every 7-segment display. This can be
// used to verify that the display is working and connected correctly.
//
// This configures a 2-digit 7-segment Common Cathode display for a arduino-nano

package main

import (
	"machine"
	"time"
)

type displayType int

const (
	CommonAnode displayType = iota
	CommonCathode
)

func main() {
	const displayConfig = CommonCathode

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
	}
	for _, pin := range segmentPins {
		pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	}

	clearPins(digitPins, segmentPins, displayConfig)

	for {
		for _, digitPin := range digitPins {
			setDigitPin(digitPin, true, displayConfig)

			for _, segmentPin := range segmentPins {
				setSegmentPin(segmentPin, true, displayConfig)
				time.Sleep(500 * time.Millisecond)
				setSegmentPin(segmentPin, false, displayConfig)
			}
			setDigitPin(digitPin, false, displayConfig)
		}
	}
}

func clearPins(digitPins, segmentPins []machine.Pin, config displayType) {
	for _, pin := range digitPins {
		setDigitPin(pin, false, config)
	}
	for _, pin := range segmentPins {
		setSegmentPin(pin, false, config)
	}
}

func setDigitPin(pin machine.Pin, active bool, config displayType) {
	if config == CommonCathode {
		if active {
			pin.Low()
		} else {
			pin.High()
		}
	} else { // CommonAnode
		if active {
			pin.High()
		} else {
			pin.Low()
		}
	}
}

func setSegmentPin(pin machine.Pin, active bool, config displayType) {
	if config == CommonCathode {
		if active {
			pin.High()
		} else {
			pin.Low()
		}
	} else { // CommonAnode
		if active {
			pin.Low()
		} else {
			pin.High()
		}
	}
}
