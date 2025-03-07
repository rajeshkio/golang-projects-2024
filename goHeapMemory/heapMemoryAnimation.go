package main

import (
	"fmt"
	"runtime"
)

func main() {
	// Allocate memory for a slice (dynamic array) of 10 integers
	memorySize := 10
	slice := make([]int, memorySize)

	// Initialize and use the allocated memory
	for i := 0; i < len(slice); i++ {
		slice[i] = 5 // Assigning a uniform value to each element
	}

	// Print the values
	for i := 0; i < len(slice); i++ {
		fmt.Printf("%d ", slice[i])
	}
	fmt.Println()

	// Force garbage collection to demonstrate memory deallocation
	runtime.GC()
}
