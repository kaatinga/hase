package hase

// decimalPairs lists "00" through "99", two bytes per value.
const decimalPairs = "00010203040506070809" +
	"10111213141516171819" +
	"20212223242526272829" +
	"30313233343536373839" +
	"40414243444546474849" +
	"50515253545556575859" +
	"60616263646566676869" +
	"70717273747576777879" +
	"80818283848586878889" +
	"90919293949596979899"

// writeDecimalPairs writes the decimal form of whole, which is at most 4 digits.
func writeDecimalPairs(buf []byte, i int, whole uint32) int {
	switch {
	case whole < 10:
		buf[i] = byte(whole) + '0'
		return i + 1
	case whole < 100:
		return writeDecimalPair(buf, i, whole)
	case whole < 1000:
		buf[i] = byte(whole/100) + '0'
		return writeDecimalPair(buf, i+1, whole%100)
	default:
		i = writeDecimalPair(buf, i, whole/100)
		return writeDecimalPair(buf, i, whole%100)
	}
}

func writeDecimalPair(buf []byte, i int, v uint32) int {
	p := v * 2
	buf[i] = decimalPairs[p]
	buf[i+1] = decimalPairs[p+1]
	return i + 2
}
