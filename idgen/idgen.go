package idgen

import "crypto/rand"

const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const base = byte(len(alphabet))

// GenerateID generates a random ID of the specified length using all characters from the alphabet.
// Minimum recommended length is 8 for up to 1 million unique IDs.
func GenerateID(length int) string {
	if length <= 0 {
		return ""
	}

	b := make([]byte, length)

	_, _ = rand.Read(b)

	for i := 0; i < length; i++ {
		b[i] = alphabet[b[i]%base]
	}

	return string(b)
}
