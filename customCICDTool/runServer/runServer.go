package runServer

import (
	"fmt"
	"log"
	"net/http"

	webhookhandler "github.com/rk280392/customCICDTool/webHookHandler"
)

func RunServer(port string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", httpHandler)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	fmt.Println("Listening on post 30480")
	log.Fatal(server.ListenAndServe())

}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/webhook":
		webhookhandler.WebhookVerifyRequest(w, r)
	default:
		http.NotFound(w, r)
	}
}
