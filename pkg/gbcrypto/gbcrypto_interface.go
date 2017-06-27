package gbcrypto

// Interface is the interface for cryptography functionality
type Interface interface {
	Encrypt(keyStr string, cryptoText string) string
	Decrypt(keyStr string, cryptoText string) string
}
