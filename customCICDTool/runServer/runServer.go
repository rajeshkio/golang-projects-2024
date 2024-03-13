package runServer

import (
	"fmt"
	"log"
	"net/http"

	webhookhandler "github.com/rk280392/customCICDTool/webHookHandler"
	webhookParser "github.com/rk280392/customCICDTool/webHookParser"
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
		parser := &webhookParser.RequestParse{}
		webhookhandler.WebhookVerifyRequest(w, r, parser)
	default:
		http.NotFound(w, r)
	}
}
