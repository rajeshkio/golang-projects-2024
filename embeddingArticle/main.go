package main

import (
	"fmt"

	"github.com/rajeshkio/embeddingArticle/address"
	"github.com/rajeshkio/embeddingArticle/contact"
	"github.com/rajeshkio/embeddingArticle/person"
)

func main() {
	addr := address.Address{
		City:    "Mumbai",
		ZipCode: 231456,
	}
	c := contact.Contact{
		Email: "rk90229@gmail.com",
		Phone: 1234567890,
	}
	p := person.Person{
		FirstName: "rajesh",
		LastName:  "kumar",
		Age:       27,
		Address:   addr,
		Contact:   c,
	}

	//fmt.Println(p.Fullname())
	//fmt.Println(p.Format())
	//fmt.Println(p.Greet())

	// Testing name conflicts
	fmt.Println(p.String()) // Only prints the Person.string() method and not Address.String() method.

	fmt.Println(p.String(), p.Address.String()) // We can access the embedded method using the field name: Address.String()
}
