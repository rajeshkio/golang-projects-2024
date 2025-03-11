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
type Person struct {
	age     int
	name    string
	address string
}

func NewPerson(age int, name, address string) Person {
	return Person{
		age:     age,
		name:    name,
		address: address,
	}
}
func (p Person) GetName() string {
	return p.name
}

func (p *Person) SetName(name string) {
	p.name = name
}

func (p Person) GetAge() int {
	return p.age
}
func (p *Person) SetAge(age int) {
	p.age = age
}

func (p Person) GetAddress() string {
	return p.address
}

func (p *Person) SetAddress(address string) {
	p.address = address
}

func (p Person) String() string {
	return fmt.Sprintf("Person(name: %s, age: %d, address: %s)", p.name, p.age, p.address)
}
