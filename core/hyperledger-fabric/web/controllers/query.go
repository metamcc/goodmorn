package controllers

import (
	"encoding/json"
	"log"
	"net/http"
)

func (app *Application) QueryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Query().Get("key")
	log.Printf("%+v\n", key)
	helloValue, err := app.Fabric.QueryHello(key)
	if err != nil {
		http.Error(w, "Unable to query the blockchain", 500)
	}

	data := &FabricData{Key: key, Value: helloValue}
	log.Printf("%+v\n", data)
	b, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(b)
}
