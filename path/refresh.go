package path

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/pukaro/tokens/db"

	"github.com/pukaro/tokens/tokens"
)

type RefreshData struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

func (header Header) Refresh(w http.ResponseWriter, r *http.Request) {
	var data RefreshData
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, statusRef, err := tokens.CheckTokens(data.Access, data.Refresh, 0)
	if err != nil || !statusRef.Alive {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	payload, err := tokens.GetPayloadAccess(data.Access)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dataRef := db.DataRefresh{
		Collection: header.Collection,
		Refresh:    data.Refresh,
		PayloadAcc: tokens.PayloadAccess{UserID: payload.UserID},
	}
	result, err := header.Session.WithTransaction(context.Background(), dataRef.RefrashTokens, header.TxnOpts)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	if err = encoder.Encode(result); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
