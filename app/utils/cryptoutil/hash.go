package cryptoutil

import (
	"crypto/md5"
	"encoding/hex"
)

// GetMD5Hash ...
func GetMD5Hash(key *string) *string {
	hasher := md5.New()
	hasher.Write([]byte(*key))
	hexResp := hex.EncodeToString(hasher.Sum(nil))

	return &hexResp
}
