package config

type Config struct {
	Addr string `json:"addr"`

	Database DatabaseConfig `json:"database"`
}

type DatabaseConfig struct {
	DBURL  string `json:"dburl"`
	DBName string `json:"dbname"`
}
