package hase

import (
	"math"
	"strconv"
)

type DeciParameter int16

// NewDeciParameter converts a real number to a deci-parameter.
func NewDeciParameter[T realNumber](value T) DeciParameter {
	return DeciParameter(math.Round(float64(value) * 10))
}

// String converts a deci-parameter to a string with one decimal place.
func (d DeciParameter) String() string {
	return strconv.FormatFloat(float64(d)/10, 'f', 1, 64)
}

// Float converts a deci-parameter to any float.
func (d DeciParameter) Float[O ~float32 | ~float64]() O {
	return O(d) / 10
}
