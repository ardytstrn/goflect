package util

const (
	base62Charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	shortCodeLen  = 8
)

func EncodeBase62(n int64) string {
	encoded := make([]byte, 0, shortCodeLen)
	for range shortCodeLen {
		encoded = append(encoded, base62Charset[n%62])
		n /= 62
	}
	return string(encoded)
}
