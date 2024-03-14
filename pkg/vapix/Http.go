package vapix

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Cacsjep/goxis/pkg/shared"
)

// on success the ResponseReader must be closed by user
func VapixGet(creds *shared.VapixCreds, url string) RequestResult {

	if creds == nil {
		return RequestResult{IsOk: false, Error: errors.New("Credentials not set eg nil")}
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return RequestResult{IsOk: false, Error: fmt.Errorf("Error creating request: %s", err.Error())}
	}

	req.SetBasicAuth(creds.Username, creds.Password)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return RequestResult{IsOk: false, Error: fmt.Errorf("Error executing request: %s", err.Error())}
	}

	if resp.StatusCode != 200 {
		return RequestResult{IsOk: false, Error: fmt.Errorf("Request not successfull, status code: %d", resp.StatusCode), StatusCode: resp.StatusCode}
	}

	return RequestResult{IsOk: true, ResponseReader: resp.Body, StatusCode: resp.StatusCode}
}

// on success the ResponseReader must be closed by user
func VapixPost(creds *shared.VapixCreds, url string, data interface{}) RequestResult {
	if creds == nil {
		return RequestResult{IsOk: false, Error: errors.New("credentials not set e.g., nil")}
	}

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
	req.SetBasicAuth(creds.Username, creds.Password)

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
