package mailtd

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"golang.org/x/crypto/argon2"
)

// DeriveAuthKey derives a 64-char hex auth_key from address and password using
// Argon2id, matching the parameters used by the Mail.td web frontend and backend.
//
//	auth_key = Argon2id(password, salt=SHA256(lower(trim(address))),
//	                    time=3, memory=16384 KiB, threads=1, keyLen=32)
//
// SDK methods that accept a password derive the auth_key locally with this
// function so the password never leaves the client.
func DeriveAuthKey(address, password string) string {
	salt := sha256.Sum256([]byte(strings.ToLower(strings.TrimSpace(address))))
	key := argon2.IDKey([]byte(password), salt[:], 3, 16384, 1, 32)
	return hex.EncodeToString(key)
}
