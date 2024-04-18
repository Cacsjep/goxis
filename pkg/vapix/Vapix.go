// Package vapix provides functions to interact with Axis camera VAPIX API.
package vapix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
)

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
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewVapixBaseMethodCall(method string) *VapixApiCall {
	return &VapixApiCall{
		ApiVersion: "1.0",
		Context:    strconv.Itoa(rand.Intn(1000)),
		Method:     method,
	}
}

// INTERNAL_VAPIX_ENDPOINT = "http://127.0.0.12" + given path
func InternalVapixUrlPathJoin(path string) string {
	return INTERNAL_VAPIX_ENDPOINT + path
}

func CheckForVapixError(vap VapixApiCall) error {
	if vap.Error != nil {
		return fmt.Errorf("vapix Api Error, error-code: %d, error-msg: %s", vap.Error.Code, vap.Error.Message)
	}
	return nil
}

// on success the ResponseReader must be closed by user
func VapixGet(username, password, url string) RequestResult {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return RequestResult{IsOk: false, Error: fmt.Errorf("cant creating request: %s", err.Error())}
	}

	req.SetBasicAuth(username, password)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return RequestResult{IsOk: false, Error: fmt.Errorf("cant executing request: %s", err.Error())}
	}

	if resp.StatusCode != 200 {
		return RequestResult{IsOk: false, Error: fmt.Errorf("request not successfull, status code: %d", resp.StatusCode), StatusCode: resp.StatusCode}
	}

	return RequestResult{IsOk: true, ResponseReader: resp.Body, StatusCode: resp.StatusCode}
}

// on success the ResponseReader must be closed by user
func VapixPost(username, password, url string, data interface{}) RequestResult {
	// Marshal the data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return RequestResult{IsOk: false, Error: fmt.Errorf("error marshaling data to JSON: %s", err.Error())}
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return RequestResult{IsOk: false, Error: fmt.Errorf("error creating request: %s", err.Error())}
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(username, password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return RequestResult{IsOk: false, Error: fmt.Errorf("error executing request: %s", err.Error())}
	}
	if resp.StatusCode != 200 {
		return RequestResult{IsOk: false, Error: fmt.Errorf("request not successful, status code: %d", resp.StatusCode), StatusCode: resp.StatusCode}
	}
	return RequestResult{IsOk: true, ResponseReader: resp.Body, StatusCode: resp.StatusCode}
}
