// Package utils provides utility functions for the goxis project.
package utils

import (
	"encoding/base64"
	"fmt"
)

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

// Helper function for basic authentication
func BasicAuthHeader(username, password string) string {
	return "Basic " + base64Encode(fmt.Sprintf("%s:%s", username, password))
}

func base64Encode(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}
