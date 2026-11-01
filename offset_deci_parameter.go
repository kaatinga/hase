package hase

import (
	"math"
	"strconv"
)

// offsetDeciBaseC is the temperature in degrees stored as zero.
const offsetDeciBaseC = 5

// OffsetDeciParameter is an experimental uint8 deci-parameter.
//
// Experimental: the encoding may change. The stored code is tenths of a
// degree above 5, so 5.0 is 0, 10.4 is 54, and 30.0 is 250.
type OffsetDeciParameter uint8

// NewOffsetDeciParameter converts degrees to an offset deci-parameter.
// value is assumed to lie between 5 and 30.
func NewOffsetDeciParameter[T realNumber](value T) OffsetDeciParameter {
	return OffsetDeciParameter(math.Round((float64(value) - offsetDeciBaseC) * 10))
}

// String converts an offset deci-parameter to a string with one decimal place.
func (d OffsetDeciParameter) String() string {
	return strconv.FormatFloat(float64(d)/10+offsetDeciBaseC, 'f', 1, 64)
}

// Float converts an offset deci-parameter to degrees.
func (d OffsetDeciParameter) Float[O ~float32 | ~float64]() O {
	return O(d)/10 + offsetDeciBaseC
}
