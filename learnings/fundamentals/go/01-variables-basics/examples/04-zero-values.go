package main

import "fmt"

func main() {
	fmt.Println("=== Example 4: Zero Values ===\n")
	fmt.Println("When you declare without initializing, Go assigns 'zero values':\n")

	// Numeric types → 0
	var count int
	var price float64
	var score int32
	fmt.Printf("int:     %d\n", count)
	fmt.Printf("float64: %.1f\n", price)
	fmt.Printf("int32:   %d\n", score)

	// String → empty string ""
	var username string
	var gameID string
	fmt.Printf("\nstring:  '%s' (empty)\n", username)
	fmt.Printf("string:  '%s' (empty)\n", gameID)

	// Boolean → false
	var isActive bool
	var hasWon bool
	fmt.Printf("\nbool:    %t\n", isActive)
	fmt.Printf("bool:    %t\n", hasWon)

	// Pointer → nil (we'll learn this in Lesson 2!)
	var ptr *int
	fmt.Printf("\npointer: %v (nil means 'points to nothing')\n", ptr)

	// This is useful for initializing counters
	var wins int // starts at 0
	wins++       // now 1
	wins++       // now 2
	fmt.Printf("\nWins after 2 games: %d\n", wins)

	// Check if string is empty
	var playerName string
	if playerName == "" {
		fmt.Println("\nNo player name set (zero value)")
	}
}
