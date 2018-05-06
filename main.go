package main

import (
	"log"

	"github.com/BurntSushi/toml"
)

var config globalConfig

const (
	configFileName   = "config.toml"
	contentMediaType = "application/vnd.git-lfs"
	metaMediaType    = contentMediaType + "+json"
)

func main() {

	_, err := toml.DecodeFile(configFileName, &config)
	if err != nil {
		log.Fatal(err)
	}

	metaStore, err := NewMetaStore(config.Database.MetaDB)
	if err != nil {
		log.Fatal(err)
	}

	app := newApp(metaStore)

	app.Serve()
}
