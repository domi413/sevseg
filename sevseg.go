//go:build tinygo

// Package sevseg is a library for controlling 7-segment displays.
package sevseg

import (
	"machine"
)

// digitCodeMap indicate which segments must be illuminated to display a
// specific digit or character.
var digitCodeMap = []uint8{
	//.GFEDCBA    Idx   ASCII   Symbol   7-segment map:
	0b00111111, // 0      0      '0'          AAA
	0b0000110,  // 1      1      '1'         F   B
	0b01011011, // 2      2      '2'         F   B
	0b01001111, // 3      3      '3'          GGG
	0b01100110, // 4      4      '4'         E   C
	0b01101101, // 5      5      '5'         E   C
	0b01111101, // 6      6      '6'          DDD
	0b00000111, // 7      7      '7'
	0b01111111, // 8      8      '8'
	0b01101111, // 9      9      '9'

	0b01110111, // 10    65      'A'
	0b01111100, // 11    66      'b'
	0b00111001, // 12    67      'C'
	0b01011110, // 13    68      'd'
	0b01111001, // 14    69      'E'
	0b01110001, // 15    70      'F'
	0b00111101, // 16    71      'G'
	0b01110110, // 17    72      'H'
	0b00110000, // 18    73      'I'
	0b00001110, // 19    74      'J'
	0b01110110, // 20    75      'K'  Same as 'H'
	0b00111000, // 21    76      'L'
	0b00000000, // 22    77      'M'  NO DISPLAY
	0b01010100, // 23    78      'n'
	0b00111111, // 24    79      'O'
	0b01110011, // 25    80      'P'
	0b01100111, // 26    81      'q'
	0b01010000, // 27    82      'r'
	0b01101101, // 28    83      'S'
	0b01111000, // 29    84      't'
	0b00111110, // 30    85      'U'
	0b00111110, // 31    86      'V'  Same as 'U'
	0b00000000, // 32    87      'W'  NO DISPLAY
	0b01110110, // 33    88      'X'  Same as 'H'
	0b01101110, // 34    89      'y'
	0b01011011, // 35    90      'Z'  Same as '2'

	0b00000000, // 36    32      ' '  BLANK
	0b01000000, // 37    45      '-'  DASH / MINUS
	0b10000000, // 38    46      '.'  PERIOD / DECIMAL POINT
	0b01100011, // 39    42      '°'  DEGREE
	0b00001000, // 40    95      '_'  UNDERSCORE
}

type displayType int

// CommonAnode and CommonCathode are used to define the type of 7-segment display.
const (
	CommonAnode displayType = iota
	CommonCathode
)

// Config holds the configuration for a 7-segment display.
type Config struct {
	// Hardware defines the type of 7-segment display.
	// It can be either CommonAnode or CommonCathode.
	Hardware displayType

	// DigitPins defines the pins used control/multiplex the digits.
	DigitPins []machine.Pin

	// SegmentPins defines the pins used to control the segments of the display.
	// Normally, these are 7 or 8 pins, depending on whether a decimal point is used.
	SegmentPins []machine.Pin

	// UseLeadingZeros defines whether leading zeros should be displayed.
	UseLeadingZeros bool
}

// SevSeg represents a 7-segment display.
type SevSeg struct {
	config                displayType
	digitPins             []machine.Pin
	segmentPins           []machine.Pin
	useLeadingZeros       bool
	updatedDisplay        []uint8
	currentDigitToRefresh uint8
}

// NewSevSeg creates a new instance of sevSeg with the provided configuration.
func NewSevSeg(cfg Config) (*SevSeg, bool) {
	if len(cfg.DigitPins) == 0 {
		return nil, false
	}

	if len(cfg.SegmentPins) < 7 || len(cfg.SegmentPins) > 8 {
		return nil, false
	}

	for _, pin := range cfg.DigitPins {
		pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	}
	for _, pin := range cfg.SegmentPins {
		pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	}

	s := &SevSeg{
		config:          cfg.Hardware,
		digitPins:       cfg.DigitPins,
		segmentPins:     cfg.SegmentPins,
		useLeadingZeros: cfg.UseLeadingZeros,
		updatedDisplay:  make([]uint8, len(cfg.DigitPins)),
	}

	s.Off()

	return s, true
}

// Clear clears the display by setting all segments to blank space.
func (s *SevSeg) Clear() {
	for i := range s.updatedDisplay {
		s.updatedDisplay[i] = digitCodeMap[36] // Blank space
	}
}

// Off turns off the display by setting all digit and segment pins to their
// respective off state, depending on the display type.
//
// This turns off the display immediately without calling Refresh.
func (s *SevSeg) Off() {
	s.clearDigitPins()
	s.clearSegmentPins()
}

// SetBrightness sets the brightness of the display.
//
// Takes the brightness level in percentage (0-100) as an argument.
func (s *SevSeg) SetBrightness(brightness uint8) error {
	// TODO:
	// An idea is that we add a new variable which holds the brightness level.
	// The Refresh method can then use this variable to update the display with
	// the according brightness level.
	return nil
}

// SetNumber sets the number to be displayed on the 7-segment display.
func (s *SevSeg) SetNumber(number int8) bool {
	if !s.checkAvailableDigits(number) {
		return false
	}

	initPattern := digitCodeMap[36] // Blank space
	if s.useLeadingZeros {
		initPattern = digitCodeMap[0] // Zero
	}
	for i := range s.updatedDisplay {
		s.updatedDisplay[i] = initPattern
	}

	isNegative := number < 0
	if isNegative {
		number = -number
	}

	position := 0
	if number == 0 {
		s.updatedDisplay[position] = digitCodeMap[0]
	} else {
		for number > 0 && position < len(s.digitPins) {
			digit := number % 10
			s.updatedDisplay[position] = digitCodeMap[digit]
			number /= 10
			position++
		}
	}

	if isNegative && position >= 0 {
		s.updatedDisplay[position] = digitCodeMap[37] // Minus
	}

	return true
}

