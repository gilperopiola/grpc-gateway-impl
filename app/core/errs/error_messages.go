package errs

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*              - Errors -             */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// NOTE -> These are just strings. Error messages - NOT actual errors.

const (

	/* -~-~-~-~-~-~ Fatal error messages (init/shutdown) ~-~-~-~-~-~- */

	FailedToStartGRPC    = "Failed to start GRPC Server: %v"
	FailedToServeGRPC    = "Failed to serve GRPC Server: %v"
	FailedToStartHTTP    = "Failed to start HTTP Server: %v"
	FailedToServeHTTP    = "Failed to serve HTTP Server: %v"
	FailedToShutdownHTTP = "Failed to shutdown HTTP Server: %v"

	FailedToCreateProtoVal = "Failed to create Proto Validator: %v"
	FailedToCreateLogger   = "Failed to create Logger: %v"
	FailedDBConn           = "Failed to connect to the DB: %v"

	FailedToLoadTLSCreds  = "Failed to load TLS Creds: %v"
	FailedToReadTLSCert   = "Failed to read TLS Cert: %v"
	FailedToAppendTLSCert = "Failed to append TLS Cert"

	/* -~-~-~-~-~ Non-Fatal error messages (init/shutdown) ~-~-~-~-~- */

	FailedToInsertDBAdmin = "Failed to insert admin to DB: %v"
	FailedToGetSQLDB      = "Failed to get SQL DB connection: %v"
	FailedToCloseSQLDB    = "Failed to close SQL DB connection: %v"

	/* -~-~-~-~-~ Repository error messages ~-~-~-~-~- */

	// Group repository errors
	FailedToCreateGroup    = "Failed to create group: %v"
	GroupNotFound          = "Group not found: %v"
	FailedToFetchGroups    = "Failed to fetch groups: %v"
	FailedToAddUserToGroup = "Failed to add user to group: %v"

	// User repository errors
	FailedToCreateUser = "Failed to create user: %v"
	UserNotFound       = "User not found: %v"
	FailedToFetchUsers = "Failed to fetch users: %v"

	// GPT Chat repository errors
	FailedToCreateChat    = "Failed to create chat: %v"
	ChatNotFound          = "Chat not found: %v"
	FailedToCreateMessage = "Failed to create message: %v"
)

const (

	/* -~-~-~-~-~ Requests lifecycle error messages ~-~-~-~-~- */

	PanicMsg       = "unexpected panic, something went wrong."
	RateLimitedMsg = "too many requests in a very short time, try again later."

	// Request Validation Errors
	ValidatingRequest           = "request validation error -> %v." // We don't use %w as grpc.Status fmt.Sprints the error,
	ValidatingRequestRuntime    = "runtime validation error -> %v." // so no error wrapping.
	ValidatingRequestUnexpected = "unexpected validation error -> %v."

	// JWT Validation
	AuthTokenNotFound  = "auth error -> token not found."
	AuthTokenMalformed = "auth error -> token malformed."
	AuthTokenInvalid   = "auth error -> token invalid."
	AuthRoleInvalid    = "auth error -> role invalid."
	AuthRouteInvalid   = "auth error -> route invalid."
	AuthUserIDInvalid  = "auth error -> user id invalid."
	AuthAPIKeyNotFound = "auth error -> api key not found."
	AuthAPIKeyInvalid  = "auth error -> api key invalid."

	// JWT Generation
	AuthGeneratingToken = "error generating token -> %v."

	// DB Errors
	DBNoQueryOpts   = "db error -> no query options"
	DBCreatingUser  = "db error -> creating user"
	DBGettingUser   = "db error -> getting user"
	DBGettingUsers  = "db error -> getting users"
	DBCountingUsers = "db error -> counting users"

	// HTTP Errors
	HTTPRouteNotFound = "route not found, URL usually follows .../v1/service/endpoint format."
	HTTPUnauthorized  = "unauthorized, authenticate first."
	HTTPForbidden     = "access to this resource is forbidden."
	HTTPInternal      = "internal server error, something failed on our end."
	HTTPUnavailable   = "service unavailable, try again later."
	HTTPConflict      = "resource already exists."
)
