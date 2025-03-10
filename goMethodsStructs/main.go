package main

import (
	"fmt"

	"github.com/rajeshkio/goMethodsStructs/task1"
)

func main() {
	fmt.Println("=== Task 1: Basic Struct with Methods ===")

	person := task1.NewPerson(33, "rajesh", "Pune")

	fmt.Println("Initial Person:", person)

	person.SetAge(31)
	fmt.Println("Second Person:", person)
	person2Age := person.GetAge()
	fmt.Println("Second Person age:", person2Age)

	fmt.Printf("Person: %s, Age: %d, Address: %s", person.GetName(), person.GetAge(), person.GetAddress())

}
