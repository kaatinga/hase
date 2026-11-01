package hase

// PointerTo copies v into *Out.
// A nil pointer stays nil. A negative value requires a signed or floating Out.
func PointerTo[Out parameterTarget, In conversionInputParameters](v *In) *Out {
	if v == nil {
		return nil
	}
	return new(Out(*v)) //nolint:gosec // v is assumed to lie within Out, and to be non-negative when Out is unsigned
}

// PointerFrom copies v into *Out.
// v is assumed to already lie within the centi or deci range, including negatives, and is converted as given.
// A nil pointer stays nil.
func PointerFrom[Out conversionInputParameters, In parameterTarget](v *In) *Out {
	if v == nil {
		return nil
	}
	return new(Out(*v)) //nolint:gosec // v is assumed to lie within the centi or deci range
}
