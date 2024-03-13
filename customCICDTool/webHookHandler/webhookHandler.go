package webhookhandler

import (
	"fmt"

	webhookParser "github.com/rk280392/customCICDTool/interfaces"
)

// Implementation of WebhookHandler
type MyWebhookHandler struct {
	parser webhookParser.WebhookParser
}

func NewMyWebhookHandler(parser webhookParser.WebhookParser) *MyWebhookHandler {
	return &MyWebhookHandler{parser: parser}
}

func (h *MyWebhookHandler) HandleWebhook(payload []byte) error {
	info, err := h.parser.Parse(payload)
	if err != nil {
		return err
	}
	// Process the parsed information
	fmt.Println("Repository:", info.Repository)
	fmt.Println("Branch:", info.Branch)
	fmt.Println("Commit ID:", info.Commit)
	fmt.Println("Repo URL:", info.Url)

	return nil
}
