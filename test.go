// This script displays a counting sequence from -9 to 99 on a 2-digit 7-segment display.
// Demonstrates basic number display functionality.
//
// This configures a 2-digit 7-segment Common Cathode display for arduino-nano

package main

import (
	"machine"

	"github.com/domi413/sevseg/sevseg"
)

func main() {
	displayConfig := sevseg.Config{
		Hardware: sevseg.CommonCathode,
		DigitPins: []machine.Pin{
			machine.GP0,
			machine.GP1,
			machine.GP2,
			machine.GP3,
		},
		SegmentPins: []machine.Pin{
			machine.GP4,  // A
			machine.GP5,  // B
			machine.GP6,  // C
			machine.GP7,  // D
			machine.GP8,  // E
			machine.GP9,  // F
			machine.GP10, // G
			machine.GP11, // DP
		},
		PWMType:         sevseg.SoftwarePWM,
		UseLeadingZeros: false,
	}

	dp, _ := sevseg.NewSevSeg(displayConfig)
	// dp.SetText("hall")
	dp.SetTemperatureWithUnit(23.234, 2, sevseg.TemperatureUnit.Celsius)

	// ctnr := uint16(0)
	// const threshold = uint16(200)

	for {
		// dp.DisplayTest(100)
		// ctnr++
		// if ctnr >= threshold {
		// 	// dp.ScrollTextLeft()
		// 	ctnr = 0
		// }

		// dp.Refresh()
		// time.Sleep(1 * time.Millisecond)
	}
}
