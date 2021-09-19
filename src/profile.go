package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Variable Block
var client *mongo.Client

const databaseName string = "example"

const collectionName string = "profile"

// Object Block
type ProfileDetail struct {
	ProfileImage string `bson:"ProfileImage,omitempty" json:"ProfileImage,omitempty"`
	Name         string `bson:"Name,omitempty" json:"Name,omitempty"`
	Email        string `bson:"Email,omitempty" json:"Email,omitempty"`
	Password     string `bson:"Password,omitempty" json:"Password,omitempty"`
}

type CreateResponse struct {
	ResponseMessage string `json:"ResponseMessage"`
	Id              string `json:"_id"`
}

type ProfileDetailComplete struct {
	Id           primitive.ObjectID `bson:"_id" json:"_id"`
	ProfileImage string             `bson:"ProfileImage,omitempty" json:"ProfileImage,omitempty"`
	Name         string             `bson:"Name,omitempty" json:"Name,omitempty"`
	Email        string             `bson:"Email,omitempty" json:"Email,omitempty"`
	Password     string             `bson:"Password,omitempty" json:"Password,omitempty"`
}

type ReadResponse struct {
	ResponseMessage string                 `json:"ResponseMessage,omitempty"`
	ProfileDetail   *ProfileDetailComplete `json:"ProfileDetail,omitempty"`
}

type FailedResponse struct {
	ResponseMessage string `default:"Program Error"`
}

func main() {
	initDatabase()
	initServer()

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
}

func initDatabase() {
	// Replace the uri string with your MongoDB deployment's connection string.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	databaseUrl := "mongodb://" + os.Getenv("MONGO_HOST") + ":27017"
	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(databaseUrl))
	if err != nil {
		panic(err)
	}

	// Ping the primary
	err1 := client.Ping(ctx, readpref.Primary())
	if err1 != nil {
		panic(err1)
	}

	fmt.Printf("Successfully connected and pinged.\n")

}

