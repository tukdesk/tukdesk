package helpers

type OutputTokenInfo struct {
	Token      string `json:"token"`
	Expiration int64  `json:"expiration"`
}

func OutputUserToken(token string, expirationSec int64) *OutputTokenInfo {
	return &OutputTokenInfo{
		Token:      token,
		Expiration: expirationSec,
	}
}
