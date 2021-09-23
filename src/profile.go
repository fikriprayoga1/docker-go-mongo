package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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

const baseUrl string = "localhost"

const port string = "8080"

// Object Block
// Database Model
type ModelProfile struct {
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ProfileImage string             `bson:"ProfileImage,omitempty" json:"ProfileImage,omitempty"`
	Name         string             `bson:"Name,omitempty" json:"Name,omitempty"`
	Email        string             `bson:"Email,omitempty" json:"Email,omitempty"`
	Password     string             `bson:"Password,omitempty" json:"Password,omitempty"`
}

// Request Model
type ModelRequest struct {
	Name     string `json:"Name,omitempty"`
	Email    string `json:"Email,omitempty"`
	Password string `json:"Password,omitempty"`
}

type ModelRequest2 struct {
	Id string `json:"Id"`
}

type ModelRequest3 struct {
	Id       string `json:"Id"`
	Name     string `json:"Name,omitempty"`
	Email    string `json:"Email,omitempty"`
	Password string `json:"Password,omitempty"`
}

// Response Model
type ModelResponse struct {
	ResponseMessage string `json:"ResponseMessage"`
	Id              string `json:"Id"`
}

type ModelResponse2 struct {
	ResponseMessage string `json:"ResponseMessage"`
	Id              string `json:"Id"`
	ProfileImage    string `json:"ProfileImage,omitempty"`
	Name            string `json:"Name,omitempty"`
	Email           string `json:"Email,omitempty"`
	Password        string `json:"Password,omitempty"`
}

type ModelResponse3 struct {
	ResponseMessage string `json:"ResponseMessage"`
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
	err2 := client.Ping(ctx, readpref.Primary())
	if err2 != nil {
		panic(err2)
	}

	fmt.Printf("Successfully connected and pinged.\n")

}

func initServer() {
	log.Printf("Server listener started.\n\n")

	http.HandleFunc("/profile/create", createProfile)
	http.HandleFunc("/profile/read", readProfile)
	http.HandleFunc("/profile/update/image", updateProfileImage)
	http.HandleFunc("/profile/update", updateProfile)
	http.HandleFunc("/profile/delete", deleteProfile)
	http.HandleFunc("/image-profile", handleRequest)

	log.Fatal(http.ListenAndServe(baseUrl+":"+port, nil))
}

