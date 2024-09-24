package gotraceit

import "fmt"

// ErrContextIsNil is returned when a trace is not found
var ErrContextIsNil = fmt.Errorf("context is nil")
