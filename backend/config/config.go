package config

type Config struct {
	Addr string `json:"addr"`

	Salt string `json:"salt"`

	Database   DatabaseConfig   `json:"database"`
	Attachment AttachmentConfig `json:"attachment"`
}

type DatabaseConfig struct {
	DBURL  string `json:"dburl"`
	DBName string `json:"dbname"`
}

type AttachmentConfig struct {
	Internal InternalAttachmentConfig `json:"internal"`
}

type InternalAttachmentConfig struct {
	Dir       string `json:"dir"`
	SizeLimit int64  `json:"size_limit"`
}
