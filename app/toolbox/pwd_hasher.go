package toolbox

import (
	"crypto/sha256"
	"encoding/base64"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
)

var _ core.PwdHasher = (*pwdHasher)(nil)

type pwdHasher struct {
	salt string
}

func NewPwdHasher(salt string) core.PwdHasher {
	return &pwdHasher{salt}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Returns a base64 encoded sha256 hash of the pwd + salt.
func (ph *pwdHasher) HashPassword(pwd string) string {
	hasher := sha256.New()
	hasher.Write([]byte(pwd + ph.salt))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

// Returns true if plainPwd hashed is equal to the hashedPwd.
func (ph *pwdHasher) PasswordsMatch(plain, hashed string) bool {
	return ph.HashPassword(plain) == hashed
}
