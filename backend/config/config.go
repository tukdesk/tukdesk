package config

type Config struct {
	Addr string `json:"addr"`

	Salt string `json:"salt"`

	Database DatabaseConfig `json:"database"`
	File     FileConfig     `json:"file"`
}

type DatabaseConfig struct {
	DBURL  string `json:"dburl"`
	DBName string `json:"dbname"`
}

type FileConfig struct {
	Internal InternalFileConfig `json:"internal"`
}

type InternalFileConfig struct {
	Dir       string `json:"dir"`
	SizeLimit int64  `json:"size_limit"`
}
