package main

import "fmt"

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
    // Create rectangle instances
    rect1 := Rectangle{width: 5.0, height: 3.0}
    rect2 := Rectangle{width: 10.0, height: 7.0}
    
    // Calculate properties for Rectangle 1
    area1 := calculateArea(rect1)
    perimeter1 := calculatePerimeter(rect1)
    
    // Calculate properties for Rectangle 2
    area2 := calculateArea(rect2)
    perimeter2 := calculatePerimeter(rect2)
    
    fmt.Printf("Rectangle 1: Area = %.2f, Perimeter = %.2f\n", area1, perimeter1)
    fmt.Printf("Rectangle 2: Area = %.2f, Perimeter = %.2f\n", area2, perimeter2)
}
