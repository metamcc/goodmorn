package web

import (
	"fmt"
	"net/http"

	"github.com/mcc/hyperledger-fabric-server/web/controllers"
)

func Serve(app *controllers.Application) {
	http.HandleFunc("/query", app.QueryHandler)
	http.HandleFunc("/invoke", app.InvokeHandler)

	fmt.Println("Listening (http://localhost:3000/) ...")
	http.ListenAndServe(":3000", nil)
}
