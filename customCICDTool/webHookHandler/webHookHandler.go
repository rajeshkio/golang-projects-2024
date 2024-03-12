package webHookHandler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func RunServer() {
	http.HandleFunc("/webhook", webHookHandler)
	fmt.Println("Starting server on port 30480")
	http.ListenAndServe(":30480", nil)
}

func webHookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	//fmt.Println(string(payload))
	w.WriteHeader(http.StatusOK)

	eventType := r.Header.Get("X-GitHub-Event")
	if eventType != "push" {
		fmt.Println("Ignoring non-push event:", eventType)
		return
	}

	WebhookParser(payload)
}

type WebhookPayload struct {
	Ref        string `json:"ref"`
	Repository struct {
		Name string `json:"name"`
	} `json:"repository"`
	Commit []struct {
		Id string `json:"id"`
	} `json:"commits"`
}

func WebhookParser(payloadData []byte) {
	var payload WebhookPayload

	if err := json.Unmarshal(payloadData, &payload); err != nil {
		fmt.Println("Error parsing JSON payload:", err)
		return
	}

	repoName := payload.Repository.Name
	branchName := payload.Ref
	commitID := payload.Commit

	fmt.Println("Repository:", repoName)
	fmt.Println("Branch:", branchName)
	fmt.Println("Commit ID:", commitID)
}

/*
TODO right now the code is tight coupled making any changes in the halder or parser need to rerun the server code.
Need to understand and decouple it.

*/
