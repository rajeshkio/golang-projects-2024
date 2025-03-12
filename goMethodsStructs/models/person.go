package models

import (
	"fmt"
)

/*
Task 1: Create a Basic Struct with Fields and Methods
Outline:

Define a Person struct with name, age, and address fields
Create methods to get and update these fields
Implement a String() method to format the struct data

Plain English Solution:
Create a Go file with a Person struct having string and integer fields. Add methods like GetName(), SetAge(), and a special String()
method (which Go uses automatically when printing). Test by creating a Person, modifying fields with your methods, and printing it.
*/
type PersonWithAddress struct {
	age  int
	name string
	Address
}

func NewPersonWithAddress(age int, name string, address Address) PersonWithAddress {
	return PersonWithAddress{
		age:     age,
		name:    name,
		Address: address,
	}
}

func (p PersonWithAddress) GetName() string {
	return p.name
}

func (p *PersonWithAddress) SetName(name string) {
	p.name = name
}

func (p PersonWithAddress) GetAge() int {
	return p.age
}
func (p *PersonWithAddress) SetAge(age int) {
	p.age = age
}

func (p PersonWithAddress) GetAddress() Address {
	return p.Address
}

func (p *PersonWithAddress) SetAddress(address Address) {
	p.Address = address
}

func (p PersonWithAddress) String() string {
	return fmt.Sprintf("Person(name: %s, age: %d, address: %s)", p.name, p.age, p.Address.Format())
}

func (p PersonWithAddress) GetFullAddress() string {
	return p.Address.Format()
}

func (p PersonWithAddress) GetCity() string {
	return p.Address.GetCity()
}

func (p *PersonWithAddress) UpdateAddressState(state string) {
	addr := p.Address
	addr.state = state
	p.Address = addr
}
