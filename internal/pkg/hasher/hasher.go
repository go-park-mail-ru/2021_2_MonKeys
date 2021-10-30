package hasher

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"

	"golang.org/x/crypto/argon2"
)

func HashAndSalt(salt []byte, plainPassword string) (string, error) {
	if salt == nil {
		salt = make([]byte, 8)
		_, err := rand.Read(salt)
		if err != nil {
			return "", err
		}
	}
	hashedPass := argon2.IDKey([]byte(plainPassword), salt, 1, 64*1024, 4, 32)
	saltAndHash := append(salt, hashedPass...)
	return string(saltAndHash[:]), nil
}

func CheckWithHash(hashedStr string, plainStr string) bool {
	salt := []byte(hashedStr[0:8])
	plainStrWithHash, _ := HashAndSalt(salt, plainStr)
	return plainStrWithHash == hashedStr
}

func GetSha1(value []byte) string {
	sum := sha1.Sum(value)
	return hex.EncodeToString(sum[:])
}
