package main

import (
	"fmt"
	"math"
)

type Rectangle struct {
	width  float64
	height float64
	x      float64
	y      float64
	color  string
}

// These methods don't modify the rectangle, so we use value receivers
func (r Rectangle) Area() float64 {
	return r.width * r.height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.width + r.height)
}

func (r Rectangle) Diagonal() float64 {
	return math.Sqrt(r.width*r.width + r.height*r.height)
}

func (r Rectangle) Description() string {
	return fmt.Sprintf("A %s rectangle at position (%.1f,%.1f) with dimensions %.1fx%.1f",
		r.color, r.x, r.y, r.width, r.height)
}

// These methods modify the rectangle, so we use pointer receivers
func (r *Rectangle) Scale(factor float64) {
	r.width *= factor
	r.height *= factor
}

func (r *Rectangle) Move(dx, dy float64) {
	r.x += dx
	r.y += dy
}

func (r *Rectangle) ChangeColor(newColor string) {
	r.color = newColor
}

func main() {
	// Create rectangle
	rect := Rectangle{
		width:  5.0,
		height: 3.0,
		x:      0.0,
		y:      0.0,
		color:  "blue",
	}

	// Print original rectangle properties
	fmt.Printf("Original rectangle: %s\n", rect.Description())
	fmt.Printf("Area: %.2f, Perimeter: %.2f, Diagonal: %.2f\n",
		rect.Area(), rect.Perimeter(), rect.Diagonal())

	// Modify the rectangle
	rect.Scale(2.0)           // Double the size
	rect.Move(1.0, 1.0)       // Move it 1 unit right and 1 unit down
	rect.ChangeColor("green") // Change color to green

	// Print modified rectangle properties
	fmt.Printf("Modified rectangle: %s\n", rect.Description())
	fmt.Printf("New area: %.2f, New perimeter: %.2f, New diagonal: %.2f\n",
		rect.Area(), rect.Perimeter(), rect.Diagonal())
}
