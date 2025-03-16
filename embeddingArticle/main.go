package main

import (
	"fmt"

	"github.com/rajeshkio/embeddingArticle/address"
	"github.com/rajeshkio/embeddingArticle/person"
)

func main() {
	addr := address.Address{
		City:    "Mumbai",
		ZipCode: "231456",
	}
	p := person.Person{
		FirstName: "rajesh",
		LastName:  "kumar",
		Address:   addr,
	}

	fmt.Println(p.Fullname())
	fmt.Println(p.Address.Format())
}
