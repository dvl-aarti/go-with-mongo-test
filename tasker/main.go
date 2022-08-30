package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func mainTest() {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		// ApplyURI("mongodb://localhost:27017").
		ApplyURI("mongodb+srv://m001-student:m001-mongodb-basics@sandbox.7zffz3a.mongodb.net/?retryWrites=true&w=majority").
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)

	demo := client.Database("demo_with_go")
	demoTestGoCollection := demo.Collection("test_go")

	fmt.Println("Insert method calling :")
	insert(demoTestGoCollection, ctx)
	fmt.Println("Retrieve method calling :")
	retrieve(demoTestGoCollection, ctx)
	fmt.Println("Query method calling :")
	query(demoTestGoCollection, ctx)
	fmt.Println("Update method calling :")
	updateDoc(demoTestGoCollection, ctx)
	fmt.Println("Delete method calling :")
	delete(demoTestGoCollection, ctx)

}

func insert(demoTestGoCollection *mongo.Collection, ctx context.Context) {
	demoTestGoResult, err := demoTestGoCollection.InsertOne(ctx, bson.D{
		{Key: "title", Value: "MongoDB with Go"},
		{Key: "author", Value: "test"},
		{Key: "tags", Value: bson.A{"development", "programming", "coding"}},
	})

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted %v documents into episode collection!\n", (demoTestGoResult.InsertedID))

	result, err2 := demoTestGoCollection.InsertMany(ctx, []interface{}{
		bson.D{
			{"podcast", demoTestGoResult.InsertedID},
			{"title", "GraphQL for API Development"},
			{"description", "Learn about GraphQL from the co-creator of GraphQL, Lee Byron."},
			{"duration", 25},
		},
		bson.D{
			{"podcast", demoTestGoResult.InsertedID},
			{"title", "Progressive Web Application Development"},
			{"description", "Learn about PWA development with Tara Manicsic."},
			{"duration", 32},
		},
	})
	if err2 != nil {
		log.Fatal(err2)
	}
	fmt.Printf("Inserted %v documents into episode collection!\n", len(result.InsertedIDs))
}

func retrieve(demoTestGoCollection *mongo.Collection, ctx context.Context) {
	cursor, err := demoTestGoCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	var episodes []bson.M
	if err = cursor.All(ctx, &episodes); err != nil {
		log.Fatal(err)
	}
	fmt.Println(episodes)
	fmt.Println("Log 2nd position data in json :")
	fmt.Println(episodes[2])
}

func query(demoTestGoCollection *mongo.Collection, ctx context.Context) {
	cursor, err := demoTestGoCollection.Find(ctx, bson.M{"author": "test"})
	if err != nil {
		log.Fatal(err)
	}
	var episodes []bson.M
	if err = cursor.All(ctx, &episodes); err != nil {
		log.Fatal(err)
	}
	log.Println(episodes)
}

func updateDoc(demoTestGoCollection *mongo.Collection, ctx context.Context) {
	//	id, _ := primitive.ObjectIDFromHex("5d9e0173c1305d2a54eb431a")
	//	tags := []string{"development", "programming", "coding"}
	result, err := demoTestGoCollection.UpdateMany(
		ctx,
		bson.M{"author": "test"},
		bson.D{
			{"$set", bson.D{{"author", "test update docuemnts"}}},
			{"$set", bson.D{{"co-author", "test update docuemnts"}}},
			//			{"$set", bson.D{"tags", tags}},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Updated %v Documents!\n", result.ModifiedCount)
}

func delete(demoTestGoCollection *mongo.Collection, ctx context.Context) {
	//	result, err := demoTestGoCollection.DeleteOne(ctx, bson.M{"author": "test update docuemnts"})
	result, err := demoTestGoCollection.DeleteOne(ctx, bson.M{"author": "test"})

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("DeleteOne removed %v document(s)\n", result.DeletedCount)
}
