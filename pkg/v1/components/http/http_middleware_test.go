package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/components/common"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/errs"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

/* ----------------------------------- */
/*      - HTTP Middleware Tests -      */
/* ----------------------------------- */

func TestHandleHTTPError(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "NotFound error",
			err:            status.Error(codes.NotFound, "resource not found"),
			expectedStatus: http.StatusNotFound,
			expectedBody:   errs.HTTPNotFoundErrBody,
		},
		{
			name:           "Internal error",
			err:            status.Error(codes.Internal, "internal error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   errs.HTTPInternalErrBody,
		},
		{
			name:           "Unauthorized error",
			err:            status.Error(codes.Unauthenticated, "unauthorized"),
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   errs.HTTPUnauthorizedErrBody,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := runtime.NewServeMux()
			marshaller := &runtime.JSONPb{MarshalOptions: protojson.MarshalOptions{UseProtoNames: true}}
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "http://example.com", nil)

			handleHTTPError(context.Background(), mux, marshaller, recorder, request, tt.err)

			if status := recorder.Code; status != tt.expectedStatus {
				t.Errorf("handleHTTPError() status = %v, want %v", status, tt.expectedStatus)
			}

			if body := recorder.Body.String(); body != tt.expectedBody {
				t.Errorf("handleHTTPError() body = %v, want %v", body, tt.expectedBody)
			}
		})
	}
}

/* ----------------------------------- */
/*     - Other Middleware Tests -      */
/* ----------------------------------- */

func TestLogHTTP(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	middleware := MiddlewareWrapper(&common.Logger{Logger: logger})

	called := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	request := httptest.NewRequest(http.MethodGet, "/test", nil)
	recorder := httptest.NewRecorder()

	handler := middleware(nextHandler)
	handler.ServeHTTP(recorder, request)

	if !called {
		t.Errorf("Expected next handler to be called")
	}

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestSetHTTPResponseHeaders(t *testing.T) {
	recorder := httptest.NewRecorder()
	ctx := context.Background()

	err := setHTTPRespHeaders(ctx, recorder, nil)
	if err != nil {
		t.Errorf("modifyHTTPResponseHeaders returned an error: %v", err)
	}

	expectedHeaders := map[string]string{
		"Content-Security-Policy":   "default-src 'self'",
		"X-XSS-Protection":          "1; mode=block",
		"X-Frame-Options":           "SAMEORIGIN",
		"X-Content-Type-Options":    "nosniff",
		"Strict-Transport-Security": "max-age=31536000; includeSubDomains; preload",
	}

	for header, expectedValue := range expectedHeaders {
		if value := recorder.Header().Get(header); value != expectedValue {
			t.Errorf("Header %s = %v, want %v", header, value, expectedValue)
		}
	}

	// Check that the gRPC-related header was removed
	if value := recorder.Header().Get("Grpc-Metadata-Content-Type"); value != "" {
		t.Errorf("Expected Grpc-Metadata-Content-Type header to be removed, got %v", value)
	}
}