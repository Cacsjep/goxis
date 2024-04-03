package utils

// NewStringPointer returns a pointer to the given string value.
func NewStringPointer(value string) *string {
	return &value
}

// NewIntPointer returns a pointer to the given int value.
func NewIntPointer(value int) *int {
	return &value
}

// NewBoolPointer returns a pointer to the given int value.
func NewBoolPointer(value bool) *bool {
	return &value
}

// NewFloat64Pointer returns a pointer to the given int value.
func NewFloat64Pointer(value float64) *float64 {
	return &value
}
