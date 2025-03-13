package main

import (
	"fmt"

	"github.com/rajeshkio/goMethodsStructs/task1"
	"github.com/rajeshkio/goMethodsStructs/task2"
	"github.com/rajeshkio/goMethodsStructs/task3"
)

func main() {
	fmt.Println("=== Task 1: Basic Struct with Methods ===")

	task1.RunExample()

	fmt.Println("=== Task 2: Embeddings in Struct with Methods ===")

	task2.RunExample()

	fmt.Println("=== Task 3: Pointers and Receivers")

	task3.RunExample()

}
