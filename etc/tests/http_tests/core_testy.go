package http_tests

import (
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"

	"github.com/stretchr/testify/assert"
)

type Testy interface {
	Prep(tcName string) string
	Run(pairs ...interface{}) (int, string, http.Header)
	Clean()

	AssertStatus(code int)
	AssertHeaders()
}

type testy struct {
	t        *testing.T
	id       string
	endpoint string
	url      string
	cleanup  func()

	// got
	status  int
	body    string
	headers http.Header
}

func NewTesty(t *testing.T, endpoint, path string) Testy {
	runApp, cleanup := app.Setup()
	runApp()
	time.Sleep(2 * time.Second)

	return &testy{
		t:        t,
		id:       "Tc" + newID() + endpoint,
		endpoint: endpoint,
		url:      "http://localhost" + core.G.HTTPPort + path,
		cleanup:  cleanup,
	}
}

func (tty *testy) Prep(tcName string) string {
	log.Printf("🔰 Testing %s 🔰\n\n", tty.id+tcName)
	return tty.id + tcName
}

func (tty *testy) Run(pairs ...interface{}) (int, string, http.Header) {
	request, err := JSONFromPairs(pairs...)
	if err != nil {
		tty.t.Error(err)
	}
	tty.status, tty.body, tty.headers = POST(tty.t, tty.url, request)
	log.Printf("🔮 GOT %d 🔮 -> %s\n\n", tty.status, tty.body)
	return tty.status, tty.body, tty.headers
}

func (tty *testy) Clean() {
	tty.cleanup()
}

func (tty *testy) AssertStatus(code int) {
	assert.Equal(tty.t, code, tty.status)
}

func (tty *testy) AssertHeaders() {
	for key, expected := range expectedHeaders {
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

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

var expectedHeaders = map[string][]string{
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
