package main

import (
	"log"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/boltdb/bolt"
	"github.com/ikmski/git-lfs3/adapter"
	"github.com/ikmski/git-lfs3/usecase"
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

	db, err := bolt.Open("meta.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}

	metaDataRepo := adapter.NewMetaDataRepository(db)
	contentRepo, err := adapter.NewContentRepository("test")
	if err != nil {
		log.Fatal(err)
	}

	batchService := usecase.NewBatchService(metaDataRepo, contentRepo)
	transferService := usecase.NewTransferService(metaDataRepo, contentRepo)

	batchController := adapter.NewBatchController(batchService)
	transferController := adapter.NewTransferController(transferService)

	app := newApp(config.Server, batchController, transferController)

	app.Serve()
}
