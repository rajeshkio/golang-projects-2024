package person

import (
	"fmt"

	"github.com/rajeshkio/embeddingArticle/address"
)

type Person struct {
	FirstName       string
	LastName        string
	address.Address // this is embedding
}

func (p Person) Fullname() string {
	return fmt.Sprintf("Fullname is : %s %s", p.FirstName, p.LastName)
}

func (p Person) Fulladdress() string {
	return fmt.Sprintf("Fulladdress is : %s %s", p.Address.City, p.Address.ZipCode)
}
