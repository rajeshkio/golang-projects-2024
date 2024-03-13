package runServer

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	webhookHandler "github.com/rk280392/customCICDTool/interfaces"
)

type Server struct {
	handler webhookHandler.WebhookHandler
}

func NewServer(handler webhookHandler.WebhookHandler) *Server {
	return &Server{handler: handler}
}

func (s *Server) Start() {
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		payload, err := s.webHookHandler(w, r)
		if err != nil {
			fmt.Println("Error handling webhook:", err)
			http.Error(w, "Error handling webhook", http.StatusInternalServerError)
			return
		}

		if err := s.handler.HandleWebhook(payload); err != nil {
			fmt.Println("Error handling the webhook", err)
			http.Error(w, "Error handling the webhook", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	fmt.Println("Running HTTP server on port 30480")
	http.ListenAndServe(":30480", nil)
}

func (s *Server) webHookHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return nil, errors.New("method not allowed")
	}

	defer r.Body.Close()
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return nil, fmt.Errorf("failed to read request body: %v", err)
	}

	eventType := r.Header.Get("X-GitHub-Event")
	if eventType != "push" {
		fmt.Println("Ignoring non-push event:", eventType)
		return payload, fmt.Errorf("ignoring non-push event: %s", eventType)
	}
	return payload, nil
	//fmt.Println(string(payload))

}
