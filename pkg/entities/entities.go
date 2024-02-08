package entities

/** -------------------- *** - Package to store all our custom data structures.
/**     pkg/entities     *** - All layers will depend on this.
/** -------------------- **/

type SignupRequest struct {
	Username string
	Password string
}

type SignupResponse struct {
	ID int
}

type LoginRequest struct {
	Username string
	Password string
}

type LoginResponse struct {
	Token string
}
