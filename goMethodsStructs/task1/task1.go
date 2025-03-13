package task1

import (
	"fmt"

	"github.com/rajeshkio/goMethodsStructs/models"
)

func RunExample() {
	address := models.NewAddress("MH", "Mumbai", 123456)
	NewPerson := models.NewPerson(30, "Rajesh")

	fmt.Printf("Address created: %s\n", address.Format())

	person := models.NewPersonWithAddress(NewPerson, address)
	fmt.Printf("Person created: %s\n", person)

	fmt.Printf("Name: %s\n", person.GetName())
	fmt.Printf("Age: %d\n", person.GetAge())
	fmt.Printf("Address: %s\n", person.GetAddress())

	person.SetAge(31)
	fmt.Printf("New Age: %d\n", person.GetAge())

	newAddress := models.NewAddress("KA", "Bangalore", 124323)
	person.SetAddress(newAddress)
	fmt.Printf("New Address: %s\n", person.GetAddress())

	person.SetName("Raj")
	fmt.Printf("New name: %s\n", person.GetName())
}
