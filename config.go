package main

type globalConfig struct {
	Server serverConfig
}

type serverConfig struct {
	Tls      bool   `toml:"tls"`
	Port     int    `toml:"port"`
	Host     string `toml:"localhost:8080"`
	CertFile string `toml:"cert_file"`
	KeyFile  string `toml:"key_file"`
}
