package helpers

type OutputBrandKey struct {
	Key string `json:"key"`
}

func OutputBrandAPIKey(key string) *OutputBrandKey {
	return &OutputBrandKey{
		Key: key,
	}
}
