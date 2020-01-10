package main

import (
	"log"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/boltdb/bolt"
	"github.com/ikmski/git-lfs3/adapter"
	"github.com/ikmski/git-lfs3/usecase"
)

const (
	configFileName = "config.toml"
)

func main() {

	var config globalConfig
	_, err := toml.DecodeFile(configFileName, &config)
	if err != nil {
		log.Fatal(err)
	}

	app, err := initializeApp(config)
	if err != nil {
		log.Fatal(err)
	}

	app.serve()
}

func initializeApp(config globalConfig) (*app, error) {

	db, err := bolt.Open("meta.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	metaDataRepo := adapter.NewMetaDataRepository(db)
	contentRepo, err := adapter.NewContentRepository("test")
	if err != nil {
		return nil, err
	}

	batchService := usecase.NewBatchService(metaDataRepo, contentRepo)
	transferService := usecase.NewTransferService(metaDataRepo, contentRepo)

	batchController := adapter.NewBatchController(batchService)
	transferController := adapter.NewTransferController(transferService)

	app := newApp(config.Server, batchController, transferController)

	return app, nil
}
