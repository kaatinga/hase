package hase

import "math"

type CentiParameter int16

// NewCentiParameter converts a real number to a centi-parameter.
func NewCentiParameter[T realNumber](value T) CentiParameter {
	return CentiParameter(math.Round(float64(value) * 100))
}

// String converts a centi-parameter to a string with two decimal places.
func (c CentiParameter) String() string {
	n := int32(c)
	neg := n < 0
	if neg {
		n = -n
	}

	var buf [8]byte
	i := 0
	if neg {
		buf[0] = '-'
		i = 1
	}
	i = writeDecimalPairs(buf[:], i, uint32(n/100))
	p := uint32(n%100) * 2
	buf[i] = '.'
	buf[i+1] = decimalPairs[p]
	buf[i+2] = decimalPairs[p+1]
	return string(buf[:i+3])
}

// Float converts a centi-parameter to any float.
func (c CentiParameter) Float[O ~float32 | ~float64]() O {
	return O(c) / 100
}
