package myCICDInterfaces

type WebhookParserInterface interface {
	WebhookRequestParse(payloadData []byte) error
}
