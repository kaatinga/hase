package hase

import (
	"math"
	"strconv"
	"strings"
	"testing"
)

// Benchmarks read inputs from package-level slices and store results in sinks.
// Constants and discarded results let the compiler delete the work.

var (
	sinkString string
	sinkCode   OffsetDeciParameter

	benchCodes = []OffsetDeciParameter{0, 49, 50, 54, 160, 250, 255, 7}
	benchFloat = []float64{5, 10.4, 10.44, 10.46, 22.5, 29.9, 30, 4.9}
	benchInt   = []int{5, 10, 21, 22, 30, 31, 0, 15}
)

const benchDigitPairs = "00010203040506070809" +
	"10111213141516171819" +
	"20212223242526272829" +
	"30313233343536373839" +
	"40414243444546474849" +
	"50515253545556575859" +
	"60616263646566676869" +
	"70717273747576777879" +
	"80818283848586878889" +
	"90919293949596979899"

func formatFloatString(d OffsetDeciParameter) string {
	return strconv.FormatFloat(float64(d)/10+offsetDeciBaseC, 'f', 1, 64)
}

func appendFloatString(d OffsetDeciParameter) string {
	var buf [8]byte
	return string(strconv.AppendFloat(buf[:0], float64(d)/10+offsetDeciBaseC, 'f', 1, 64))
}

// digitsString is the unrolled Byte2Bytes form: one division per digit.
func digitsString(d OffsetDeciParameter) string {
	n := uint(d) + offsetDeciBaseC*10
	frac := byte(n%10) + '0'
	whole := n / 10
	if whole < 10 {
		return string([]byte{byte(whole) + '0', '.', frac})
	}
	return string([]byte{byte(whole/10) + '0', byte(whole%10) + '0', '.', frac})
}

// digitPairsString writes the integer part from a two-digit lookup.
func digitPairsString(d OffsetDeciParameter) string {
	n := uint(d) + offsetDeciBaseC*10
	frac := byte(n%10) + '0'
	whole := n / 10
	if whole < 10 {
		return string([]byte{byte(whole) + '0', '.', frac})
	}
	i := whole * 2
	return string([]byte{benchDigitPairs[i], benchDigitPairs[i+1], '.', frac})
}

type packedDeci struct {
	one  string
	two  string
	flat string
	lens string
	tab  [256]string
}

func buildPackedDeci() packedDeci {
	var one, two, flat, lens strings.Builder
	var tab [256]string
	for code := range 256 {
		s := strconv.FormatFloat(float64(code)/10+offsetDeciBaseC, 'f', 1, 64)
		tab[code] = s
		if code < 50 {
			one.WriteString(s)
		} else {
			two.WriteString(s)
		}
		flat.WriteString(s)
		if len(s) == 3 {
			flat.WriteByte(' ')
		}
		lens.WriteByte(byte(len(s)))
	}
	return packedDeci{one: one.String(), two: two.String(), flat: flat.String(), lens: lens.String(), tab: tab}
}

var packed = buildPackedDeci()

func packedString(d OffsetDeciParameter) string {
	if d < 50 {
		i := int(d) * 3
		return packed.one[i : i+3]
	}
	i := int(d-50) * 4
	return packed.two[i : i+4]
}

func tableString(d OffsetDeciParameter) string {
	return packed.tab[d]
}

func fixedStrideString(d OffsetDeciParameter) string {
	i := int(d) * 4
	return packed.flat[i : i+int(packed.lens[d])]
}

func roundNew(value float64) OffsetDeciParameter {
	return OffsetDeciParameter(math.Round((value - offsetDeciBaseC) * 10))
}

// addHalfNew rounds half away from zero by truncation for magnitudes whose
// +0.5 addition is exact, and falls back to math.Round otherwise.
func addHalfNew(value float64) OffsetDeciParameter {
	scaled := (value - offsetDeciBaseC) * 10
	if scaled >= 0 && scaled < 1<<51 {
		return OffsetDeciParameter(scaled + 0.5)
	}
	return OffsetDeciParameter(math.Round(scaled))
}

func bitsNew(value float64) OffsetDeciParameter {
	return OffsetDeciParameter(roundHalfAway((value - offsetDeciBaseC) * 10))
}

