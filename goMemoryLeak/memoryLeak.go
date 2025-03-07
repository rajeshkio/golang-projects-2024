package main

import (
	"fmt"
	"time"
)

// Converts bytes to a human-readable string with appropriate units.
func humanReadableSize(bytes int) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

func main() {
	var leakySlice []int
	incrementSize := 10000000 // Number of elements to add each iteration

	for i := 0; ; i++ {
		for j := 0; j < incrementSize; j++ {
			leakySlice = append(leakySlice, j)
		}
		lengthInBytes := len(leakySlice) * 4 // Each int is 4 bytes
		capacityInBytes := cap(leakySlice) * 4
		fmt.Printf("Length: %s, Capacity: %s\n", humanReadableSize(lengthInBytes), humanReadableSize(capacityInBytes))
		time.Sleep(1 * time.Second) // Slow down the loop for easier monitoring
	}
}
