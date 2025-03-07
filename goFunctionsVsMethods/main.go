package main

import (
	"fmt"
	"math"
)

type Rectangle struct {
	width  float64
	height float64
	x      float64 // X position
	y      float64 // Y position
	color  string  // Color
}

func calculateArea(rect Rectangle) float64 {
	return rect.width * rect.height
}

func calculatePerimeter(rect Rectangle) float64 {
	return 2 * (rect.width + rect.height)
}

func calculateDiagonal(rect Rectangle) float64 {
	return math.Sqrt(rect.width*rect.width + rect.height*rect.height)
}

func main() {
	// Create rectangle instances
	rect1 := Rectangle{width: 5.0, height: 3.0, x: 0.0, y: 0.0, color: "blue"}
	rect2 := Rectangle{width: 10.0, height: 7.0, x: 3.0, y: 2.0, color: "red"}

	// Calculate properties for Rectangle 1
	area1 := calculateArea(rect1)
	perimeter1 := calculatePerimeter(rect1)
	diagonal1 := calculateDiagonal(rect1)

	// Calculate properties for Rectangle 2
	area2 := calculateArea(rect2)
	perimeter2 := calculatePerimeter(rect2)
	diagonal2 := calculateDiagonal(rect2)

	fmt.Printf("Rectangle 1: Area = %.2f, Perimeter = %.2f, Diagonal = %.2f\n", area1, perimeter1, diagonal1)
	fmt.Printf("Rectangle 2: Area = %.2f, Perimeter = %.2f, Diagonal = %.2f\n", area2, perimeter2, diagonal2)
}
