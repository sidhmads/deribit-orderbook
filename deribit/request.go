package deribit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type httpMethod string

const (
	get    httpMethod = http.MethodGet
	post   httpMethod = http.MethodPost
	PUT    httpMethod = http.MethodPut
	delete httpMethod = http.MethodDelete
	patch  httpMethod = http.MethodPatch
)

type request struct {
	Method  httpMethod
	Url     string
	Headers map[string]string
	Body    []byte
	Output  interface{}
}

func createNewRequest(method httpMethod, url string, headers map[string]string, body []byte, output interface{}) *request {
	return &request{
		Method:  method,
		Url:     url,
		Headers: headers,
		Body:    body,
		Output:  output,
	}
}

func (r *request) sendHTTPRequest() error {
	req, err := http.NewRequest(string(r.Method), r.Url, bytes.NewBuffer(r.Body))
	if err != nil {
		return err
	}

	for key, value := range r.Headers {
		req.Header.Set(key, value)
	}
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: status code %d", resp.StatusCode)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(responseBody, r.Output)
}
