package main

import (
	"github.com/rk280392/customCICDTool/runServer"
	webhookhandler "github.com/rk280392/customCICDTool/webHookHandler"
	"github.com/rk280392/customCICDTool/webHookParser"
)

func main() {

	/* 	path := "/tmp/customCICDPath"
	   	url := "https://github.com/git-fixtures/basic.git"
	   	branch := "master"
	   	commitHash := ""

	   	os.RemoveAll(path)
	   	err := cloneRepo.CloneRepo(path, branch, url)
	   	if err != nil {
	   		fmt.Println("Failed to clone repo: ", url, branch)
	   	}
	   	err = cloneRepo.CheckoutCommit(commitHash, path, branch)
	   	if err != nil {
	   		fmt.Println("Failed to checkout")
	   	}
	   	fmt.Println("Repository cloned successfully")
	*/
	parser := &webHookParser.MyWebhookParser{}
	handler := webhookhandler.NewMyWebhookHandler(parser)

	server := runServer.NewServer(handler)
	server.Start()
	
}
