package contact

import "fmt"

type Contact struct {
	Email string
	Phone int64
}

func (c Contact) ContactInfo() string {
	return fmt.Sprintf("Email: %s, Phone: %v", c.Email, c.Phone)
}
