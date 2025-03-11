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
package main

type Address struct {
	street string
	city   string
	zip    string
}

func test2() {

}
