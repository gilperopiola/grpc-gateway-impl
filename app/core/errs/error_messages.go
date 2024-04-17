package errs

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*              - Errors -             */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// -> These are just strings, not actual errors <-

// T0D0 put Err at the end of the vars name. Also perfect the service/db errors thingy.

const (

	// Fatal error messages. Most of them are used when initializing the application.
	FatalErrMsgGettingWorkingDir = "Failed to get working directory: %v"
	FatalErrMsgCreatingValidator = "Failed to create validator: %v"
	FatalErrMsgCreatingLogger    = "Failed to create logger: %v"
	FatalErrMsgLoadingTLSCreds   = "Failed to load server TLS credentials: %v"
	FatalErrMsgReadingTLSCert    = "Failed to read TLS certificate: %v"
	FatalErrMsgAppendingTLSCert  = "Failed to append TLS certificate"
	FatalErrMsgStartingGRPC      = "Failed to start gRPC server: %v"
	FatalErrMsgServingGRPC       = "Failed to serve gRPC server: %v"
	FatalErrMsgStartingHTTP      = "Failed to start HTTP server: %v"
	FatalErrMsgServingHTTP       = "Failed to serve HTTP server: %v"
	FatalErrMsgShuttingDownHTTP  = "Failed to shutdown HTTP server: %v"
	FatalErrMsgConnectingDB      = "Failed to connect to the database: %v"

	// Non-fatal initialization error messages.
	InsertingDBAdmin    = "Failed to insert admin user to the database: %v"
	ErrMsgGettingDBConn = "Failed to get database connection: %v"

	// Non-fatal shutdown error messages.
	ErrMsgGettingSqlDB = "Failed to get SQL database connection: %v"

	// Request lifecycle error messages.
	ErrMsgInValidation           = "validation error: %v."
	ErrMsgInValidationRuntime    = "runtime validation error: %v."
	ErrMsgInValidationUnexpected = "unexpected validation error: %v."
	ErrMsgPanic                  = "unexpected panic, something went wrong."
	ErrMsgRateLimitExceeded      = "too many requests in a very short time, try again later."

	// Storage Layer error messages.
	ErrMsgRepoNoQueryOpts = "external error -> no query options"

	ErrMsgRepoCreatingUser  = "external error -> creating user"
	ErrMsgRepoGettingUser   = "external error -> getting user"
	ErrMsgRepoGettingUsers  = "external error -> getting users"
	ErrMsgRepoCountingUsers = "external error -> counting users"

	// HTTP error response bodies.
	// They are what gets sent as the HTTP Response's body when an error occurs.
	// Bad Requests (400) are handled by the HTTP Error Handler.
	HTTPNotFoundErrBody       = `{"error": "resource not found."}`
	HTTPUnauthorizedErrBody   = `{"error": "unauthorized, authenticate first."}`
	HTTPForbiddenErrBody      = `{"error": "access to this resource is forbidden."}`
	HTTPInternalErrBody       = `{"error": "internal server error, something failed on our end."}`
	HTTPServiceUnavailErrBody = `{"error": "service unavailable, try again later."}`
)
