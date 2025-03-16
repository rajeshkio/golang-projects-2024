package address

import "fmt"

type Address struct {
	City    string
	ZipCode string
}

func (a Address) Format() string {
	return fmt.Sprintf("%s %s", a.City, a.ZipCode)
}
