package main

import "fmt"

// Simulates a function that returns multiple values
func getPlayerInfo() (string, int, bool) {
	return "alice", 25, true
}

// Simulates a function that might fail (returns error)
func connectToGame(gameID string) (string, error) {
	if gameID == "" {
		return "", fmt.Errorf("gameID cannot be empty")
	}
	return "Connected to game: " + gameID, nil
}

func main() {
	fmt.Println("=== Example 3: Multiple Return Values ===\n")

	// Pattern 1: Receive multiple values
	name, age, isActive := getPlayerInfo()
	fmt.Printf("Player: %s, Age: %d, Active: %t\n", name, age, isActive)

	// Pattern 2: Error handling (super common in Go!)
	message, err := connectToGame("game123")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Success: %s\n", message)
	}

	// Pattern 3: Try with empty gameID
	message, err = connectToGame("")
	if err != nil {
		fmt.Printf("Error: %v\n", err) // This will print
	} else {
		fmt.Printf("Success: %s\n", message)
	}

	// Pattern 4: Ignore values with _
	playerName, _, _ := getPlayerInfo() // Only care about name
	fmt.Printf("\nPlayer name only: %s\n", playerName)

	// Pattern 5: Reusing variables with :=
	// At least ONE variable must be new for := to work
	name2 := "bob"                   // name2 is new
	name2, score := "charlie", 100   // name2 exists, but score is new - OK!
	fmt.Printf("\nReused variable: %s, New variable: %d\n", name2, score)
}
