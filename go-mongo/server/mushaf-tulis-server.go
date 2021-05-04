package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/* Variable Block */
var db *mongo.Database

const databaseName string = "mushafTulis"

var collectionName = [6]string{"seriesOne", "seriesTwo", "seriesThree", "seriesFour", "seriesFive", "seriesSix"}

/* Object Block */
type DatabaseObject struct {
	SerialNumber string
	UUIDList     []string
}

type AccessRequestObject struct {
	Series       uint8
	SerialNumber string
	UUID         string
}

/* Function Block */
func initServerListener() {
	log.Println("Server listener started")

	mux := http.NewServeMux()
	mux.HandleFunc("/requestAccess", requestAccessHandler)
	http.ListenAndServe(":8080", mux)
}

func initDatabase() {
	log.Println("Init database started")
	client, err0 := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err0 != nil {
		log.Fatal(err0)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err1 := client.Connect(ctx)
	if err1 != nil {
		log.Fatal(err1)
	}
	db = client.Database(databaseName)

}

func initSerialNumber() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err0 := db.ListCollectionNames(ctx, bson.M{})
	if err0 != nil {
		log.Fatal(err0)
	}

	if len(result) < 1 {
		insertSerialNumber()
	}

}

func insertSerialNumber() {
	sn1, _ := ioutil.ReadFile("/usr/local/go/src/server/serial-number-series-1.txt")
	sn2, _ := ioutil.ReadFile("/usr/local/go/src/server/serial-number-series-2.txt")
	sn3, _ := ioutil.ReadFile("/usr/local/go/src/server/serial-number-series-3.txt")
	sn4, _ := ioutil.ReadFile("/usr/local/go/src/server/serial-number-series-4.txt")
	sn5, _ := ioutil.ReadFile("/usr/local/go/src/server/serial-number-series-5.txt")
	sn6, _ := ioutil.ReadFile("/usr/local/go/src/server/serial-number-series-6.txt")

	snList1 := strings.Split(string(sn1), "\n")
	snList2 := strings.Split(string(sn2), "\n")
	snList3 := strings.Split(string(sn3), "\n")
	snList4 := strings.Split(string(sn4), "\n")
	snList5 := strings.Split(string(sn5), "\n")
	snList6 := strings.Split(string(sn6), "\n")

	collection1 := db.Collection(collectionName[0])
	collection2 := db.Collection(collectionName[1])
	collection3 := db.Collection(collectionName[2])
	collection4 := db.Collection(collectionName[3])
	collection5 := db.Collection(collectionName[4])
	collection6 := db.Collection(collectionName[5])
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	insertDatabase(collection1, ctx, snList1)
	insertDatabase(collection2, ctx, snList2)
	insertDatabase(collection3, ctx, snList3)
	insertDatabase(collection4, ctx, snList4)
	insertDatabase(collection5, ctx, snList5)
	insertDatabase(collection6, ctx, snList6)

}

func insertDatabase(collection *mongo.Collection, ctx context.Context, dataList []string) {
	var snData []interface{}

	log.Println(len(dataList))
	// _ is index, where a is value
	for _, a := range dataList {
		snData = append(snData, bson.D{
			{Key: "SerialNumber", Value: a},
			{Key: "UUID", Value: bson.A{}},
		})
	}

	res, err0 := collection.InsertMany(ctx, snData)
	if err0 != nil {
		log.Fatal(err0)
	}
	log.Println(res.InsertedIDs...)
}

func requestAccessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var databaseObject DatabaseObject
	var requestObject AccessRequestObject
	err0 := json.NewDecoder(r.Body).Decode(&requestObject)
	if err0 != nil {
		log.Fatal(err0)
	}

	series := requestObject.Series
	serialNumber := requestObject.SerialNumber
	uuid := requestObject.UUID

	var mCollectionName string

	switch series {
	case 1:
		mCollectionName = collectionName[0]
	case 2:
		mCollectionName = collectionName[1]
	case 3:
		mCollectionName = collectionName[2]
	case 4:
		mCollectionName = collectionName[3]
	case 5:
		mCollectionName = collectionName[4]
	case 6:
		mCollectionName = collectionName[5]
	default:
		mCollectionName = collectionName[0]

	}

	collection := db.Collection(mCollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	log.Println(series)
	log.Println(serialNumber)
	log.Println(uuid)
	err1 := collection.FindOne(ctx, bson.M{"SerialNumber": serialNumber, "uuid": bson.M{"$all": bson.A{uuid}}}).Decode(&databaseObject)
	log.Println(databaseObject)
	if err1 != nil {
		log.Println("Ada error")
		err2 := collection.FindOne(ctx, bson.M{"SerialNumber": serialNumber, "uuid": bson.M{"$not": bson.M{"$size": 4}}}).Decode(&databaseObject)

		if err2 != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{ "Response": "Access Failed" }`))
			log.Println(err2)
			return
		}
		uuidList := databaseObject.UUIDList
		uuidList = append(uuidList, uuid)

		filter := bson.M{"SerialNumber": serialNumber}
		update := bson.M{"$set": bson.M{"uuid": uuidList}}

		_, err3 := collection.UpdateOne(ctx, filter, update)
		if err3 != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{ "Response": "Access Failed" }`))
			log.Println(err3)
			return
		}

	}
	w.Write([]byte(`{ "Response": "Access Granted" }`))

}

func main() {
	initDatabase()
	initSerialNumber()
	initServerListener()

}
