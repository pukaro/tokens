package path

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/pukaro/tokens/db"
)

func (header Header) DeleteOne(w http.ResponseWriter, r *http.Request) {
	data := new(db.DataDelete)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(data); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data.Collection = header.Collection
	_, err := header.Session.WithTransaction(context.Background(), data.DeleteOne, header.TxnOpts)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
