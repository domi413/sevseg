package main

import (
	"machine"

	"sevseg/sevseg"
)

func main() {
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	displayConfig := sevseg.Config{
		Hardware: sevseg.CommonCathode,
		DigitPins: []machine.Pin{
			machine.D3,
			machine.D2,
		},
		SegmentPins: []machine.Pin{
			machine.D4,  // A
			machine.D5,  // B
			machine.D6,  // C
			machine.D7,  // D
			machine.D8,  // E
			machine.D9,  // F
			machine.D10, // G
			machine.D11, // DP
		},
		UseLeadingZeros: false,
	}

	dp, ok := sevseg.NewSevSeg(displayConfig)
	if !ok {
		led.High()
	}

	dp.SetNumber(3)
	// dp.SetNumberWithDecimal(12, 1)
	// dp.SetNumberWithMultipleDecimals(42, []uint8{0})
	// dp.SetSegment([]uint8{0b00111001, 0b10001111})
	// dp.SetText("Er")

	for {
		dp.Refresh()
	}
	// for count := int8(0); ; count = (count + 1) % 100 {
	// 	if err := display.SetNumber(count); !err {
	// 		panic(err)
	// 	}

	// 	for range 300 {
	// 		display.Refresh()
	// 		time.Sleep(time.Millisecond)
	// 	}
	// }
}