func initServer() {
	log.Printf("Server listener started.\n\n")

	http.HandleFunc("/profile/create", createProfile)
	http.HandleFunc("/profile/read", readProfile)
	http.HandleFunc("/profile/update", updateProfile)
	http.HandleFunc("/profile/delete", deleteProfile)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createProfile(w http.ResponseWriter, r *http.Request) {
	// Parse body
	err := r.ParseMultipartForm(512000)
	if err != nil {
		errorHandler(err, w, FailedResponse{ResponseMessage: "Program Error"})
		return
	}

	profileImage := r.FormValue("ProfileImage")
	name := r.FormValue("Name")
	email := r.FormValue("Email")
	password := r.FormValue("Password")

	fmt.Printf("profileImage : %v\n", profileImage)
	fmt.Printf("name : %v\n", name)
	fmt.Printf("email : %v\n", email)
	fmt.Printf("password : %v\n", password)

	fmt.Printf("Request insert document\n")

	// Create structure
	profileData := ProfileDetail{
		ProfileImage: profileImage,
		Name:         name,
		Email:        email,
		Password:     password,
	}

	// Insert to database
	coll := client.Database(databaseName).Collection(collectionName)
	res, err1 := coll.InsertOne(context.TODO(), profileData)
	if err1 != nil {
		errorHandler(err1, w, FailedResponse{ResponseMessage: "Program Error"})
		return
	}
	insertId := res.InsertedID.(primitive.ObjectID).Hex()

	fmt.Printf("Document inserted with id : %v\n\n", insertId)

	// Create response
	responseData := CreateResponse{
		ResponseMessage: "Profile data inserted to database",
		Id:              insertId,
	}
	responseJson, err2 := json.Marshal(responseData)
	if err2 != nil {
		errorHandler(err2, w, FailedResponse{ResponseMessage: "Program Error"})
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(responseJson)
}

func readProfile(w http.ResponseWriter, r *http.Request) {
	// Parse body
	err := r.ParseMultipartForm(512000)
	if err != nil {
		errorHandler(err, w, FailedResponse{ResponseMessage: "Program Error"})
		return
	}

	id := r.FormValue("Id")

	fmt.Printf("Request read document by id : %v\n", id)

	// begin find
	myCollection := client.Database(databaseName).Collection(collectionName)
	mId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{Key: "_id", Value: mId}}

	// Find the document for which the _id field matches id.
	// Specify the Sort option to sort the documents by age.
	// The first document in the sorted order will be returned.
	var profileDetailComplete ProfileDetailComplete
	err1 := myCollection.FindOne(
		context.TODO(),
		filter,
	).Decode(&profileDetailComplete)
	if err1 != nil {
		if err1 == mongo.ErrNoDocuments {
			errorHandler(err1, w, FailedResponse{ResponseMessage: "Can't find your data"})
			return
		}
		errorHandler(err1, w, FailedResponse{ResponseMessage: "Program Error"})
		return
	}
	// end find

	fmt.Printf("Document finded\n\n")

	// Create response
	responseData := ReadResponse{
		ResponseMessage: "Profile data inserted to database",
		ProfileDetail:   &profileDetailComplete,
	}
	responseJson, err2 := json.Marshal(responseData)
	if err2 != nil {
		errorHandler(err2, w, FailedResponse{ResponseMessage: "Program Error"})
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(responseJson)

}

func updateProfile(w http.ResponseWriter, r *http.Request) {
	// Parse body
	err := r.ParseMultipartForm(512000)
	if err != nil {
		errorHandler(err, w, FailedResponse{ResponseMessage: "Program Error"})
		return
	}

	id := r.FormValue("Id")
	profileImage := r.FormValue("ProfileImage")
	name := r.FormValue("Name")
	email := r.FormValue("Email")
	password := r.FormValue("Password")

	fmt.Printf("profileImage : %v\n", profileImage)
	fmt.Printf("name : %v\n", name)
	fmt.Printf("email : %v\n", email)
	fmt.Printf("password : %v\n", password)

	fmt.Printf("Request update document by id : %v\n", id)

	podcastsCollection := client.Database(databaseName).Collection(collectionName)
	mId, _ := primitive.ObjectIDFromHex(id)

	// Create structure
	profileData := ProfileDetail{
		ProfileImage: profileImage,
		Name:         name,
		Email:        email,
		Password:     password,
	}

	filter := bson.D{{Key: "_id", Value: mId}}

	result, err1 := podcastsCollection.ReplaceOne(context.TODO(), filter, profileData)

	if err1 != nil {
		errorHandler(err1, w, FailedResponse{ResponseMessage: "Program Error"})
		return
	}

	if result.ModifiedCount == 0 {
		errorHandler(nil, w, FailedResponse{ResponseMessage: "Can't find your data"})
		return
	}

	fmt.Printf("Document updated \n\n")

	// Create response
	profileDetailComplete := ProfileDetailComplete{
		Id:           mId,
		ProfileImage: profileImage,
		Name:         name,
		Email:        email,
		Password:     password,
	}
	responseData := ReadResponse{
		ResponseMessage: "Profile data updated",
		ProfileDetail:   &profileDetailComplete,
	}
	responseJson, err2 := json.Marshal(responseData)
	if err2 != nil {
		errorHandler(err2, w, FailedResponse{ResponseMessage: "Program Error"})
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(responseJson)
}

func deleteProfile(w http.ResponseWriter, r *http.Request) {
	// Parse body
	err := r.ParseMultipartForm(512000)
	if err != nil {
		errorHandler(err, w, FailedResponse{ResponseMessage: "Program Error"})
		return
	}

	id := r.FormValue("Id")

	fmt.Printf("Request delete document by id : %v\n", id)

	coll := client.Database(databaseName).Collection(collectionName)
	mId, _ := primitive.ObjectIDFromHex(id)

	filter := bson.D{{Key: "_id", Value: mId}}
	result, err1 := coll.DeleteOne(context.TODO(), filter)
	if err1 != nil {
		errorHandler(err1, w, FailedResponse{ResponseMessage: "Program Error"})
		return
	}

	if result.DeletedCount == 0 {
		errorHandler(nil, w, FailedResponse{ResponseMessage: "Can't find your data"})
		return
	}

	fmt.Printf("Document deleted\n\n")

	// Create response
	responseData := CreateResponse{
		ResponseMessage: "Profile data deleted",
		Id:              id,
	}
	responseJson, err2 := json.Marshal(responseData)
	if err2 != nil {
		errorHandler(err2, w, FailedResponse{ResponseMessage: "Program Error"})
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(responseJson)
}

func errorHandler(err error, w http.ResponseWriter, message FailedResponse) {
	fmt.Printf("Proccess failed\n\n")
	responseJson, err1 := json.Marshal(message)
	if err1 != nil {
		fmt.Println(err1)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	w.Write(responseJson)
	fmt.Println(err)

}
