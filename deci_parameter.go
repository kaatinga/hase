package hase

import "math"

type DeciParameter int16

// NewDeciParameter converts a real number to a deci-parameter.
func NewDeciParameter[T realNumber](value T) DeciParameter {
	return DeciParameter(math.Round(float64(value) * 10))
}

// String converts a deci-parameter to a string with one decimal place.
func (d DeciParameter) String() string {
	n := int32(d)
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
	i = writeDecimalPairs(buf[:], i, uint32(n/10))
	buf[i] = '.'
	buf[i+1] = byte(n%10) + '0'
	return string(buf[:i+2])
}

// Float converts a deci-parameter to any float.
func (d DeciParameter) Float[O ~float32 | ~float64]() O {
	return O(d) / 10
}
