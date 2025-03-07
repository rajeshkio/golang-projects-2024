package main

import (
	"fmt"
	"os"
	"strconv"
)

type Rectangle struct {
	width  float64
	height float64
}

func calculateArea(rect Rectangle) float64 {
	return rect.width * rect.height
}

func calculatePerimeter(rect Rectangle) float64 {
	return 2 * (rect.width + rect.height)
}

func main() {

	if len(os.Args) < 3 || (len(os.Args)-1)%2 != 0 {
		fmt.Println("Usage: go run program.go <width1> <height1> [<width2> <height2> ...]")
		fmt.Println("Example: go run program.go 5.0 3.0 10.0 7.0 15.0 12.0")
		fmt.Println("You must provide width and height pairs for each rectangle")
		os.Exit(1)
	}

	numOfRectangles := (len(os.Args) - 1) / 2
	rectangles := make([]Rectangle, numOfRectangles)

	for i := range numOfRectangles {
		widthPos := i*2 + 1
		heightPos := i*2 + 2

		width, err := strconv.ParseFloat(os.Args[widthPos], 64)
		if err != nil {
			fmt.Printf("Width %d  should be a valid number: ", i+1)
		}
		height, err := strconv.ParseFloat(os.Args[heightPos], 64)
		if err != nil {
			fmt.Printf("Height %d  should be a valid number: ", i+1)
		}
		/*
			creating rectangle and adding to slice
			Initial state: [{0 0} {0 0} {0 0}]
			After iteration 0: [{10 5} {0 0} {0 0}]
			After iteration 1: [{10 5} {20 10} {0 0}]
			After iteration 2: [{10 5} {20 10} {30 15}]
		*/
		rectangles[i] = Rectangle{width: width, height: height}
	}

	// calculate and show properties of each rectangle
	for i, rect := range rectangles {
		area := calculateArea(rect)
		perimeter := calculatePerimeter(rect)

		fmt.Printf("Rectangle %d: Area = %.2f  Perimeter = %.2f\n", i, area, perimeter)
	}
}
