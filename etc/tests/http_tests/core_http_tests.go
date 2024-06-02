package http_tests

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"testing"

	"github.com/gilperopiola/god"
)

/* -~-~-~-~-~    Test Setup    ~-~-~-~-~-~- */

func newID() string {
	id := ""
	for i := 0; i < 3; i++ {
		id += god.MapIntToLetter(rand.Intn(26))
	}
	return strings.ToUpper(id)
}

/* -~-~-~-~-~    HTTP    ~-~-~-~-~-~- */

func GET(url string) (*http.Response, error) {
	return http.Get(url)
}

func POST(t *testing.T, url, request string) (int, string, http.Header) {
	resp, err := http.Post(url, "application/json", strings.NewReader(request))
	if err != nil {
		t.Fatalf("Failed POST request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed reading POST body: %v", err)
	}

	return resp.StatusCode, string(body), resp.Header
}

/* -~-~-~-~-~    JSON    ~-~-~-~-~-~- */

// JSON creates a JSON string from key-value pairs
func JSONFromPairs(pairs ...interface{}) (string, error) {
	if len(pairs)%2 != 0 {
		return "", fmt.Errorf("pairs must be even, got %d", len(pairs))
	}

	data := make(map[string]interface{})

	for keyIndex := 0; keyIndex < len(pairs); keyIndex += 2 {
		valIndex := keyIndex + 1

		key, ok := pairs[keyIndex].(string)
		if !ok {
			return "", fmt.Errorf("pair keys must be strings, got %T", pairs[keyIndex])
		}

		data[key] = pairs[valIndex]
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal json, %w", err)
	}

	return string(jsonData), nil
}
