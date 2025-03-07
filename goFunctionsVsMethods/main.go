package main

import (
	"fmt"
	"os"
	"strconv"
)

func calculateArea(w, h float64) float64 {
	return w * h
}

func main() {

	if len(os.Args) != 3 {
		fmt.Println("Usage: go run program.go <width> <height>")
		fmt.Println("Example: go run program.go 5.0 3.0")
		os.Exit(1)
	}
	width, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		fmt.Println("Error: Width must be a valid number")
	}
	height, err := strconv.ParseFloat(os.Args[2], 64)
	if err != nil {
		fmt.Println("Error: height must be a valid number")
	}

	area := calculateArea(float64(width), float64(height))
	fmt.Printf("The area of the rectangle is: %.2f\n", area)
}
