package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewStorage(uri, dbName string) (mongo.Client) {
	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	
	if err := client.Database(dbName).RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		panic(err)
	}

	log.Println("Batabase connected successfully")
	return *client
}

// func main() {

// 	client, err := mongo.Connect(context.TODO(), options.Client().
// 		ApplyURI(uri))
// 	if err != nil {
// 		panic(err)
// 	}

// 	defer func() {
// 		if err := client.Disconnect(context.TODO()); err != nil {
// 			panic(err)
// 		}
// 	}()

// 	coll := client.Database("sample_mflix").Collection("movies")
// 	title := "Back to the Future"

// 	var result bson.M
// 	err = coll.FindOne(context.TODO(), bson.D{{"title", title}}).
// 		Decode(&result)
// 	if err == mongo.ErrNoDocuments {
// 		fmt.Printf("No document was found with the title %s\n", title)
// 		return
// 	}
// 	if err != nil {
// 		panic(err)
// 	}

// 	jsonData, err := json.MarshalIndent(result, "", "    ")
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("%s\n", jsonData)
// }
