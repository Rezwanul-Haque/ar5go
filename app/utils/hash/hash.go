package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
)

// GetMD5Hash ...
func GetMD5Hash(key *string) *string {
	hasher := md5.New()
	hasher.Write([]byte(*key))
	hexResp := hex.EncodeToString(hasher.Sum(nil))

	return &hexResp
}

// Get Sha 1 Hash ...
func GetSha1Hash(key string) string {
	hasher := sha1.New()
	_, _ = hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}
