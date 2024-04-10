package utils

// StrPtr returns a pointer to the given string value.
func StrPtr(value string) *string {
	return &value
}

// IntPtr returns a pointer to the given int value.
func IntPtr(value int) *int {
	return &value
}

// BoolPtr returns a pointer to the given int value.
func BoolPtr(value bool) *bool {
	return &value
}

// Float64Ptr returns a pointer to the given int value.
func Float64Ptr(value float64) *float64 {
	return &value
}
