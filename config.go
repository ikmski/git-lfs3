package main

type globalConfig struct {
	Server   serverConfig
	Database databaseConfig
	S3       s3Config
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

type s3Config struct {
	AwsAccessKeyID     string `toml:"aws_access_key_id"`
	AwsSecretAccessKey string `toml:"aws_secret_access_key"`
	Region             string `toml:"region"`
	Bucket             string `toml:"bucket"`
}
