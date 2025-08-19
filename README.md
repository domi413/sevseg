# TinyGo 7-Segment Display Library

A comprehensive TinyGo library for controlling multiplexed 7-segment displays
with support for both common-anode and common-cathode configurations.

This library is inspired by the popular
[SevSeg](https://github.com/DeanIsMe/SevSeg/) library for Arduino adapted and
improved for TinyGo.

## Installation

Include `import "github.com/domi413/sevseg"` in your TinyGo project and run the following commands to set up the module:

```bash
go mod init {your-project} # May already be done
go get github.com/domi413/sevseg@v0.1.1
```

## Hardware Setup

### 7-Segment Display pin-out

```
 AAA
F   B
F   B
 GGG
E   C
E   C
 DDD  DP
```

### Pin Mapping

The library expects segment pins in this order:

1. Segment A
2. Segment B
3. Segment C
4. Segment D
5. Segment E
6. Segment F
7. Segment G
8. Decimal Point (DP)

## Simple example for a 2-digit common-cathode display on an arduino nano

```go
package main

import (
  "machine"

  "sevseg"
)

func main() {
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

	display, ok := sevseg.NewSevSeg(displayConfig)
	if !ok {
		panic("Failed to create display")
	}

	if !display.SetNumber(69) {
		panic("Failed to set number")
	}

	for {
		// Keep refreshing the display
		display.Refresh()
	}
}
```

## API Reference

### Configuration

```go
type Config struct {
    Hardware        displayType     // CommonAnode or CommonCathode
    DigitPins       []machine.Pin   // Pins for multiplexing the segments
    SegmentPins     []machine.Pin   // Pins controlling segments (A-G + DP)
    UseLeadingZeros bool            // Whether to display leading zeros for numbers
}
```

### Methods

#### `NewSevSeg(config Config) (*SevSeg, bool)`

Creates a new SevSeg instance with the provided configuration. Returns the
display instance and a boolean indicating success (true) or failure (false).

#### `Clear()`

Clears the display, turning off all segments by setting them to blank space.

#### `Off()`

Turns off the display completely by clearing all digit and segment pins. This
provides immediate shutdown without requiring a `Refresh()` call.

#### `SetBrightness(brightness uint8) error`

Sets the brightness of the display. Currently not implemented.

_Note: This method is planned for future implementation where brightness would
be controlled through PWM or timing adjustments._

#### `SetNumber(number int8) bool`

Sets the number to be displayed. The number can be positive or negative. Returns
`true` on success, `false` if the number exceeds the display's capacity.

#### `SetNumberWithDecimal(number int8, decimalPoint uint8) bool`

Sets the number to be displayed with a decimal point at the specified position.
The `decimalPoint` parameter is zero-indexed from the right.

E.g., `SetNumberWithDecimal(1234, 1)` will display `123.4`.

Returns `true` on success, `false` if the number exceeds capacity, decimal point
position is invalid, or display doesn't have all 8 segment pins defined.

#### `SetNumberWithMultipleDecimals(number int8, decimalPoints []uint8) bool`

Sets the number to be displayed with multiple decimal points at specified
positions. The `decimalPoints` slice contains zero-indexed positions from the
right.

E.g., `SetNumberWithMultipleDecimals(1234, []uint8{1, 2})` will display `12.3.4`.

Returns `true` on success, `false` if the number exceeds capacity, any decimal
point position is invalid, or display doesn't have 8 segment pins.

#### `SetSegment(pattern []uint8) bool`

Sets a custom segment pattern to be displayed. The `pattern` slice should
contain bit values representing each segment (A-G, DP) in the order defined in
the pin mapping.

E.g. to display the following pattern:

```
‾│ │‾  │‾  ‾│
_|.│_  │_. _│
```

This would equal the bit pattern:

`[]uint8{0b10001111, 0b00111001, 0b10111001, 0b00001111}`

The segments are displayed from right to left (index 0 = rightmost digit). If
fewer segments are provided than available digits, remaining digits on the left
will be cleared.

Returns `true` on success, `false` if the pattern length exceeds the number of digits.

#### `SetText(text string) bool`

Displays text on the 7-segment display.

Text is displayed from left to right. If the text is shorter than the number of
digits, the remaining segments on the right will be cleared.

Returns `true` on success, `false` if the text is longer than the number of
digits or contains unsupported characters.

Supported characters:

- Digits: `0-9`
- Letters: `A-Z` (case insensitive)
- Special: ` ` (space), `-` (minus), `.` (period), `°` (degree symbol), `_` (underscore)

#### `Refresh() bool`

Refreshes the display by cycling through each digit. Must be called frequently
(recommended >100Hz) to maintain a stable, flicker-free display.

Returns `true` on success, `false` if the display is not properly initialized.

## Troubleshooting

### Display is dim or flickering

- Increase the refresh rate by calling `Refresh()` more frequently (aim for >100Hz)
- Check your resistor values (too high resistance can cause dimming)
- Verify power supply can handle the current draw

### Numbers appear backwards

- Check if you've connected digit pins in the correct order
- Verify segment pin connections match the expected order (A-G, DP)

### Display shows wrong characters

- Verify you're using the correct display type (CommonAnode vs CommonCathode)
- Check segment pin wiring order

### Some segments don't light up

- Check for loose connections
- Verify resistor values and connections
- Test individual segments with a multimeter

### Method returns false

- Check that your number fits within the available digits
- For decimal operations, ensure you have 8 segment pins (including DP)
- For text operations, verify all characters are supported
- For segment patterns, ensure pattern length doesn't exceed digit count

## Error Handling

This library uses boolean return values instead of Go's standard error interface
for memory efficiency in embedded environments:

- `true` indicates successful operation
- `false` indicates an error occurred

It's recommended to check return values for methods that can fail:

```go
if !display.SetNumber(123) {
    // Handle error - number too large for display
}

if !display.SetText("HELLO") {
    // Handle error - text too long or unsupported characters
}
```
