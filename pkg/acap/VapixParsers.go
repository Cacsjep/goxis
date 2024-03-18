package acap

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

func JsonResponseParser(responseReader io.ReadCloser) (*VapixApiCall, error) {
	var response VapixApiCall
	content, err := io.ReadAll(responseReader)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(content, &response); err != nil {
		return nil, err
	}
	if err := CheckForVapixError(response); err != nil {
		return nil, err
	}
	return &response, nil
}

func ParseKeyValueRequestBody(body io.Reader) (map[string]string, error) {
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}

	bodyString := string(bodyBytes)
	lines := strings.Split(bodyString, "\n")
	keyValuePairs := make(map[string]string)

	for _, line := range lines {
		if line == "" {
			continue // skip empty lines
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid line (not a key=value pair): %s", line)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		keyValuePairs[key] = value
	}

	return keyValuePairs, nil
}

func ParseUpdateResponse(body io.Reader) error {
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	bodyString := strings.TrimSpace(string(bodyBytes))
	if bodyString == "OK" {
		return nil
	} else if strings.HasPrefix(bodyString, "Error:") {
		return errors.New(bodyString)
	} else {
		return fmt.Errorf("response is neither 'OK' nor an error: %s", bodyString)
	}
}
