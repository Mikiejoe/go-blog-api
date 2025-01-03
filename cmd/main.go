package main

import (
	"context"
	"log"

	"github.com/Mikiejoe/go-blog-api/cmd/api"
	"github.com/Mikiejoe/go-blog-api/config"
	"github.com/Mikiejoe/go-blog-api/db"
)


func main() {

	client := db.NewStorage(config.Envs.MongoUrl,config.Envs.DBName)
	defer client.Disconnect(context.TODO())
	server := api.NewApiServer(config.Envs.Port, &client)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}


}
