package task1

import (
	"fmt"

	person "github.com/rajeshkio/goMethodsStructs/models"
)

func RunExample() {
	p := person.NewPerson(
		30,
		"Rajesh Kumar",
		"Mumbai",
	)
	fmt.Println("Person created: ", p)
	fmt.Printf("Name: %s\n", p.GetName())
	fmt.Printf("Age: %d\n", p.GetAge())
	fmt.Printf("Address: %s\n", p.GetAddress())

	p.SetAge(31)
	fmt.Printf("New Age: %d\n", p.GetAge())

	p.SetAddress("Pune")
	fmt.Printf("New Address: %s\n", p.GetAddress())

	p.SetName("Raj")
	fmt.Printf("New name: %s\n", p.GetName())
}
