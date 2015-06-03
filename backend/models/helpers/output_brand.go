package helpers

type OutputBrandKEyInfo struct {
	Key string `json:"key"`
}

func OutputBrandAPIKey(key string) *OutputBrandKEyInfo {
	return &OutputBrandKEyInfo{
		Key: key,
	}
}
