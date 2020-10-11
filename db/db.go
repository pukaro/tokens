package db

import (
	"context"
	"encoding/base64"

	"github.com/globalsign/mgo/bson"
	"github.com/pukaro/tokens/tokens"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ed25519"
)

type DataTokens struct {
	ID      primitive.ObjectID `bson:"_id"`
	UserID  int64              `bson:"userId"`
	Access  string             `bson:"access"`
	Refresh []byte             `bson:"refresh"`
}

type DataRefresh struct {
	Collection *mongo.Collection
	Refresh    string
	PayloadAcc tokens.PayloadAccess
}

type DataDelete struct {
	Collection *mongo.Collection
	Refresh    string `json:"refresh"`
}

// Генерирует Access, Refresh tokens, а так же записывает их в БД.
func Generate(ctx context.Context, collection *mongo.Collection, payload tokens.PayloadAccess) (access, refresh string, err error) {
	access, err = payload.GenerateAccess()
	if err != nil {
		return
	}
	refresh, err = tokens.GenerateRefresh(access)
	if err != nil {
		return
	}

	dataTokens := new(DataTokens)
	dataTokens.ID = primitive.NewObjectID()
	dataTokens.UserID = payload.UserID
	signAcc := ed25519.Sign(tokens.PrivateKey, []byte(access))
	dataTokens.Access = base64.StdEncoding.EncodeToString(signAcc)
	dataTokens.Refresh, err = bcrypt.GenerateFromPassword([]byte(refresh), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	_, err = collection.InsertOne(ctx, dataTokens)
	return
}

// Удаляет из базы данных старые токени и генерирует новую пару токенов на основе PayloadAcc
// Новые токины записываются в БД, а также возвращаются ввиде map[string]string{"access": "accessToken", "refresh": "refreshToken"}
func (data DataRefresh) RefrashTokens(sessCtx mongo.SessionContext) (interface{}, error) {
	payloadRef, err := tokens.GetPayliadRefresh(data.Refresh)
	if err != nil {
		return nil, err
	}

	dataToken := new(DataTokens)
	result := data.Collection.FindOneAndDelete(sessCtx, bson.M{"userId": data.PayloadAcc.UserID, "access": payloadRef.Access})
	if result.Err() != nil {
		return nil, result.Err()
	}
	if err = result.Decode(dataToken); err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(dataToken.Refresh, []byte(data.Refresh))
	if err != nil {
		return nil, err
	}

	newTokens := map[string]string{}
	newTokens["access"], newTokens["refresh"], err = Generate(sessCtx, data.Collection, data.PayloadAcc)
	if err != nil {
		return nil, err
	}
	return newTokens, nil
}

// Удаляет из базы данных Refresh token.
func (data DataDelete) DeleteOne(ctx mongo.SessionContext) (interface{}, error) {
	payload, err := tokens.GetPayliadRefresh(data.Refresh)
	if err != nil {
		return nil, err
	}

	dataTokens := new(DataTokens)
	result := data.Collection.FindOneAndDelete(ctx, bson.M{"access": payload.Access})
	if result.Err() != nil {
		return nil, result.Err()
	}
	if err = result.Decode(dataTokens); err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword(dataTokens.Refresh, []byte(data.Refresh))
	return nil, err
}
