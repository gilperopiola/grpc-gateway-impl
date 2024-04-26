package errs

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*              - Errors -             */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// NOTE -> These are just strings. Error messages, NOT the actual errors.

const (

	/* -~-~-~-~-~-~ Fatal (init/shutdown) ~-~-~-~-~-~- */

	FailedToStartGRPC    = "Failed to start GRPC Server: %v"
	FailedToServeGRPC    = "Failed to serve GRPC Server: %v"
	FailedToStartHTTP    = "Failed to start HTTP Server: %v"
	FailedToServeHTTP    = "Failed to serve HTTP Server: %v"
	FailedToShutdownHTTP = "Failed to shutdown HTTP Server: %v"

	FailedToCreateProtoVal = "Failed to create Proto Validator: %v"
	FailedToCreateLogger   = "Failed to create Logger: %v"
	FailedToConnectToDB    = "Failed to connect to the DB: %v"

	FailedToLoadTLSCreds  = "Failed to load TLS Creds: %v"
	FailedToReadTLSCert   = "Failed to read TLS Cert: %v"
	FailedToAppendTLSCert = "Failed to append TLS Cert"

	/* -~-~-~-~-~ Non-Fatal (init/shutdown) ~-~-~-~-~- */

	FailedToInsertDBAdmin = "Failed to insert admin to DB: %v"
	FailedToGetSQLDB      = "Failed to get SQL DB connection: %v"
	FailedToCloseSQLDB    = "Failed to close SQL DB connection: %v"

	/* -~-~-~-~-~ Requests lifecycle errors ~-~-~-~-~- */

	PanicMsg       = "unexpected panic, something went wrong."
	RateLimitedMsg = "too many requests in a very short time, try again later."

	InReqValidation           = "request validation error -> %v."
	InReqValidationRuntime    = "runtime validation error -> %v."
	InReqValidationUnexpected = "unexpected validation error -> %v."

	DBNoQueryOpts   = "DB error -> no query options"
	DBCreatingUser  = "DB error -> creating user"
	DBGettingUser   = "DB error -> getting user"
	DBGettingUsers  = "DB error -> getting users"
	DBCountingUsers = "DB error -> counting users"

	// These are the actual HTTP error response bodies when an error happens.
	// Bad Requests (400) are handled by the HTTP Error Handler.
	HTTPNotFound     = `{"error": "resource not found."}`
	HTTPUnauthorized = `{"error": "unauthorized, authenticate first."}`
	HTTPForbidden    = `{"error": "access to this resource is forbidden."}`
	HTTPInternal     = `{"error": "internal server error, something failed on our end."}`
	HTTPUnavailable  = `{"error": "service unavailable, try again later."}`
)
