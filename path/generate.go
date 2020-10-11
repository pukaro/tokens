package path

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/pukaro/tokens/db"
	"github.com/pukaro/tokens/tokens"
)

func (header Header) Generate(w http.ResponseWriter, r *http.Request) {
	var data tokens.PayloadAccess
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	access, refresh, err := db.Generate(context.Background(), header.Collection, data)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tokensData := map[string]string{
		"access":  access,
		"refresh": refresh,
	}
	encoder := json.NewEncoder(w)
	if err = encoder.Encode(tokensData); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
