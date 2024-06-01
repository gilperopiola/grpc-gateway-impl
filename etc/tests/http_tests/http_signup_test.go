package http_tests

import (
	"net/http"
	"testing"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
)

type Testy interface {
	Prep(tcID string) string
	Run(request string) (int, string, http.Header)
	AssertHeaders()
}

type testy struct {
	t        *testing.T
	id       string
	endpoint string
	url      string
	cleanup  func()

	status  int
	body    string
	headers http.Header
}

func NewTesty(t *testing.T, endpoint, path string) (Testy, string) {
	id := newID()

	runApp, cleanup := app.NewApp()
	runApp()

	time.Sleep(1 * time.Second)

	url := "http://localhost" + core.HTTPPort
	return &testy{
		t:        t,
		id:       id,
		endpoint: endpoint,
		url:      url + path,
		cleanup:  cleanup,
	}, id
}

func (t *testy) Prep(tcID string) string {
	return tcID + t.id
}

func (t *testy) Run(request string) (int, string, http.Header) {
	t.status, t.body, t.headers = POST(t.t, t.url, request)
	return t.status, t.body, t.headers
}

func (t *testy) AssertHeaders() {
	expecteds := map[string][]string{
		"Access-Control-Allow-Headers": {"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"},
		"Access-Control-Allow-Methods": {"POST, GET, OPTIONS, PUT, DELETE"},
		"Access-Control-Allow-Origin":  {"*"},
		"Content-Security-Policy":      {"default-src 'self'"},
		"Content-Type":                 {"application/json"},
		"X-Content-Type-Options":       {"nosniff"},
		"Strict-Transport-Security":    {"max-age=31536000; includeSubDomains; preload"},
		"X-Frame-Options":              {"SAMEORIGIN"},
		"X-Xss-Protection":             {"1; mode=block"},
	}

	for key, expected := range expecteds {
		got, ok := t.headers[key]
		if !ok {
			t.t.Errorf("Header %s is missing", key)
			continue
		}

		for i, exp := range expected {
			if got[i] != exp {
				t.t.Errorf("Header %s: expected value %s, got %s", key, exp, got[i])
			}
		}
	}
}

func TestHTTPSignup(t *testing.T) {

	// -> 🏠 Prepare
	testy, id := NewTesty(t, "/v1/auth/signup")

	// -> 🚀 Act
	status, body, headers := testy.Run(JSON("username", id, "password", "password"))

	// -> 📡 Assert
	if status != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", status, body)
	}

	assertHeaders(t, headers)

	// -> 🧹 defer can s*ck my di*k
	clean()
}
