package task3

import (
	"fmt"

	"github.com/rajeshkio/goMethodsStructs/models"
)

func RunExample() {

	fmt.Println("Simple Person Example")
	person := models.NewPerson(34, "Rajesh")

	fmt.Println("Original Person: ", person)

	// Simple receiver example
	updatedPerson := models.UpdateAge(person, 31)
	fmt.Printf("After UpdateAge operation: original: %d, returned: %d\n", person.Age, updatedPerson.Age)

	// Pointer receiver example
	models.UpdateAgePointer(&person, 32)
	fmt.Printf("After UpdateAgePointer operation: original: %d\n", person.Age)

}
