package config

type Config struct {
	DataDir     string `json:"data_dir"`
	ServerAddr  string `json:"server_addr"`
	ServerToken string `json:"server_token"`
}
