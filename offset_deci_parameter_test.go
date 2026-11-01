package hase

import (
	"math"
	"testing"
)

func TestNewOffsetDeciParameterScalesAboveFive(t *testing.T) {
	check := func(t *testing.T, name string, got, want OffsetDeciParameter) {
		t.Helper()
		if got != want {
			t.Errorf("%s = %d, want %d", name, got, want)
		}
	}

	check(t, "five", NewOffsetDeciParameter(5), 0)
	check(t, "ten point four", NewOffsetDeciParameter(10.4), 54)
	check(t, "thirty", NewOffsetDeciParameter(30), 250)
	check(t, "int", NewOffsetDeciParameter(21), 160)
	check(t, "int8", NewOffsetDeciParameter(int8(21)), 160)
	check(t, "int16", NewOffsetDeciParameter(int16(21)), 160)
	check(t, "int32", NewOffsetDeciParameter(int32(21)), 160)
	check(t, "int64", NewOffsetDeciParameter(int64(21)), 160)
	check(t, "uint", NewOffsetDeciParameter(uint(21)), 160)
	check(t, "uint8", NewOffsetDeciParameter(uint8(21)), 160)
	check(t, "uint16", NewOffsetDeciParameter(uint16(21)), 160)
	check(t, "uint32", NewOffsetDeciParameter(uint32(21)), 160)
	check(t, "uint64", NewOffsetDeciParameter(uint64(21)), 160)
	check(t, "float32", NewOffsetDeciParameter(float32(22.5)), 175)
	check(t, "float64", NewOffsetDeciParameter(22.5), 175)
	check(t, "one tenth above base", NewOffsetDeciParameter(5.1), 1)
	check(t, "one tenth below max", NewOffsetDeciParameter(29.9), 249)

	type celsius int
	type degrees float64
	check(t, "defined int", NewOffsetDeciParameter(celsius(21)), 160)
	check(t, "defined float", NewOffsetDeciParameter(degrees(10.4)), 54)
}

func TestNewOffsetDeciParameterRoundsFloatFraction(t *testing.T) {
	// float32(10.4) is slightly under 10.4. Times 10 above 5 it is just under 54.
	if got := NewOffsetDeciParameter(float32(10.4)); got != 54 {
		t.Errorf("NewOffsetDeciParameter(float32(10.4)) = %d, want 54", got)
	}
	if got := NewOffsetDeciParameter(10.44); got != 54 {
		t.Errorf("NewOffsetDeciParameter(10.44) = %d, want 54", got)
	}
	if got := NewOffsetDeciParameter(10.46); got != 55 {
		t.Errorf("NewOffsetDeciParameter(10.46) = %d, want 55", got)
	}
}

func TestNewOffsetDeciParameterWrapsOutsideUint8(t *testing.T) {
	check := func(t *testing.T, name string, got, want OffsetDeciParameter) {
		t.Helper()
		if got != want {
			t.Errorf("%s = %d, want %d", name, got, want)
		}
	}

	check(t, "30.6", NewOffsetDeciParameter(30.6), 0)
	check(t, "31", NewOffsetDeciParameter(31), 4)
	check(t, "4.9", NewOffsetDeciParameter(4.9), 255)
	check(t, "0", NewOffsetDeciParameter(0), 206)
	check(t, "+Inf", NewOffsetDeciParameter(math.Inf(1)), 255)
	check(t, "-Inf", NewOffsetDeciParameter(math.Inf(-1)), 0)
	check(t, "NaN", NewOffsetDeciParameter(math.NaN()), 0)
}

func TestOffsetDeciParameterFloat(t *testing.T) {
	for _, tc := range []struct {
		in     OffsetDeciParameter
		want32 float32
		want64 float64
	}{
		{in: 0, want32: 5, want64: 5},
		{in: 54, want32: 10.4, want64: 10.4},
		{in: 1, want32: 5.1, want64: 5.1},
		{in: 250, want32: 30, want64: 30},
		{in: 255, want32: 30.5, want64: 30.5},
	} {
		if got := tc.in.Float[float32](); got != tc.want32 {
			t.Errorf("OffsetDeciParameter(%d).Float[float32]() = %v, want %v", tc.in, got, tc.want32)
		}
		if got := tc.in.Float[float64](); got != tc.want64 {
			t.Errorf("OffsetDeciParameter(%d).Float[float64]() = %v, want %v", tc.in, got, tc.want64)
		}
	}
}

func TestOffsetDeciParameterFloatInvertsExactInputs(t *testing.T) {
	for _, value := range []float64{5, 5.1, 10.4, 21, 22.5, 29.9, 30} {
		if got := NewOffsetDeciParameter(value).Float[float64](); got != value {
			t.Errorf("Float[float64](NewOffsetDeciParameter(%v)) = %v, want %v", value, got, value)
		}
	}
	for _, value := range []float32{5, 5.1, 10.4, 21, 22.5, 30} {
		if got := NewOffsetDeciParameter(value).Float[float32](); got != value {
			t.Errorf("Float[float32](NewOffsetDeciParameter(%v)) = %v, want %v", value, got, value)
		}
	}
}

func TestOffsetDeciParameterString(t *testing.T) {
	for _, tc := range []struct {
		in   OffsetDeciParameter
		want string
	}{
		{in: 0, want: "5.0"},
		{in: 54, want: "10.4"},
		{in: 1, want: "5.1"},
		{in: 160, want: "21.0"},
		{in: 250, want: "30.0"},
		{in: 255, want: "30.5"},
	} {
		if got := tc.in.String(); got != tc.want {
			t.Errorf("OffsetDeciParameter(%d).String() = %q, want %q", tc.in, got, tc.want)
		}
	}
}
