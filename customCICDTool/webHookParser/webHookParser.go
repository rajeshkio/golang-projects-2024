package webHookParser

import (
	"encoding/json"
	"fmt"
)

type RepositoryInfo struct {
	Repository string
	Branch     string
	Commit     string
	Url        string
}

type MyWebhookParser struct{}

type WebhookPayload struct {
	Ref        string `json:"ref"`
	Repository struct {
		Name string `json:"name"`
		Url  string `json:"clone_url"`
	} `json:"repository"`
	Commit []struct {
		Id string `json:"id"`
	} `json:"commits"`
}

func (p *MyWebhookParser) Parse(payloadData []byte) (RepositoryInfo, error) {
	var payload WebhookPayload

	if err := json.Unmarshal(payloadData, &payload); err != nil {
		return RepositoryInfo{}, err
	}
	fmt.Println("Parsing payload...")
	repoName := payload.Repository.Name
	branchName := payload.Ref
	commitID := payload.Commit[0].Id
	repoUrl := payload.Repository.Url

	return RepositoryInfo{
		Repository: repoName,
		Branch:     branchName,
		Commit:     commitID,
		Url:        repoUrl,
	}, nil
}

/*
TODO right now the code is tight coupled making any changes in the halder or parser need to rerun the server code.
Need to understand and decouple it.

*/
