package helpers

func OuputToken(token string, expirationSec int64) M {
	return M{
		"token":      token,
		"expiration": expirationSec,
	}
}
