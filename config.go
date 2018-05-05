package main

type globalConfig struct {
	Server   serverConfig
	Database databaseConfig
}

type serverConfig struct {
	Tls      bool   `toml:"tls"`
	Port     int    `toml:"port"`
	Host     string `toml:"host"`
	CertFile string `toml:"cert_file"`
	KeyFile  string `toml:"key_file"`
}

type databaseConfig struct {
	MetaDB string `toml:"meta_db"`
}
