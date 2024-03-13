package webhookParser

import (
	"encoding/json"
	"fmt"
)

type RequestParse struct {
	Ref    string `json:"ref"`
	Repo   repo   `json:"repository"`
	Commit commit `json:"commits"`
}

type repo struct {
	Name string `json:"name"`
	Url  string `json:"clone_url"`
}
type commit []struct {
	ID string `json:"id"`
}

func (rp RequestParse) WebhookRequestParse(payloadData []byte) error {
	var payload RequestParse
	err := json.Unmarshal(payloadData, &payload)
	if err != nil {
		fmt.Println("Failed to parse the request")
		return err
	}
	fmt.Println(payload.Commit[0].ID)
	fmt.Println(payload.Repo.Url)
	return nil
}
