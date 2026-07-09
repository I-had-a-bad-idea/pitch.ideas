package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	timeCost    uint32 = 3
	memoryCost  uint32 = 64 * 1024 // 64 MB
	parallelism uint8  = 4
	saltLength         = 16
	keyLength          = 32
)

// HashPassword hashes a password using Argon2id.
func HashPassword(password string) (string, error) {
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		timeCost,
		memoryCost,
		parallelism,
		keyLength,
	)

	b64 := base64.RawStdEncoding

	return fmt.Sprintf(
		"$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		memoryCost,
		timeCost,
		parallelism,
		b64.EncodeToString(salt),
		b64.EncodeToString(hash),
	), nil
}

// VerifyPassword returns true if password matches the stored hash.
func VerifyPassword(password, encodedHash string) bool {
	ok, err := verifyPassword(password, encodedHash)
	return err == nil && ok
}

func verifyPassword(password, encodedHash string) (bool, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, errors.New("invalid hash")
	}

	var memory uint32
	var time uint32
	var threads uint8

	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
	if err != nil {
		return false, err
	}

	b64 := base64.RawStdEncoding

	salt, err := b64.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	hash, err := b64.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	otherHash := argon2.IDKey(
		[]byte(password),
		salt,
		time,
		memory,
		threads,
		uint32(len(hash)),
	)

	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}

	return false, nil
}

// IsValidUsername returns true if the username contains only ASCII
// letters, digits, '_' or '-'.
func IsValidUsername(username string) bool {
	for _, r := range username {
		if r > 127 {
			return false
		}

		if (r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			r == '_' ||
			r == '-' {
			continue
		}

		return false
	}

	return true
}