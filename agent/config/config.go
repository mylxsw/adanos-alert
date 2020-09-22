package config

// Config Agent 配置对象
type Config struct {
	// DataDir Agent 数据存储目录
	DataDir string `json:"data_dir"`
	// ServerAddr Adanos Server GRPC 监听地址
	ServerAddr string `json:"server_addr"`
	// ServerToken Adanos Server GRPC 访问秘钥
	ServerToken string `json:"server_token"`

	// Listen Agent 监听地址
	Listen string `json:"listen"`
	// LogPath Agent 日志目录
	LogPath string `json:"log_path"`
}
