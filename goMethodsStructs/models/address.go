package models

import "fmt"

type Address struct {
	state string
	city  string
	zip   int
}

func NewAddress(state, city string, zip int) Address {
	return Address{
		state: state,
		city:  city,
		zip:   zip,
	}
}

func (a Address) GetState() string {
	return a.state
}

func (a Address) GetCity() string {
	return a.city
}

func (a Address) GetZip() int {
	return a.zip
}

func (a *Address) SetState(state string) {
	a.state = state
}

func (a *Address) SetCity(city string) {
	a.city = city
}

func (a *Address) SetZip(zip int) {
	a.zip = zip
}

func (a Address) Format() string {
	return fmt.Sprintf("%s, %s, %d", a.state, a.city, a.zip)
}

func (a Address) String() string {
	return fmt.Sprintf("Address(State: %s, City: %s, Zip: %d)", a.state, a.city, a.zip)
}
