package hase

import (
	"math"
	"strconv"
	"testing"
)

func TestNewDeciParameterScalesByTen(t *testing.T) {
	check := func(t *testing.T, name string, got, want DeciParameter) {
		t.Helper()
		if got != want {
			t.Errorf("%s = %d, want %d", name, got, want)
		}
	}

	check(t, "int", NewDeciParameter(21), 210)
	check(t, "int8", NewDeciParameter(int8(21)), 210)
	check(t, "negative int8", NewDeciParameter(int8(-12)), -120)
	check(t, "int16", NewDeciParameter(int16(21)), 210)
	check(t, "int32", NewDeciParameter(int32(21)), 210)
	check(t, "int64", NewDeciParameter(int64(21)), 210)
	check(t, "uint", NewDeciParameter(uint(21)), 210)
	check(t, "uint8", NewDeciParameter(uint8(21)), 210)
	check(t, "uint16", NewDeciParameter(uint16(21)), 210)
	check(t, "uint32", NewDeciParameter(uint32(21)), 210)
	check(t, "uint64", NewDeciParameter(uint64(21)), 210)
	check(t, "float32", NewDeciParameter(float32(22.5)), 225)
	check(t, "float64", NewDeciParameter(22.5), 225)
	check(t, "zero", NewDeciParameter(0), 0)
	check(t, "negative int", NewDeciParameter(-12), -120)
	check(t, "negative float", NewDeciParameter(-12.5), -125)
	check(t, "largest whole value that fits", NewDeciParameter(3276), 32760)
	check(t, "smallest whole value that fits", NewDeciParameter(-3276), -32760)
	check(t, "largest value that fits", NewDeciParameter(3276.7), 32767)
	check(t, "smallest value that fits", NewDeciParameter(-3276.8), -32768)

	type celsius int
	type degrees float64
	check(t, "defined int", NewDeciParameter(celsius(21)), 210)
	check(t, "defined float", NewDeciParameter(degrees(22.5)), 225)
}

func TestNewDeciParameterRoundsFloatFraction(t *testing.T) {
	// 1/16 is exact in binary. Times 10 it is 0.625, which rounds to 1.
	// 0.29 * 10 is 2.9, which rounds to 3.
	if got := NewDeciParameter(float64(1) / 16); got != 1 {
		t.Errorf("NewDeciParameter(1/16) = %d, want 1", got)
	}
	if got := NewDeciParameter(float32(1) / 16); got != 1 {
		t.Errorf("NewDeciParameter(float32(1)/16) = %d, want 1", got)
	}
	if got := NewDeciParameter(float64(-1) / 16); got != -1 {
		t.Errorf("NewDeciParameter(-1/16) = %d, want -1", got)
	}
	if got := NewDeciParameter(0.29); got != 3 {
		t.Errorf("NewDeciParameter(0.29) = %d, want 3", got)
	}
	if got := NewDeciParameter(-0.29); got != -3 {
		t.Errorf("NewDeciParameter(-0.29) = %d, want -3", got)
	}
}

func TestNewDeciParameterWrapsOutsideInt16(t *testing.T) {
	check := func(t *testing.T, name string, got, want DeciParameter) {
		t.Helper()
		if got != want {
			t.Errorf("%s = %d, want %d", name, got, want)
		}
	}

	check(t, "3276.8", NewDeciParameter(3276.8), -32768)
	check(t, "-3276.9", NewDeciParameter(-3276.9), 32767)
	check(t, "3277", NewDeciParameter(3277), -32766)
	check(t, "+Inf", NewDeciParameter(math.Inf(1)), -1)
	check(t, "-Inf", NewDeciParameter(math.Inf(-1)), 0)
	check(t, "NaN", NewDeciParameter(math.NaN()), 0)
}

func TestDeciParameterFloatScalesByTenth(t *testing.T) {
	for _, tc := range []struct {
		in     DeciParameter
		want32 float32
		want64 float64
	}{
		{in: 0, want32: 0, want64: 0},
		{in: 225, want32: 22.5, want64: 22.5},
		{in: 1, want32: 0.1, want64: 0.1},
		{in: -123, want32: -12.3, want64: -12.3},
		{in: 32760, want32: 3276, want64: 3276},
		{in: -32760, want32: -3276, want64: -3276},
		{in: 32767, want32: 3276.7, want64: 3276.7},
		{in: -32768, want32: -3276.8, want64: -3276.8},
	} {
		if got := tc.in.Float[float32](); got != tc.want32 {
			t.Errorf("DeciParameter(%d).Float[float32]() = %v, want %v", tc.in, got, tc.want32)
		}
		if got := tc.in.Float[float64](); got != tc.want64 {
			t.Errorf("DeciParameter(%d).Float[float64]() = %v, want %v", tc.in, got, tc.want64)
		}
	}
}

func TestDeciParameterFloatInvertsExactInputs(t *testing.T) {
	for _, value := range []float64{0, 0.1, 21, 22.5, -12.3, -12.5, -40, 125, 3276, -3276, 3276.7, -3276.8} {
		if got := NewDeciParameter(value).Float[float64](); got != value {
			t.Errorf("Float[float64](NewDeciParameter(%v)) = %v, want %v", value, got, value)
		}
	}
	for _, value := range []float32{0, 0.1, 21, 22.5, -12.3, -12.5, 3276.7, -3276.8} {
		if got := NewDeciParameter(value).Float[float32](); got != value {
			t.Errorf("Float[float32](NewDeciParameter(%v)) = %v, want %v", value, got, value)
		}
	}
}

func TestDeciParameterString(t *testing.T) {
	for _, tc := range []struct {
		in   DeciParameter
		want string
	}{
		{in: 0, want: "0.0"},
		{in: 210, want: "21.0"},
		{in: 225, want: "22.5"},
		{in: 1, want: "0.1"},
		{in: -123, want: "-12.3"},
		{in: -1, want: "-0.1"},
		{in: 32767, want: "3276.7"},
		{in: -32768, want: "-3276.8"},
	} {
		if got := tc.in.String(); got != tc.want {
			t.Errorf("DeciParameter(%d).String() = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestDeciParameterStringMatchesFormatFloat(t *testing.T) {
	for n := int32(math.MinInt16); n <= math.MaxInt16; n++ {
		d := DeciParameter(n)
		want := strconv.FormatFloat(float64(d)/10, 'f', 1, 64)
		if got := d.String(); got != want {
			t.Fatalf("DeciParameter(%d).String() = %q, want %q", n, got, want)
		}
	}
}
