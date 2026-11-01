package hase

import (
	"math"
	"strconv"
	"testing"
)

// Benchmarks read inputs from package-level slices and store results in sinks.
// Constants and discarded results let the compiler delete the work.

const pairDigits = "00010203040506070809" +
	"10111213141516171819" +
	"20212223242526272829" +
	"30313233343536373839" +
	"40414243444546474849" +
	"50515253545556575859" +
	"60616263646566676869" +
	"70717273747576777879" +
	"80818283848586878889" +
	"90919293949596979899"

var (
	sinkParameterString string

	benchDeci  = []DeciParameter{0, 1, -1, 54, -123, 210, 32767, -32768}
	benchCenti = []CentiParameter{0, 1, -1, 54, -1234, 2100, 32767, -32768}
)

func deciFormatFloat(d DeciParameter) string {
	return strconv.FormatFloat(float64(d)/10, 'f', 1, 64)
}

func centiFormatFloat(c CentiParameter) string {
	return strconv.FormatFloat(float64(c)/100, 'f', 2, 64)
}

func deciDigits(d DeciParameter) string {
	n := int32(d)
	neg := n < 0
	if neg {
		n = -n
	}
	frac := byte(n%10) + '0'
	whole := uint32(n / 10)

	var buf [8]byte
	i := 0
	if neg {
		buf[0] = '-'
		i = 1
	}
	i = putDigits(buf[:], i, whole)
	buf[i] = '.'
	buf[i+1] = frac
	return string(buf[:i+2])
}

func centiDigits(c CentiParameter) string {
	n := int32(c)
	neg := n < 0
	if neg {
		n = -n
	}
	frac := uint32(n % 100)
	whole := uint32(n / 100)

	var buf [8]byte
	i := 0
	if neg {
		buf[0] = '-'
		i = 1
	}
	i = putDigits(buf[:], i, whole)
	buf[i] = '.'
	buf[i+1] = byte(frac/10) + '0'
	buf[i+2] = byte(frac%10) + '0'
	return string(buf[:i+3])
}

func putDigits(buf []byte, i int, whole uint32) int {
	if whole == 0 {
		buf[i] = '0'
		return i + 1
	}
	var tmp [4]byte
	n := 0
	for whole > 0 {
		tmp[n] = byte(whole%10) + '0'
		n++
		whole /= 10
	}
	for n > 0 {
		n--
		buf[i] = tmp[n]
		i++
	}
	return i
}

func deciPairs(d DeciParameter) string {
	n := int32(d)
	neg := n < 0
	if neg {
		n = -n
	}
	frac := byte(n%10) + '0'
	whole := uint32(n / 10)

	var buf [8]byte
	i := 0
	if neg {
		buf[0] = '-'
		i = 1
	}
	i = putPairs(buf[:], i, whole)
	buf[i] = '.'
	buf[i+1] = frac
	return string(buf[:i+2])
}

func centiPairs(c CentiParameter) string {
	n := int32(c)
	neg := n < 0
	if neg {
		n = -n
	}
	frac := uint32(n % 100)
	whole := uint32(n / 100)

	var buf [8]byte
	i := 0
	if neg {
		buf[0] = '-'
		i = 1
	}
	i = putPairs(buf[:], i, whole)
	p := frac * 2
	buf[i] = '.'
	buf[i+1] = pairDigits[p]
	buf[i+2] = pairDigits[p+1]
	return string(buf[:i+3])
}

func putPairs(buf []byte, i int, whole uint32) int {
	switch {
	case whole < 10:
		buf[i] = byte(whole) + '0'
		return i + 1
	case whole < 100:
		return putPair(buf, i, whole)
	case whole < 1000:
		buf[i] = byte(whole/100) + '0'
		return putPair(buf, i+1, whole%100)
	default:
		i = putPair(buf, i, whole/100)
		return putPair(buf, i, whole%100)
	}
}

func putPair(buf []byte, i int, v uint32) int {
	p := v * 2
	buf[i] = pairDigits[p]
	buf[i+1] = pairDigits[p+1]
	return i + 2
}

func TestParameterDigitsMatchFormatFloat(t *testing.T) {
	for n := int32(math.MinInt16); n <= math.MaxInt16; n++ {
		d := DeciParameter(n)
		want := deciFormatFloat(d)
		if got := deciDigits(d); got != want {
			t.Fatalf("deci digits %d = %q, FormatFloat = %q", n, got, want)
		}
		if got := deciPairs(d); got != want {
			t.Fatalf("deci pairs %d = %q, FormatFloat = %q", n, got, want)
		}
		if got := d.String(); got != want {
			t.Fatalf("deci production %d = %q, FormatFloat = %q", n, got, want)
		}

		c := CentiParameter(n)
		want = centiFormatFloat(c)
		if got := centiDigits(c); got != want {
			t.Fatalf("centi digits %d = %q, FormatFloat = %q", n, got, want)
		}
		if got := centiPairs(c); got != want {
			t.Fatalf("centi pairs %d = %q, FormatFloat = %q", n, got, want)
		}
		if got := c.String(); got != want {
			t.Fatalf("centi production %d = %q, FormatFloat = %q", n, got, want)
		}
	}
}

func BenchmarkDeciString_Production(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkParameterString = benchDeci[i&7].String()
	}
}

func BenchmarkDeciString_FormatFloat(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkParameterString = deciFormatFloat(benchDeci[i&7])
	}
}

func BenchmarkDeciString_Digits(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkParameterString = deciDigits(benchDeci[i&7])
	}
}

func BenchmarkDeciString_Pairs(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkParameterString = deciPairs(benchDeci[i&7])
	}
}

func BenchmarkCentiString_Production(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkParameterString = benchCenti[i&7].String()
	}
}

func BenchmarkCentiString_FormatFloat(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkParameterString = centiFormatFloat(benchCenti[i&7])
	}
}

func BenchmarkCentiString_Digits(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkParameterString = centiDigits(benchCenti[i&7])
	}
}

func BenchmarkCentiString_Pairs(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkParameterString = centiPairs(benchCenti[i&7])
	}
}
