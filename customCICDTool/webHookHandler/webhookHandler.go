package webhookhandler

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/rk280392/customCICDTool/myCICDInterfaces"
)

func WebhookVerifyRequest(w http.ResponseWriter, r *http.Request, parser myCICDInterfaces.WebhookParserInterface) {
	if r.Method != "POST" {
		http.Error(w, "Only POST request in valid", http.StatusMethodNotAllowed)
		return
	}

	eventType := r.Header.Get("X-GitHub-Event")
	if eventType != "push" {
		fmt.Println("Only push events are supported")
		return
	}

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Cannot read the webhook payload")
	}
	defer r.Body.Close()

	err = parser.WebhookRequestParse(payload)
	if err != nil {
		http.Error(w, "Failed to parse webhook payload", http.StatusInternalServerError)
		log.Println("Failed to parse webhook payload: ", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
