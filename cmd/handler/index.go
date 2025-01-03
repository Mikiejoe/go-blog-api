package handler

import (
	"context"
	"log"
	"net/http"
	"fmt"

	"github.com/Mikiejoe/go-blog-api/cmd/api"
	"github.com/Mikiejoe/go-blog-api/config"
	"github.com/Mikiejoe/go-blog-api/db"
)

var server *api.ApiServer

func Handler(w http.ResponseWriter, r *http.Request) {
	// Set up the database connection
	client := db.NewStorage(config.Envs.MongoUrl, config.Envs.DBName)
	defer client.Disconnect(context.TODO())

	server = api.NewApiServer(config.Envs.Port, &client)
	if err := server.Run(); err != nil {
		http.Error(w, fmt.Sprintf("Error starting server: %s", err), http.StatusInternalServerError)
		log.Fatal(err)
	}
}


