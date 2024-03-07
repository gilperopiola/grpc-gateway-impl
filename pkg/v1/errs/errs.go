package errs

// Note these are just error messages and not actual errors.

const (

	// Fatal error messages.
	FatalErrMsgGettingWorkingDir      = "Failed to get working directory: %v"
	FatalErrMsgCreatingProtoValidator = "Failed to create proto validator: %v"
	FatalErrMsgCreatingLogger         = "Failed to create logger: %v"
	FatalErrMsgLoadingTLSCredentials  = "Failed to load server TLS credentials: %v"
	FatalErrMsgReadingTLSCert         = "Failed to read TLS certificate: %v"
	FatalErrMsgAppendingTLSCert       = "Failed to append TLS certificate"
	FatalErrMsgStartingGRPC           = "Failed to start gRPC server: %v"
	FatalErrMsgServingGRPC            = "Failed to serve gRPC server: %v"
	FatalErrMsgStartingHTTP           = "Failed to start HTTP server: %v"
	FatalErrMsgServingHTTP            = "Failed to serve HTTP server: %v"
	FatalErrMsgShuttingDownHTTP       = "Failed to shutdown HTTP server: %v"

	// Request lifecycle error messages.
	ErrMsgInValidation           = "validation error: %v."
	ErrMsgInValidationRuntime    = "unexpected runtime validation error: %v."
	ErrMsgInValidationUnexpected = "unexpected validation error: %v."
	ErrMsgPanic                  = "unexpected panic, something went wrong."
	ErrMsgRateLimitExceeded      = "rate limit exceeded, try again later."

	// HTTP error response bodies.
	// These strings are the JSON representations of a middleware.httpErrorResponseBody.
	// They are what gets sent as the HTTP Response's body when an error occurs.
	HTTPNotFoundErrBody       = `{"error": "not found, check the docs for the correct path and method."}`
	HTTPUnauthorizedErrBody   = `{"error": "unauthorized, authenticate first."}`
	HTTPForbiddenErrBody      = `{"error": "forbidden, not allowed to access this resource."}`
	HTTPInternalErrBody       = `{"error": "internal server error, something went wrong on our end."}`
	HTTPServiceUnavailErrBody = `{"error": "service unavailable, try again later."}`
)
