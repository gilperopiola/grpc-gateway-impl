package entities

/* ----------------------------------- */
/*            - Entities -             */
/* ----------------------------------- */
/*
 * Here we store all our custom data structures => our entities.
 * All layers kind of depend on these entities, so we want to keep them in a separate package.
 */

// SignupRequest is the request entity for the Signup API method.
// Validations can be found on the .proto file.
type SignupRequest struct {
	Username string
	Password string
}

// SignupResponse is the response entity for the Signup API method.
// ID would be the user's ID if this were a real-life project.
type SignupResponse struct {
	ID int
}

// LoginRequest is the request entity for the Login API method.
// Validations can be found on the .proto file.
type LoginRequest struct {
	Username string
	Password string
}

// LoginResponse is the response entity for the Login API method.
// Token would be the user's token if this were a real-life project.
type LoginResponse struct {
	Token string
}
