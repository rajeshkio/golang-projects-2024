package address

import "fmt"

type Address struct {
	City    string
	ZipCode int64
	State   string
}

func (a Address) Format() string {
	return fmt.Sprintf("%s %v", a.City, a.ZipCode)
}

// Testing name conflicts

func (a Address) String() string {
	return fmt.Sprintf("%s, %v, %s", a.City, a.ZipCode, a.State)
}
