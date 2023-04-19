package externalapi

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"loan-api/errors"
	"net/http"
	"strings"
)

type RequestHeader struct {
	ContentType string `json:"Content-Type"`
	ClientId    string `json:"client-id"`
	Signature   string `json:"signature"`
	RequestTime string `json:"Request-Time"`
}

// GET Method
func Get(url string, header string) (string, error) {
	return RestCall("GET", url, nil, header)
}

// POST Method
func Post(url string, jsonBodyReq string, header string) (string, error) {
	// Format request body to io.Reader
	body := strings.NewReader(string(jsonBodyReq))

	return RestCall("POST", url, body, header)
}

func setHeader(req *http.Request, header string) {
	if header != "" {
		// Declared an empty map interface
		var headerMap map[string]string

		// Unmarshal or Decode the JSON to the interface.
		json.Unmarshal([]byte(header), &headerMap)

		for k, v := range headerMap {
			req.Header.Set(k, v)
		}
	}
}

// Execute the external API call
func RestCall(method string, url string, payload io.Reader, header string) (string, error) {

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return "", errors.InvalidCurlRequest
	}

	// set header parameters
	setHeader(req, header)

	resp, err := client.Do(req)
	if err != nil {
		return "", errors.InvalidCurlRequest
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.InvalidCurlRequest
	}

	// Convert response body to string
	bodyString := string(bodyBytes)
	return bodyString, nil
}
