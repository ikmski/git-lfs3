package main

import (
	"log"
	"time"

    "github.com/aws/aws-sdk-go/aws/session"

	"github.com/boltdb/bolt"
	"github.com/BurntSushi/toml"
	"github.com/ikmski/git-lfs3/api"
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

    apiConfig := convertToApiConfig(&config.Server)

    sess, err := session.NewSession()
	if err != nil {
		log.Fatal(err)
	}

	db, err := bolt.Open("meta.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}

    metaDataRepo := adapter.NewMetaDataRepository(db)
    contentRepo, err := adapter.NewContentRepository(sess, "test")
	if err != nil {
		log.Fatal(err)
	}

    batchService := usecase.NewBatchService(metaDataRepo, contentRepo)
    transferService := usecase.NewTransferService(contentRepo)

    batchController := adapter.NewBatchController(batchService)
    transferController := adapter.NewTransferController(transferService)

    batchHandler := api.NewBatchHandler(batchController)
    transferHandler := api.NewTransferHandler(transferController)

	app := api.NewAPI(apiConfig, batchHandler, transferHandler)

	app.Serve()
}

func convertToApiConfig(conf *serverConfig) *api.Config {

    c := &api.Config {
        Tls: conf.Tls,
        Port: conf.Port,
        Host: conf.Host,
        CertFile: conf.CertFile,
        KeyFile: conf.KeyFile,
    }

    return c
}
