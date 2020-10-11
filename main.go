package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/mux"
	"github.com/pukaro/tokens/path"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

func main() {
	var header path.Header
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(
		"mongodb+srv://medods:yD0D3vEXWCtnFkvZ@cluster0.ufzsk.mongodb.net/<dbname>?retryWrites=true&w=majority",
	))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(context.TODO())

	db := client.Database("medods")
	header.Collection = db.Collection("tokens")
	db.RunCommand(context.TODO(), bson.D{{"create", "tokens"}})

	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	header.TxnOpts = options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	header.Session, err = client.StartSession()
	if err != nil {
		panic(err)
	}
	defer header.Session.EndSession(context.Background())

	r := mux.NewRouter()
	r.HandleFunc("/generate", header.Generate).Methods("POST")
	r.HandleFunc("/refresh", header.Refresh).Methods("POST")
	r.HandleFunc("/delete", header.DeleteOne).Methods("POST")
	r.HandleFunc("/deleteall", header.DeleteAll).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	if err = http.ListenAndServe(":"+port, r); err != nil {
		log.Println(err)
		return
	}
}
