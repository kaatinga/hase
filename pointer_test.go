package hase

import "testing"

func TestPointerToAndFrom(t *testing.T) {
	if PointerTo[int32, CentiParameter](nil) != nil {
		t.Error("PointerTo(nil centi) is not nil")
	}
	if PointerTo[int32, DeciParameter](nil) != nil {
		t.Error("PointerTo(nil deci) is not nil")
	}
	if PointerFrom[CentiParameter, int16](nil) != nil {
		t.Error("PointerFrom(nil) is not nil")
	}

	for _, value := range []CentiParameter{0, 1, -1, 2250, -1234, 32767, -32768} {
		checkPointer(t, value)
	}
	for _, value := range []DeciParameter{0, 1, -1, 225, -123, 32767, -32768} {
		checkPointer(t, value)
	}
}

func checkPointer[In conversionInputParameters](t *testing.T, value In) {
	t.Helper()
	asInt32 := PointerTo[int32](&value)
	asInt64 := PointerTo[int64](&value)
	asFloat := PointerTo[float64](&value)
	asInt16 := PointerTo[int16](&value)
	if asInt32 == nil || *asInt32 != int32(value) {
		t.Errorf("PointerTo[int32](%d) = %v", value, asInt32)
	}
	if asInt64 == nil || *asInt64 != int64(value) {
		t.Errorf("PointerTo[int64](%d) = %v", value, asInt64)
	}
	if asFloat == nil || *asFloat != float64(value) {
		t.Errorf("PointerTo[float64](%d) = %v", value, asFloat)
	}
	if asInt16 == nil || *asInt16 != int16(value) {
		t.Errorf("PointerTo[int16](%d) = %v", value, asInt16)
	}

	backs := []In{
		*PointerFrom[In](asInt32),
		*PointerFrom[In](asInt64),
		*PointerFrom[In](asFloat),
		*PointerFrom[In](asInt16),
	}
	if value >= 0 {
		asUint16 := PointerTo[uint16](&value)
		asUint32 := PointerTo[uint32](&value)
		if asUint16 == nil || *asUint16 != uint16(value) {
			t.Errorf("PointerTo[uint16](%d) = %v", value, asUint16)
		}
		if asUint32 == nil || *asUint32 != uint32(value) {
			t.Errorf("PointerTo[uint32](%d) = %v", value, asUint32)
		}
		backs = append(backs, *PointerFrom[In](asUint16), *PointerFrom[In](asUint32))
	}
	for _, back := range backs {
		if back != value {
			t.Errorf("PointerFrom(%d) = %d", value, back)
		}
	}
}
