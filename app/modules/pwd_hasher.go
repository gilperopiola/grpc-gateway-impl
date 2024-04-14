package modules

import (
	"crypto/sha256"
	"encoding/base64"
)

// pwdHasher is our concrete implementation of the PwdHasher interface.
type pwdHasher struct {
	salt string
}

// NewPwdHasher returns a new instance of the pwdHasher.
func NewPwdHasher(salt string) *pwdHasher {
	return &pwdHasher{salt: salt}
}

// Hash returns a base64 encoded sha256 hash of the pwd + salt.
func (p *pwdHasher) Hash(pwd string) string {
	hasher := sha256.New()
	hasher.Write([]byte(pwd + p.salt))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

// ComparePasswords returns true if plainPwd hashed is equal to the hashedPwd.
func (p *pwdHasher) Compare(plainPwd, hashedPwd string) bool {
	return p.Hash(plainPwd) == hashedPwd
}
