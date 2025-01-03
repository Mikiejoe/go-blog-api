package api

import (
	"log"
	"net/http"

	"github.com/Mikiejoe/go-blog-api/config"
	"github.com/Mikiejoe/go-blog-api/services/blog"
	"github.com/Mikiejoe/go-blog-api/services/user"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type ApiServer struct {
	addr   string
	client *mongo.Client
	router *mux.Router
}

func NewApiServer(addr string, client *mongo.Client) *ApiServer {
	router := mux.NewRouter()
	return &ApiServer{addr: addr, client: client,router:router}
}

func (s *ApiServer) Run() error {
	router := mux.NewRouter()

	subRouter := router.PathPrefix("/api/v1").Subrouter()
	database := s.client.Database(config.Envs.DBName)
	userCollection:= database.Collection(config.Collections.Users)
	blogsCollection:= database.Collection(config.Collections.Blogs)

	blogStore:= blog.NewStore(blogsCollection)
	userStore := user.NewStore(userCollection)

	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subRouter)
	blogHander:= blog.NewHandler(blogStore,userStore)
	blogHander.RegisterRoutes(subRouter)
	log.Println("Listening on port: ", s.addr)
	return http.ListenAndServe(":3000", subRouter)
}

func (s *ApiServer) setupRoutes() {
	
	subRouter := s.router.PathPrefix("/api/v1").Subrouter()
	database := s.client.Database(config.Envs.DBName)
	userCollection:= database.Collection(config.Collections.Users)
	blogsCollection:= database.Collection(config.Collections.Blogs)

	blogStore:= blog.NewStore(blogsCollection)
	userStore := user.NewStore(userCollection)

	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subRouter)
	blogHander:= blog.NewHandler(blogStore,userStore)
	blogHander.RegisterRoutes(subRouter)
}

func (s *ApiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.setupRoutes()
	s.router.ServeHTTP(w, r)
}