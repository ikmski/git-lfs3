package main

import (
	"log"

	"github.com/BurntSushi/toml"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
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

	credentials := credentials.NewStaticCredentials(
		config.S3.AwsAccessKeyID,
		config.S3.AwsSecretAccessKey,
		"")

	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials,
		Region:      aws.String(config.S3.Region),
	}))

	contentStore, err := NewContentStore(sess, config.S3.Bucket)
	if err != nil {
		log.Fatal(err)
	}

	app := newApp(metaStore, contentStore)

	app.Serve()
}
