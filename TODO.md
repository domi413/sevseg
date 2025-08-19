# TODO's

- Readme:
  - mention that Refresh() should not be spammed but called at a reasonable rate
  - mention that brightness-control requires pwm pines
  - update api
  - mention that we can display number up to int32, which fit perfectly for a 8 digit display.

```go
// // ---------------------------------------------------

// //nolint:all
// type (
// Pin uint8
// PinMode uint8
// )

// //nolint:all
// const (
// PinOutput PinMode = iota
// PinInput
// PinInputPullup
// PinInputPulldown
// )

// //nolint:all
// type PinConfig struct{ Mode PinMode }

// //nolint:all
// func (p Pin) Configure(config PinConfig) {}

// //nolint:all
// func (p Pin) High() {}

// //nolint:all
// func (p Pin) Low() {}

// //nolint:all
// var machine = struct {
// Pin Pin
// PinConfig PinConfig
// PinOutput PinMode
// }{PinOutput: PinOutput}

// // ------------------------------------------------------
```
