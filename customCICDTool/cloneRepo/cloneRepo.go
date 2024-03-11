package cloneRepo

import (
	"fmt"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func CloneRepo(path, branch, url string) error {
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:           url,
		ReferenceName: plumbing.ReferenceName(branch),
	})
	if err != nil {
		return err
	}
	return nil
}

func CheckoutCommit(commitHash, path, branch string) error {

	readFile, err := git.PlainOpen(path)
	if err != nil {
		fmt.Println("Failed to open the git repo")
		return err
	}

	worktree, err := readFile.Worktree()
	if err != nil {
		fmt.Println("Failed to get worktree:")
		return err
	}
	commit := plumbing.NewHash(commitHash)
	if commit != plumbing.ZeroHash {
		err = worktree.Checkout(&git.CheckoutOptions{
			Hash: commit,
		})
	} else {
		headRef, err := readFile.Head()
		if err != nil {
			fmt.Println("Failed to get HEAD reference:")
			return err
		}
		err = worktree.Checkout(&git.CheckoutOptions{
			Hash: headRef.Hash(),
		})
		if err != nil {
			fmt.Println("Failed to checkout to head of branch: ")
			return err
		}
	}
	if err != nil {
		fmt.Println("Failed to checkout commit: ")
		return err
	}
	return nil
}
