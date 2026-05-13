package helper

// ToPtr returns a pointer to the given value
func ToPtr[T any](v T) *T {
	return &v
}

// FromPtrOrZero dereferences a pointer, returning the zero value of T if nil.
func FromPtrOrZero[T any](v *T) T {
	if v == nil {
		var zero T
		return zero
	}
	return *v
}
