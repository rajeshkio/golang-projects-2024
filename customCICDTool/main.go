package main

import (
	"io"
	"log"
	"os"

	memfs "github.com/go-git/go-billy/v5/memfs"
	git "github.com/go-git/go-git/v5"
	memory "github.com/go-git/go-git/v5/storage/memory"
)

func main() {
	fs := memfs.New()
	storer := memory.NewStorage()

	_, err := git.Clone(storer, fs, &git.CloneOptions{
		URL: "https://github.com/git-fixtures/basic.git",
	})
	if err != nil {
		log.Fatal(err)
	}

	changelog, err := fs.Open("CHANGELOG")
	if err != nil {
		log.Fatal(err)
	}
	io.Copy(os.Stdout, changelog)
}
