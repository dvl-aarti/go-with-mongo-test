package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ctx          = context.Background()
	kmsProviders map[string]map[string]interface{}
	schemaMap    bson.M
)

func createDataKey() primitive.Binary {
	kvClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://m001-student:m001-mongodb-basics@sandbox.7zffz3a.mongodb.net/?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	kvClient.Database("keyvault").Collection("datakeys").Drop(ctx)
	clientEncryptionOpts := options.ClientEncryption().SetKeyVaultNamespace("keyvault.datakeys").SetKmsProviders(kmsProviders)
	clientEncryption, err := mongo.NewClientEncryption(kvClient, clientEncryptionOpts)
	if err != nil {
		log.Fatal(err)
	}
	defer clientEncryption.Close(ctx)
	dataKeyId, err := clientEncryption.CreateDataKey(ctx, "local", options.DataKey())
	if err != nil {
		log.Fatal(err)
	}
	return dataKeyId
}

func createEncryptedClient(dataKeyIdBase64 string) *mongo.Client {
	schemaMap = readSchemaFromFile("schema.json", dataKeyIdBase64)
	mongocryptdOpts := map[string]interface{}{
		"mongodcryptdBypassSpawn": true,
	}
	autoEncryptionOpts := options.AutoEncryption().
		SetKeyVaultNamespace("keyvault.datakeys").
		SetKmsProviders(kmsProviders).
		SetSchemaMap(schemaMap).
		SetExtraOptions(mongocryptdOpts)
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://m001-student:m001-mongodb-basics@sandbox.7zffz3a.mongodb.net/?retryWrites=true&w=majority").SetAutoEncryptionOptions(autoEncryptionOpts))
	if err != nil {
		log.Fatal(err)
	}
	return mongoClient
}

func readSchemaFromFile(file string, dataKeyIdBase64 string) bson.M {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	content = []byte(fmt.Sprintf(string(content), dataKeyIdBase64))
	var doc bson.M
	if err = bson.UnmarshalExtJSON(content, false, &doc); err != nil {
		log.Fatal(err)
	}
	return doc
}

func main() {
	fmt.Println("Starting the application...")
	localKey := make([]byte, 96)
	if _, err := rand.Read(localKey); err != nil {
		log.Fatal(err)
	}
	kmsProviders = map[string]map[string]interface{}{
		"local": {
			"key": localKey,
		},
	}
	dataKeyId := createDataKey()
	client := createEncryptedClient(base64.StdEncoding.EncodeToString(dataKeyId.Data))
	defer client.Disconnect(ctx)
	collection := client.Database("fle-example").Collection("people")
	collection.Drop(context.TODO())
	if _, err := collection.InsertOne(context.TODO(), bson.M{"name": "Nic Raboy", "ssn": "123456"}); err != nil {
		log.Fatal(err)
	}
	result, err := collection.FindOne(context.TODO(), bson.M{"ssn": "123456"}).DecodeBytes()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
}
