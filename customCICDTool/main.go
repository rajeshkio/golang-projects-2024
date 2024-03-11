package main

import (
	"fmt"
	"os"

	"github.com/rk280392/customCICDTool/cloneRepo"
)

func main() {

	path := "/tmp/customCICDPath"
	url := "https://github.com/git-fixtures/basic.git"
	branch := "refs/heads/master"

	os.RemoveAll(path)
	err := cloneRepo.CloneRepo(path, branch, url)
	if err != nil {
		fmt.Println("Failed to clone repo: ", url, branch)
	}
	err = cloneRepo.CheckoutCommit("", path, branch)
	if err != nil {
		fmt.Println("Failed to checkout")
	}
	fmt.Println("Repository cloned successfully")

}
