package tools

import (
	"crypto/sha256"
	"encoding/base64"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
)

var _ core.PwdHasher = (*pwdHasher)(nil)

// pwdHasher is our concrete implementation of the core.PwdHasher interface.
type pwdHasher struct {
	salt string
}

// NewPwdHasher returns a new instance of the pwdHasher.
func NewPwdHasher(salt string) *pwdHasher {
	return &pwdHasher{salt: salt}
}

// HashPassword returns a base64 encoded sha256 hash of the pwd + salt.
func (p *pwdHasher) HashPassword(pwd string) string {
	hasher := sha256.New()
	hasher.Write([]byte(pwd + p.salt))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

// ComparePwdsPasswords returns true if plainPwd hashed is equal to the hashedPwd.
func (p *pwdHasher) PasswordsMatch(plain, hashed string) bool {
	return p.HashPassword(plain) == hashed
}
