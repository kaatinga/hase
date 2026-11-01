package hase

import (
	"math"
	"testing"
)

func TestNewCentiParameterScalesByHundred(t *testing.T) {
	check := func(t *testing.T, name string, got, want CentiParameter) {
		t.Helper()
		if got != want {
			t.Errorf("%s = %d, want %d", name, got, want)
		}
	}

	check(t, "int", NewCentiParameter(21), 2100)
	check(t, "int8", NewCentiParameter(int8(21)), 2100)
	check(t, "negative int8", NewCentiParameter(int8(-12)), -1200)
	check(t, "int16", NewCentiParameter(int16(21)), 2100)
	check(t, "int32", NewCentiParameter(int32(21)), 2100)
	check(t, "int64", NewCentiParameter(int64(21)), 2100)
	check(t, "uint", NewCentiParameter(uint(21)), 2100)
	check(t, "uint8", NewCentiParameter(uint8(21)), 2100)
	check(t, "uint16", NewCentiParameter(uint16(21)), 2100)
	check(t, "uint32", NewCentiParameter(uint32(21)), 2100)
	check(t, "uint64", NewCentiParameter(uint64(21)), 2100)
	check(t, "float32", NewCentiParameter(float32(22.5)), 2250)
	check(t, "float64", NewCentiParameter(22.5), 2250)
	check(t, "zero", NewCentiParameter(0), 0)
	check(t, "negative int", NewCentiParameter(-12), -1200)
	check(t, "negative float", NewCentiParameter(-12.5), -1250)
	check(t, "largest whole value that fits", NewCentiParameter(327), 32700)
	check(t, "smallest whole value that fits", NewCentiParameter(-327), -32700)
	check(t, "largest value that fits", NewCentiParameter(327.67), 32767)
	check(t, "smallest value that fits", NewCentiParameter(-327.68), -32768)

	type celsius int
	type degrees float64
	check(t, "defined int", NewCentiParameter(celsius(21)), 2100)
	check(t, "defined float", NewCentiParameter(degrees(22.5)), 2250)
}

func TestNewCentiParameterRoundsFloatFraction(t *testing.T) {
	// 1/64 is exact in binary. Times 100 it is 1.5625, which rounds to 2.
	// 0.29 * 100 is 28.999999999999996 in float64, so rounding is what keeps 29.
	if got := NewCentiParameter(float64(1) / 64); got != 2 {
		t.Errorf("NewCentiParameter(1/64) = %d, want 2", got)
	}
	if got := NewCentiParameter(float32(1) / 64); got != 2 {
		t.Errorf("NewCentiParameter(float32(1)/64) = %d, want 2", got)
	}
	if got := NewCentiParameter(float64(-1) / 64); got != -2 {
		t.Errorf("NewCentiParameter(-1/64) = %d, want -2", got)
	}
	if got := NewCentiParameter(0.29); got != 29 {
		t.Errorf("NewCentiParameter(0.29) = %d, want 29", got)
	}
	if got := NewCentiParameter(-0.29); got != -29 {
		t.Errorf("NewCentiParameter(-0.29) = %d, want -29", got)
	}
}

func TestNewCentiParameterWrapsOutsideInt16(t *testing.T) {
	check := func(t *testing.T, name string, got, want CentiParameter) {
		t.Helper()
		if got != want {
			t.Errorf("%s = %d, want %d", name, got, want)
		}
	}

	check(t, "327.68", NewCentiParameter(327.68), -32768)
	check(t, "-327.69", NewCentiParameter(-327.69), 32767)
	check(t, "328", NewCentiParameter(328), -32736)
	check(t, "+Inf", NewCentiParameter(math.Inf(1)), -1)
	check(t, "-Inf", NewCentiParameter(math.Inf(-1)), 0)
	check(t, "NaN", NewCentiParameter(math.NaN()), 0)
}

func TestCentiParameterFloatScalesByHundredth(t *testing.T) {
	for _, tc := range []struct {
		in     CentiParameter
		want32 float32
		want64 float64
	}{
		{in: 0, want32: 0, want64: 0},
		{in: 2250, want32: 22.5, want64: 22.5},
		{in: 1, want32: 0.01, want64: 0.01},
		{in: -1234, want32: -12.34, want64: -12.34},
		{in: 32700, want32: 327, want64: 327},
		{in: -32700, want32: -327, want64: -327},
		{in: 32767, want32: 327.67, want64: 327.67},
		{in: -32768, want32: -327.68, want64: -327.68},
	} {
		if got := tc.in.Float[float32](); got != tc.want32 {
			t.Errorf("CentiParameter(%d).Float[float32]() = %v, want %v", tc.in, got, tc.want32)
		}
		if got := tc.in.Float[float64](); got != tc.want64 {
			t.Errorf("CentiParameter(%d).Float[float64]() = %v, want %v", tc.in, got, tc.want64)
		}
	}
}

func TestCentiParameterFloatInvertsExactInputs(t *testing.T) {
	for _, value := range []float64{0, 0.01, 21, 22.5, -12.34, -12.5, -40, 125, 327, -327, 327.67, -327.68} {
		if got := NewCentiParameter(value).Float[float64](); got != value {
			t.Errorf("Float[float64](NewCentiParameter(%v)) = %v, want %v", value, got, value)
		}
	}
	for _, value := range []float32{0, 0.01, 21, 22.5, -12.34, -12.5, 327.67, -327.68} {
		if got := NewCentiParameter(value).Float[float32](); got != value {
			t.Errorf("Float[float32](NewCentiParameter(%v)) = %v, want %v", value, got, value)
		}
	}
}

func TestCentiParameterString(t *testing.T) {
	for _, tc := range []struct {
		in   CentiParameter
		want string
	}{
		{in: 0, want: "0.00"},
		{in: 2100, want: "21.00"},
		{in: 2250, want: "22.50"},
		{in: 1, want: "0.01"},
		{in: -1234, want: "-12.34"},
		{in: -1, want: "-0.01"},
		{in: 32767, want: "327.67"},
		{in: -32768, want: "-327.68"},
	} {
		if got := tc.in.String(); got != tc.want {
			t.Errorf("CentiParameter(%d).String() = %q, want %q", tc.in, got, tc.want)
		}
	}
}
