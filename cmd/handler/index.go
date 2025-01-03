package handler

import (
	"context"
	"net/http"

	"github.com/Mikiejoe/go-blog-api/cmd/api"
	"github.com/Mikiejoe/go-blog-api/config"
	"github.com/Mikiejoe/go-blog-api/db"
)

// Handler for Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	client := db.NewStorage(config.Envs.MongoUrl, config.Envs.DBName)
	defer client.Disconnect(context.TODO())

	server := api.NewApiServer(config.Envs.Port, &client)

	// Pass the request to the API server's handler
	server.ServeHTTP(w, r)
	
}
