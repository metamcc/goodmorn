package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func (app *Application) InvokeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Unmarshal
	var rData *FabricData
	err = json.Unmarshal(b, &rData)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	log.Println(rData.Key, rData.Value)

	data := &struct {
		TransactionId string
		Success       bool
		Response      bool
	}{
		TransactionId: "",
		Success:       false,
		Response:      false,
	}
	txid, err := app.Fabric.InvokeHello(rData.Key, rData.Value)
	if err != nil {
		log.Panicln(err)
		http.Error(w, "Unable to invoke hello in the blockchain", 500)
	}
	data.TransactionId = txid
	data.Success = true
	data.Response = true
	w.Write(b)
}
