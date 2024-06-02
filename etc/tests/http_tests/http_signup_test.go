package http_tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"

	"github.com/stretchr/testify/assert"
)

func TestHTTPSignup(t *testing.T) {

	testy := NewTesty(t, "Signup", "/v1/auth/signup")
	defer testy.Clean()

	type testCase struct {
		name     string
		username any // -> By default it's the testID
		password any
		status   int
	}

	for _, tc := range []testCase{
		{
			name: "OK", password: "password", status: http.StatusOK,
		},
		{
			name: "UsrEmpty", username: "", password: "password", status: http.StatusBadRequest,
		},
		{
			name: "PwdEmpty", password: "", status: http.StatusBadRequest,
		},
		{
			name: "PwdNil", status: http.StatusBadRequest,
		},
		{
			name: "UsrTooShrt", username: "Tc", password: "password", status: http.StatusBadRequest, // -> Overrides username
		},
		{
			name: "PwdTooShrt", password: "pwd", status: http.StatusBadRequest,
		},
		{
			name: "UsrExists", password: "password", status: http.StatusOK, // -> Prep for next test
		},
		{
			name: "UsrExists", password: "password", status: http.StatusConflict,
		},
	} {

		// -> 🏠 Prepare
		testID := testy.Prep(tc.name)
		username := testID
		if tmp, ok := tc.username.(string); ok {
			username = tmp
		}

		// -> 🚀 Act
		_, body, _ := testy.Run("username", username, "password", tc.password)

		// -> 📡 Assert
		testy.AssertStatus(tc.status)
		testy.AssertHeaders()

		if tc.status >= 400 {
			continue
		}

		var typedBody pbs.SignupResponse
		if err := json.Unmarshal([]byte(body), &typedBody); err != nil {
			t.Errorf("Error unmarshaling response: %s", body)
		}

		assert.Greater(t, typedBody.Id, int32(0))
	}
}
