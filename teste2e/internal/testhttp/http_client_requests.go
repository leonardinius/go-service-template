package testhttp

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// MustGET performs an HTTP GET request to the specified address and returns the response.
// It is used in tests to make sure the request is successful, otherwise it will fail the test.
func MustGET(ctx context.Context, t *testing.T, addr string) *http.Response {
	t.Helper()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, addr, http.NoBody)
	require.NoError(t, err, "Failed to create request: %v", err)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err, "Failed to perform request: %v", err)

	t.Cleanup(func() {
		_ = resp.Body.Close()
	})

	return resp
}

func MustPost(ctx context.Context, t *testing.T, addr, contentType string, body io.Reader) *http.Response {
	t.Helper()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, addr, body)
	require.NoError(t, err, "Failed to create request: %v", err)
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err, "Failed to perform request: %v", err)

	t.Cleanup(func() {
		_ = resp.Body.Close()
	})

	return resp
}

// MustReadFullyString reads all the data from the specified `io.ReadCloser` and returns it as a string.
// If there is an error while reading or converting to string, it returns an empty string and the error.
func MustReadFullyString(t *testing.T, resp *http.Response) string {
	t.Helper()
	data, err := ReadFullyString(resp)
	require.NoError(t, err, "Failed to read fully: %v", err)
	return data
}

// MustReadFullyJSON reads the response body from an HTTP response and populates the target variable with the JSON data.
// It is used in tests to ensure that the response body is successfully read and parsed as JSON, otherwise it will fail the test.
// The resp parameter is the HTTP response that contains the JSON data to be read.
// The target parameter is the variable where the parsed JSON data will be stored.
// The function returns the target variable with the JSON data.
// If there is an error reading or parsing the JSON data, the function will fail the test and return an error message.
func MustReadFullyJSON[T ~*E, E any](t *testing.T, resp *http.Response, target T) {
	t.Helper()
	require.Contains(t, resp.Header.Get("Content-Type"), "application/json")
	err := ReadFullyJSON(resp, target)
	require.NoError(t, err, "Failed to read fully: %v", err)
}

// ReadFullyJSON reads the response body from an HTTP response and populates the target variable with the JSON data.
// The resp parameter is the HTTP response that contains the JSON data to be read.
// The target parameter is the variable where the parsed JSON data will be stored.
// The function returns an error if there is a problem reading or parsing the JSON data.
func ReadFullyJSON[T ~*E, E any](resp *http.Response, target T) (err error) {
	defer func() {
		_ = resp.Body.Close()
	}()

	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(target)
}

// ReadFullyString reads all the data from the specified `io.ReadCloser` and returns it as a string.
// If there is an error while reading or converting to string, it returns an empty string and the error.
func ReadFullyString(resp *http.Response) (string, error) {
	data, err := ReadFully(resp)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// ReadFully reads all the data from the specified `io.ReadCloser` and returns it as a byte slice.
// If there is an error while reading, it returns `nil` and the error.
func ReadFully(resp *http.Response) ([]byte, error) {
	defer func() {
		_ = resp.Body.Close()
	}()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
