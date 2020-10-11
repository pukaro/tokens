package path

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/globalsign/mgo/bson"
)

func (header Header) DeleteAll(w http.ResponseWriter, r *http.Request) {
	data := struct {
		User int64 `json:"user_id"`
	}{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := header.Collection.DeleteMany(ctx, bson.M{"userId": data.User})
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
