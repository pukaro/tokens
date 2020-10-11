package path

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Header struct {
	Collection *mongo.Collection
	TxnOpts    *options.TransactionOptions
	Session    mongo.Session
}