func roundHalfAway(x float64) float64 {
	const (
		mask     = 0x7FF
		bias     = 1023
		shift    = 52
		signMask = 1 << 63
		fracMask = 1<<shift - 1
		oneBits  = uint64(0x3FF0000000000000)
	)
	bits := math.Float64bits(x)
	e := uint(bits>>shift) & mask
	if e < bias {
		bits &= signMask
		if e == bias-1 {
			bits |= oneBits
		}
	} else if e < bias+shift {
		const half = 1 << (shift - 1)
		e -= bias
		bits += half >> e
		bits &^= fracMask >> e
	}
	return math.Float64frombits(bits)
}

// intArithNew matches math.Round for small integers. It disagrees once the
// value is no longer an exact float64, so it is not a substitute.
func intArithNew(value int) OffsetDeciParameter {
	return OffsetDeciParameter((value - offsetDeciBaseC) * 10)
}

func TestOffsetDeciStringApproachesMatchFormatFloat(t *testing.T) {
	for code := 0; code <= 255; code++ {
		d := OffsetDeciParameter(code)
		want := formatFloatString(d)
		checks := []struct {
			name string
			got  string
		}{
			{"digits", digitsString(d)},
			{"digit pairs", digitPairsString(d)},
			{"packed", packedString(d)},
			{"table", tableString(d)},
			{"fixed stride", fixedStrideString(d)},
			{"append float", appendFloatString(d)},
			{"production", d.String()},
		}
		for _, check := range checks {
			if check.got != want {
				t.Fatalf("%s(%d) = %q, FormatFloat = %q", check.name, code, check.got, want)
			}
		}
	}
}

func TestOffsetDeciNewApproachesMatchRound(t *testing.T) {
	check := func(t *testing.T, name string, got, want OffsetDeciParameter) {
		t.Helper()
		if got != want {
			t.Fatalf("%s = %d, Round = %d", name, got, want)
		}
	}

	values := make([]float64, 0, 10029)
	values = append(values,
		0, 4.9, 5, 5.1, 10.4, 10.44, 10.45, 10.46, 22.5, 29.9, 30, 30.6, 31,
		math.Inf(1), math.Inf(-1), math.NaN(), math.Copysign(0, -1),
		-1, 1e-12, 0.5, -0.5, 1.5, -1.5, 255, 1<<40, -(1 << 40), 1e20, -1e20,
	)
	for i := range 10001 {
		values = append(values, float64(i)/100-5)
	}
	for _, value := range values {
		want := roundNew(value)
		check(t, "add half", addHalfNew(value), want)
		check(t, "bits", bitsNew(value), want)
	}
	for value := -1000; value <= 1000; value++ {
		want := NewOffsetDeciParameter(value)
		check(t, "int arith", intArithNew(value), want)
	}
}

func BenchmarkOffsetDeciString_Production(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkString = benchCodes[i&7].String()
	}
}

func BenchmarkOffsetDeciString_FormatFloat(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkString = formatFloatString(benchCodes[i&7])
	}
}

func BenchmarkOffsetDeciString_AppendFloat(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkString = appendFloatString(benchCodes[i&7])
	}
}

func BenchmarkOffsetDeciString_Digits(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkString = digitsString(benchCodes[i&7])
	}
}

func BenchmarkOffsetDeciString_DigitPairs(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkString = digitPairsString(benchCodes[i&7])
	}
}

func BenchmarkOffsetDeciString_Packed(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkString = packedString(benchCodes[i&7])
	}
}

func BenchmarkOffsetDeciString_Table(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkString = tableString(benchCodes[i&7])
	}
}

func BenchmarkOffsetDeciString_FixedStride(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkString = fixedStrideString(benchCodes[i&7])
	}
}

func BenchmarkOffsetDeciNew_Round(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkCode = roundNew(benchFloat[i&7])
	}
}

func BenchmarkOffsetDeciNew_AddHalf(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkCode = addHalfNew(benchFloat[i&7])
	}
}

func BenchmarkOffsetDeciNew_Bits(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkCode = bitsNew(benchFloat[i&7])
	}
}

func BenchmarkOffsetDeciNew_ProductionFloat(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkCode = NewOffsetDeciParameter(benchFloat[i&7])
	}
}

func BenchmarkOffsetDeciNew_ProductionInt(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkCode = NewOffsetDeciParameter(benchInt[i&7])
	}
}

func BenchmarkOffsetDeciNew_IntArith(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkCode = intArithNew(benchInt[i&7])
	}
}
