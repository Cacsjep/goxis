package vapix

import (
	"fmt"
	"io"
	"math/rand"
	"strconv"
)

// https://axiscommunications.github.io/acap-documentation/docs/develop/VAPIX-access-for-ACAP-applications.html
const INTERNAL_VAPIX_ENDPOINT = "http://127.0.0.12"

type RequestResult struct {
	IsOk           bool
	Error          error
	Password       string
	ResponseReader io.ReadCloser
	StatusCode     int
}

type Param struct {
	Key   string
	Value string
}

type VapixApiCall struct {
	ApiVersion string      `json:"apiVersion"`
	Context    string      `json:"context"`
	Method     string      `json:"method"`
	Error      *VapixError `json:"error"`
	Data       struct {
		PropertiesList map[string]interface{} `json:"propertyList"`
	} `json:"data"`
}

type VapixError struct {
	Code    int `json:"code"`
	Message int `json:"message"`
}

func NewVapixBaseMethodCall(method string) *VapixApiCall {
	return &VapixApiCall{
		ApiVersion: "1.0",
		Context:    strconv.Itoa(rand.Intn(1000)),
		Method:     method,
	}
}

// INTERNAL_VAPIX_ENDPOINT = "http://127.0.0.12" + given path
func InternalUrlConstruct(path string) string {
	return INTERNAL_VAPIX_ENDPOINT + path
}

func CheckForVapixError(vap VapixApiCall) error {
	if vap.Error != nil {
		return fmt.Errorf("vapix Api Error, error-code: %s, error-msg: %s", vap.Error.Code, vap.Error.Message)
	}
	return nil
}
