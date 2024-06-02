package http_tests

import (
	"net/http"
	"testing"
)

func TestHTTPSignup(t *testing.T) {
	testy := NewTesty(t, "Signup", "/v1/auth/signup")

	for _, tc := range []struct {
		name     string
		password any
		status   int
	}{
		{name: "OK", password: "password", status: http.StatusOK}, // Signup: OK

		{name: "UsrnmExists", password: "password", status: http.StatusOK}, // Signup: Username Exists
		{name: "UsrnmExists", password: "password", status: http.StatusConflict},

		{name: "PwdTooShrt", password: "pass", status: http.StatusBadRequest}, // Signup:x Password Too Short
	} {

		// -> 🏠 Prepare
		txxID := testy.Prep(tc.name)

		// -> 🚀 Act
		testy.Run("username", txxID, "password", tc.password)

		// -> 📡 Assert
		testy.AssertStatus(tc.status)
		testy.AssertHeaders()

		// -> 🧹 defer can s*ck my di*k
		testy.Clean()
	}
}
