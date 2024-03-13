package interfaces

import (
	webhookparser "github.com/rk280392/customCICDTool/webHookParser"
)

type WebhookHandler interface {
	HandleWebhook(payload []byte) error
	//HandleWebhookAndClone(payload []byte) error
}

type WebhookParser interface {
	Parse(payload []byte) (webhookparser.RepositoryInfo, error)
}
