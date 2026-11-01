package hase

type realNumber interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// conversionInputParameters is a centi or deci parameter.
type conversionInputParameters interface {
	DeciParameter | CentiParameter
}

// parameterTarget is a numeric type used on either side of a pointer conversion.
// Unsigned types are for values that are known to be non-negative.
type parameterTarget interface {
	~int | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}
