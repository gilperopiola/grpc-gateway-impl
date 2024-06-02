package http_tests

import (
	"net/http"
	"testing"
)

func TestHTTPSignup(t *testing.T) {
	testy := NewTesty(t, "Signup", "/v1/auth/signup")

	for _, tc := range []struct {
		name     string
		username any
		password any
		status   int
	}{
		{name: "OK", username: "test", password: "password", status: http.StatusOK},
		{name: "UsrnmExists", username: "test", password: "password", status: http.StatusConflict},
		{name: "PwdTooShrt", username: "test", password: "pass", status: http.StatusBadRequest},
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
