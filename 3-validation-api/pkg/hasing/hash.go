package hasing

import (
	"crypto/sha256"
	"encoding/hex"
)

func GetHashString(input string) string {
	hash256 := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash256[:])
}
