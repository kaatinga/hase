package hase

import (
	"math"
	"strconv"
)

type CentiParameter int16

// NewCentiParameter converts a real number to a centi-parameter.
func NewCentiParameter[T realNumber](value T) CentiParameter {
	return CentiParameter(math.Round(float64(value) * 100))
}

// String converts a centi-parameter to a string with two decimal places.
func (c CentiParameter) String() string {
	return strconv.FormatFloat(float64(c)/100, 'f', 2, 64)
}

// Float converts a centi-parameter to any float.
func (c CentiParameter) Float[O ~float32 | ~float64]() O {
	return O(c) / 100
}