// SetNumberWithDecimal sets the number to be displayed on the 7-segment display,
// including a decimal point at a specified position.
//
// decimalPointPlace specifies the position of the decimal point from right to left.
//
// E.g. for a 4-digit display, decimalPointPlace = 1 would look like this: 000.0
func (s *SevSeg) SetNumberWithDecimal(number int8, decimalPoint uint8) bool {
	return s.SetNumberWithMultipleDecimals(number, []uint8{decimalPoint})
}

// SetNumberWithMultipleDecimals sets the number to be displayed on the 7-segment display,
// including multiple decimal points at specified positions.
//
// decimalPointPlaces is a slice of positions for the decimal points from right to left.
//
// E.g. for a 4-digit display, decimalPointPlaces = []uint{1, 2} would look like this: 00.0.0
func (s *SevSeg) SetNumberWithMultipleDecimals(number int8, decimalPoints []uint8) bool {
	if len(decimalPoints) == 0 {
		return false
	}

	if len(s.segmentPins) < 8 {
		return false
	}

	if !s.SetNumber(number) {
		return false
	}

	for _, decimalPos := range decimalPoints {
		if decimalPos > uint8(len(s.digitPins)) {
			return false
		}
		s.updatedDisplay[decimalPos] |= digitCodeMap[38] // Decimal point
	}

	return true
}

// SetSegment can be used to display any arbitrary segment pattern.
// E.g. to display the following pattern:
//
//		‾│ │‾  │‾  ‾│
//		_|.│_  │_. _│
//	    4   3   2  1
//
// This would equal the bit pattern:
//
// {0b00001111, 0b10111001, 0b00111001, 0b10001111}
//
//	1           2           3           4
//
// Note that the left most segment is the MSB and is the last pattern in the slice.
//
// The segments will be displayed from right to left. This means that if fewer
// segments are defined than digits, the remaining segments (on the left) will be cleared.
func (s *SevSeg) SetSegment(pattern []uint8) bool {
	if len(pattern) > len(s.digitPins) {
		return false
	}

	for i, segment := range pattern {
		s.updatedDisplay[i] = segment
	}

	return true
}

// SetText displays a text.
//
// If the text is longer than the number of digits, an error is returned.
//
// The text is written from left to right, meaning that if the text is shorter
// than the number of digits, the remaining segments (on the right) will be
// cleared.
func (s *SevSeg) SetText(text string) bool {
	if len(text) > len(s.digitPins) {
		return false
	}

	s.Clear()

	for i, char := range []byte(text) {
		segment, ok := s.charToSegmentPattern(char)
		if !ok {
			return false
		}
		s.updatedDisplay[len(s.digitPins)-1-i] = segment
	}

	return true
}

// Refresh updates the display. Must be called periodically, ideally with >100Hz.
func (s *SevSeg) Refresh() bool {
	if len(s.updatedDisplay) == 0 {
		return false
	}

	s.clearDigitPins()

	// Get segment pattern for current digit
	pattern := s.updatedDisplay[s.currentDigitToRefresh]
	for i, pin := range s.segmentPins {
		if segmentOn := (pattern & (1 << i)) != 0; s.config == CommonCathode {
			if segmentOn {
				pin.High()
			} else {
				pin.Low()
			}
		} else { // CommonAnode
			if segmentOn {
				pin.Low()
			} else {
				pin.High()
			}
		}
	}

	// Turn on the current digit
	if s.currentDigitToRefresh < uint8(len(s.digitPins)) {
		if s.config == CommonCathode {
			s.digitPins[s.currentDigitToRefresh].Low()
		} else {
			s.digitPins[s.currentDigitToRefresh].High()
		}
	}

	s.currentDigitToRefresh = (s.currentDigitToRefresh + 1) % uint8(len(s.digitPins))

	return true
}

// checkAvailableDigits checks if the number can fit within the specified number of digits.
func (s *SevSeg) checkAvailableDigits(number int8) bool {
	count := uint8(0)

	if number == 0 {
		count = 1
	} else {
		if number < 0 {
			count++
			number = -number
		}

		for ; number > 0; number /= 10 {
			count++
		}
	}

	if count > uint8(len(s.digitPins)) {
		return false
	}

	return true
}

// charToSegmentPattern converts a character to its corresponding segment pattern
func (s *SevSeg) charToSegmentPattern(char byte) (uint8, bool) {
	if char >= 'a' && char <= 'z' {
		// Since we can't differ between upper and lower case letters, we convert
		// lower case letters to upper case.
		char = char - 'a' + 'A'
	}

	switch {
	case char >= '0' && char <= '9':
		return digitCodeMap[char-'0'], true
	case char >= 'A' && char <= 'Z':
		return digitCodeMap[char-'A'+10], true
	case char == ' ':
		return digitCodeMap[36], true
	case char == '-':
		return digitCodeMap[37], true
	case char == '.':
		return digitCodeMap[38], true
	case char == '*':
		return digitCodeMap[39], true
	case char == '_':
		return digitCodeMap[40], true
	}

	return 0, false
}

// clearDigitPins turns off all digit pins
func (s *SevSeg) clearDigitPins() {
	for _, pin := range s.digitPins {
		if s.config == CommonCathode {
			pin.High()
		} else {
			pin.Low()
		}
	}
}

// clearSegmentPins turns off all segment pins
func (s *SevSeg) clearSegmentPins() {
	for _, pin := range s.segmentPins {
		if s.config == CommonCathode {
			pin.Low()
		} else {
			pin.High()
		}
	}
}
