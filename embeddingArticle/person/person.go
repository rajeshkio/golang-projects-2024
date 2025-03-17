package person

import (
	"fmt"

	"github.com/rajeshkio/embeddingArticle/address"
	"github.com/rajeshkio/embeddingArticle/contact"
)

type Person struct {
	FirstName       string
	LastName        string
	Age             int
	address.Address // this is embedding
	contact.Contact
}

func (p Person) Fullname() string {
	return fmt.Sprintf("Fullname is : %s %s", p.FirstName, p.LastName)
}

func (p Person) Fulladdress() string {
	return fmt.Sprintf("%s %v", p.Address.City, p.Address.ZipCode)
}

func (p Person) Greet() string {
	return fmt.Sprintf("Hi I am %s, and I live in %s and my contact information is %s", p.Fullname(), p.Fulladdress(), p.ContactInfo())
}

func (p Person) String() string {
	return fmt.Sprintf("%s, %d - %s", p.Fullname(), p.Age, p.Fulladdress())
}
