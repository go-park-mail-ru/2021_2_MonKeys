package hasher

import (
	"crypto/sha1"
	"encoding/hex"
	"math/rand"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func HashAndSalt(salt []byte, password string) string {
	if salt == nil {
		salt = make([]byte, 8)
		for i := range salt {
			salt[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
		}
	}
	// hashedPass := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	hashedPass := []byte(GetSha1([]byte(password)))
	saltAndHash := append(salt, hashedPass...)
	return string(saltAndHash[:])
}

func CheckWithHash(hashedStr string, plainStr string) bool {
	salt := []byte(hashedStr[0:8])
	plainStrWithHash := HashAndSalt(salt, plainStr)
	return plainStrWithHash == hashedStr
}

func GetSha1(value []byte) string {
	sum := sha1.Sum(value)
	return hex.EncodeToString(sum[:])
}
