/*
Task 2: Embedding Structs (Composition)
Outline:

Create an Address struct
Embed it within your Person struct
Create methods for both structs
Observe field and method promotion

Plain English Solution:
First, separate the address into its own struct with fields like street, city, and zip code.
Then modify your Person struct to include an Address field. Notice how you can access Address fields directly
from Person (field promotion) and how methods work with nested structs. This demonstrates Go's composition approach versus inheritance.
*/
package task2

import (
	"fmt"

	"github.com/rajeshkio/goMethodsStructs/models"
)

func RunExample() {
	address := models.NewAddress("MH", "Pune", 412110)
	fmt.Println("Address created: ", address)

	newPerson := models.NewPersonWithAddress(31, "Rajesh", address)
	fmt.Println("new person created: ", newPerson)

	fmt.Printf("%s aged %d lives in %s\n", newPerson.GetName(), newPerson.GetAge(), newPerson.GetFullAddress())

	oldAddress := newPerson.GetAddress()
	newAddress := oldAddress
	newAddress.SetState("KA")
	newAddress.SetCity("Bangalore")
	newAddress.SetZip(321234)

	newPerson.SetAddress(newAddress)

	fmt.Printf("%s moved to %s from %s", newPerson.GetName(), newPerson.GetFullAddress(), oldAddress.Format())
}