func createProfile(w http.ResponseWriter, r *http.Request) {
	// Parse body
	var modelRequest ModelRequest
	err := json.NewDecoder(r.Body).Decode(&modelRequest)
	if err != nil {
		panic(err)
	}

	fmt.Printf("name : %v\n", modelRequest.Name)
	fmt.Printf("email : %v\n", modelRequest.Email)
	fmt.Printf("password : %v\n", modelRequest.Password)

	fmt.Printf("Request insert document\n")

	// Create structure
	modelProfile := ModelProfile{
		Name:     modelRequest.Name,
		Email:    modelRequest.Email,
		Password: modelRequest.Password,
	}

	// Insert to database
	coll := client.Database(databaseName).Collection(collectionName)
	res, err2 := coll.InsertOne(context.TODO(), modelProfile)
	if err2 != nil {
		errorHandler(err2, w, "Program Error")
		return
	}
	insertId := res.InsertedID.(primitive.ObjectID).Hex()

	fmt.Printf("Document inserted with id : %v\n\n", insertId)

	// Create response
	modelResponse := ModelResponse{
		ResponseMessage: "Profile data inserted to database",
		Id:              insertId,
	}
	responseJson, err3 := json.Marshal(modelResponse)
	if err3 != nil {
		errorHandler(err3, w, "Program Error")
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(responseJson)
}

func readProfile(w http.ResponseWriter, r *http.Request) {
	// Parse body
	queryParameter := r.URL.Query()
	id := queryParameter["Id"][0]

	fmt.Printf("Request read document by id : %v\n", id)

	// begin find
	myCollection := client.Database(databaseName).Collection(collectionName)
	mId, err2 := primitive.ObjectIDFromHex(id)
	if err2 != nil {
		errorHandler(err2, w, "Program Error")
		return
	}
	filter := bson.D{{Key: "_id", Value: mId}}

	var modelProfile ModelProfile
	err3 := myCollection.FindOne(
		context.TODO(),
		filter,
	).Decode(&modelProfile)
	if err3 != nil {
		if err3 == mongo.ErrNoDocuments {
			errorHandler(err3, w, "Can't find your data")
			return
		}
		errorHandler(err3, w, "Program Error")
		return
	}
	// end find

	fmt.Printf("Document finded\n\n")

	// Create response
	modelResponse2 := ModelResponse2{
		ResponseMessage: "Profile data available in database",
		Id:              modelProfile.Id.Hex(),
		ProfileImage:    modelProfile.ProfileImage,
		Name:            modelProfile.Name,
		Email:           modelProfile.Email,
		Password:        modelProfile.Password,
	}

	responseJson, err4 := json.Marshal(modelResponse2)
	if err4 != nil {
		errorHandler(err4, w, "Program Error")
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(responseJson)

}

func updateProfileImage(w http.ResponseWriter, r *http.Request) {
	// Parse body
	r.ParseMultipartForm(512000)

	id := r.FormValue("Id")
	fmt.Printf("Request update image profile by id : %v\n", id)
	mId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		errorHandler(err, w, "Program Error")
		return
	}

	file, handler, err4 := r.FormFile("ImageProfile")
	if err4 != nil {
		errorHandler(err4, w, "Program Error")
		return
	}

	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	if(handler.Size > 512000) {
		errorHandler(nil, w, "Your file size above 512 Kb")
		return
	}

	// Create structure
	modelProfile := ModelProfile{
		ProfileImage: baseUrl + ":" + port + "/image-profile?Id=" + id,
	}
	update := bson.D{{Key: "$set", Value: modelProfile}}
	myCollection := client.Database(databaseName).Collection(collectionName)

	var modelProfileSample ModelProfile
	filter := bson.D{{Key: "_id", Value: mId}}
	err2 := myCollection.FindOne(
		context.TODO(),
		filter,
	).Decode(&modelProfileSample)
	if err2 != nil {
		if err2 == mongo.ErrNoDocuments {
			errorHandler(err2, w, "Can't find your data")
			return
		}
		errorHandler(err2, w, "Program Error")
		return
	}

	_, err3 := myCollection.UpdateByID(context.TODO(), mId, update)
	if err3 != nil {
		errorHandler(err3, w, "Program Error")
		return
	}

	fmt.Printf("Document updated \n\n")

	// Create file
	dst, err5 := os.Create("image-profile/" + id + ".jpg")
	if err5 != nil {
		errorHandler(err5, w, "Program Error")
		return
	}

	// Copy the uploaded file to the created file on the filesystem
	if _, err6 := io.Copy(dst, file); err6 != nil {
		errorHandler(err6, w, "Program Error")
		return
	}

	// Create response
	modelResponse3 := ModelResponse3{
		ResponseMessage: "Image profile updated",
	}
	responseJson, err7 := json.Marshal(modelResponse3)
	if err7 != nil {
		errorHandler(err7, w, "Program Error")
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(responseJson)
}

func updateProfile(w http.ResponseWriter, r *http.Request) {
	// Parse body
	var modelRequest3 ModelRequest3
	err := json.NewDecoder(r.Body).Decode(&modelRequest3)
	if err != nil {
		panic(err)
	}

	fmt.Printf("modelRequest3 : %v\n", modelRequest3)
	mId, err2 := primitive.ObjectIDFromHex(modelRequest3.Id)
	if err2 != nil {
		errorHandler(err2, w, "Program Error")
		return
	}
	fmt.Printf("Request update document by id : %v\n", modelRequest3.Id)

	podcastsCollection := client.Database(databaseName).Collection(collectionName)

	// Create structure
	modelProfile := ModelProfile{
		Name:     modelRequest3.Name,
		Email:    modelRequest3.Email,
		Password: modelRequest3.Password,
	}
	update := bson.D{{Key: "$set", Value: modelProfile}}

	result, err2 := podcastsCollection.UpdateByID(context.TODO(), mId, update)
	if err2 != nil {
		errorHandler(err2, w, "Program Error")
		return
	}

	if result.ModifiedCount == 0 {
		errorHandler(nil, w, "Can't find your data")
		return
	}

	fmt.Printf("Document updated \n\n")

	// Create response
	modelResponse3 := ModelResponse3{
		ResponseMessage: "Profile data updated",
	}

	responseJson, err3 := json.Marshal(modelResponse3)
	if err3 != nil {
		errorHandler(err3, w, "Program Error")
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(responseJson)
}

func deleteProfile(w http.ResponseWriter, r *http.Request) {
	// Parse body
	var modelRequest2 ModelRequest2
	err := json.NewDecoder(r.Body).Decode(&modelRequest2)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Request delete document by id : %v\n", modelRequest2.Id)

	coll := client.Database(databaseName).Collection(collectionName)
	mId, _ := primitive.ObjectIDFromHex(modelRequest2.Id)

	filter := bson.D{{Key: "_id", Value: mId}}
	result, err2 := coll.DeleteOne(context.TODO(), filter)
	if err2 != nil {
		errorHandler(err2, w, "Program Error")
		return
	}

	if result.DeletedCount == 0 {
		errorHandler(nil, w, "Can't find your data")
		return
	}

	fmt.Printf("Document deleted\n\n")

	// Create response
	modelResponse3 := ModelResponse3{
		ResponseMessage: "Profile data deleted",
	}
	responseJson, err3 := json.Marshal(modelResponse3)
	if err3 != nil {
		errorHandler(err3, w, "Program Error")
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(responseJson)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	queryParameter := r.URL.Query()
	id := queryParameter["Id"][0]

	fileBytes, err := ioutil.ReadFile("image-profile/" + id + ".jpg")
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes)
}

func errorHandler(err error, w http.ResponseWriter, responseMessage string) {
	fmt.Printf("Proccess failed\n\n")
	fmt.Println(err)
	modelResponse3 := ModelResponse3{ResponseMessage: responseMessage}
	responseJson, err2 := json.Marshal(modelResponse3)
	if err2 != nil {
		fmt.Println(err2)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	w.Write(responseJson)

}
