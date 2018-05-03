package main

import (
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
)

var config globalConfig

const (
	configFileName = "config.toml"
)

func main() {

	_, err := toml.DecodeFile(configFileName, &config)
	if err != nil {
		log.Fatal(err)
	}

	e := newEngine()

	if config.Server.Tls {
		e.RunTLS(
			fmt.Sprintf(":%d", config.Server.Port),
			config.Server.CertFile,
			config.Server.KeyFile)
	} else {
		e.Run(fmt.Sprintf(":%d", config.Server.Port))
	}
}
