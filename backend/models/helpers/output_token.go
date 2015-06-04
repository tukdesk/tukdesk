package helpers

type OutputUserToken struct {
	Token      string `json:"token"`
	Expiration int64  `json:"expiration"`
}

func OutputUserTokenInfo(token string, expirationSec int64) *OutputUserToken {
	return &OutputUserToken{
		Token:      token,
		Expiration: expirationSec,
	}
}
