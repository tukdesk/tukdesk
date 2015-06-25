package helpers

type OutputToken struct {
	Token      string `json:"token"`
	Expiration int64  `json:"expiration"`
}

func OutputTokenInfo(token string, expirationSec int64) *OutputToken {
	return &OutputToken{
		Token:      token,
		Expiration: expirationSec,
	}
}
