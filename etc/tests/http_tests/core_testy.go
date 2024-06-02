package http_tests

import (
	"net/http"
	"testing"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"

	"github.com/stretchr/testify/assert"
)

type Testy interface {
	Prep(tcID string) string
	Run(pairs ...interface{}) (int, string, http.Header)

	AssertStatus(code int)
	AssertHeaders()

	Clean()
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

func NewTesty(t *testing.T, endpoint, path string) Testy {
	testID := newID()
	runApp, cleanup := app.NewApp()
	runApp()
	time.Sleep(1 * time.Second)

	return &testy{
		t:        t,
		id:       "Tc" + testID,
		endpoint: endpoint,
		url:      "http://localhost" + core.HTTPPort + path,
		cleanup:  cleanup,
	}
}

func (tty *testy) Prep(testCase string) string {
	txxID := tty.id + tty.endpoint + testCase
	return txxID
}

func (tty *testy) Run(pairs ...interface{}) (int, string, http.Header) {
	request, err := JSONFromPairs(pairs...)
	if err != nil {
		tty.t.Error(err)
	}
	tty.status, tty.body, tty.headers = POST(tty.t, tty.url, request)
	return tty.status, tty.body, tty.headers
}

func (tty *testy) Clean() {
	tty.cleanup()
}

func (tty *testy) AssertStatus(code int) {
	assert.Equal(tty.t, code, tty.status)
}

func (tty *testy) AssertHeaders() {
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
		got, ok := tty.headers[key]
		if !ok {
			tty.t.Errorf("Header %s is missing", key)
			continue
		}

		for i, exp := range expected {
			if got[i] != exp {
				tty.t.Errorf("Header %s: expected value %s, got %s", key, exp, got[i])
			}
		}
	}
}
