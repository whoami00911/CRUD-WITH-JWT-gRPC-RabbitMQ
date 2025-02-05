package hasher

import (
	"crypto/sha1"
	"fmt"
)

type Hash struct {
	salt string
}

func NewHashInit(salt string) *Hash {
	return &Hash{
		salt: salt,
	}
}
func (h *Hash) GenHashPass(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(h.salt)))
}
