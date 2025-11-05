package main

import "fmt"

// TODO: Complete the exercises below
// Run with: go run practice.go

// Exercise 1: Fix the errors in this function
func exercise1() {
	fmt.Println("\n=== Exercise 1: Fix the Errors ===")

	// TODO: Fix the following code
	// Hint: There are 3 errors

	var playerName string
	playerName = "alice"

	// This should create a new variable 'score'
	// playerName := "bob" // ERROR: playerName already exists
	// Fix: Use = for assignment, := only for new variables

	score := 100
	// score := 200 // ERROR: score already declared
	// Fix: Use = instead of :=

	fmt.Printf("Player: %s, Score: %d\n", playerName, score)

	// TODO: Declare a variable 'level' with value 5
	// Write your code here:
	var level int = 5

	fmt.Printf("Level: %d\n", level)
}

// Exercise 2: Complete the function
func exercise2() {
	fmt.Println("\n=== Exercise 2: Multiple Return Values ===")

	// This function returns (username, email, age)
	name, email, age := getUserData()

	// TODO: Print the values
	// Expected output: "User: alice, Email: alice@example.com, Age: 25"
	// Write your code here:
	fmt.Printf("User: %s, Email: %s, Age: %d\n", name, email, age)

	// TODO: Call getUserData again but only get the username
	// Ignore the other values using _
	// Write your code here:
	username, _, _ := getUserData()

	fmt.Printf("Username only: %s\n", username)
}

// Helper function for exercise 2
func getUserData() (string, string, int) {
	return "alice", "alice@example.com", 25
}

// Exercise 3: Error handling
func exercise3() {
	fmt.Println("\n=== Exercise 3: Error Handling ===")

	// TODO: Call validateMove with "e2e4" and handle the result
	// If error is nil, print "Valid move: [move]"
	// If error is not nil, print "Invalid: [error]"
	// Write your code here:

	move, err := validateMove("e2e4")
	if err != nil {
		fmt.Printf("Invalid: %v\n", err)
	} else {
		fmt.Printf("Valid move: %s\n", move)
	}

	// TODO: Call validateMove with empty string "" and handle the error
	// Write your code here:

	move, err = validateMove("")
	if err != nil {
		fmt.Printf("Invalid: %v\n", err)
	} else {
		fmt.Printf("Valid move: %s\n", move)
	}
}

// Helper function for exercise 3
func validateMove(move string) (string, error) {
	if move == "" {
		return "", fmt.Errorf("move cannot be empty")
	}
	return move, nil
}

// Exercise 4: Zero values
func exercise4() {
	fmt.Println("\n=== Exercise 4: Zero Values ===")

	// TODO: Declare variables without initialization and print their zero values
	// Declare: wins (int), playerName (string), isActive (bool)
	// Write your code here:

	var wins int
	var playerName string
	var isActive bool

	fmt.Printf("wins: %d (should be 0)\n", wins)
	fmt.Printf("playerName: '%s' (should be empty)\n", playerName)
	fmt.Printf("isActive: %t (should be false)\n", isActive)

	// TODO: Increment wins 3 times and print
	// Write your code here:

	wins++
	wins++
	wins++

	fmt.Printf("After 3 games, wins: %d\n", wins)
}

// Exercise 5: Real websocket pattern
func exercise5() {
	fmt.Println("\n=== Exercise 5: Websocket Pattern ===")

	// TODO: Complete the websocket connection pattern
	// 1. Call connectToServer("player1")
	// 2. Check for errors
	// 3. If no error, print "Connected: [connection id]"
	// 4. If error, print "Failed: [error]" and return
	// Write your code here:

	conn, err := connectToServer("player1")
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
		return
	}
	fmt.Printf("Connected: %s\n", conn)

	// TODO: Call receiveMessage and handle the 3 return values
	// Print: "Received: type=[type], data=[data]"
	// Write your code here:

	msgType, data, err := receiveMessage()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Received: type=%d, data=%s\n", msgType, string(data))
}

// Helper functions for exercise 5
func connectToServer(playerID string) (string, error) {
	if playerID == "" {
		return "", fmt.Errorf("playerID required")
	}
	return "conn_" + playerID, nil
}

func receiveMessage() (int, []byte, error) {
	return 1, []byte("Hello!"), nil
}

// Bonus Exercise: Challenge yourself!
func bonusExercise() {
	fmt.Println("\n=== Bonus: Create Your Own ===")

	// TODO: Write a function that:
	// 1. Takes a player name and score
	// 2. Returns a formatted string and an error
	// 3. Returns error if name is empty or score is negative
	// 4. Otherwise returns "Player [name] scored [score] points"

	// Write your function here:
	formatPlayerScore := func(name string, score int) (string, error) {
		if name == "" {
			return "", fmt.Errorf("name cannot be empty")
		}
		if score < 0 {
			return "", fmt.Errorf("score cannot be negative")
		}
		return fmt.Sprintf("Player %s scored %d points", name, score), nil
	}

	// Test your function:
	result, err := formatPlayerScore("alice", 100)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println(result)
	}

	// Test with invalid input:
	result, err = formatPlayerScore("", 50)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println(result)
	}
}

func main() {
	fmt.Println("╔════════════════════════════════════════╗")
	fmt.Println("║   Lesson 1: Practice Exercises         ║")
	fmt.Println("╚════════════════════════════════════════╝")

	// Uncomment each exercise as you complete it:
	exercise1()
	exercise2()
	exercise3()
	exercise4()
	exercise5()
	bonusExercise()

	fmt.Println("\n✅ All exercises complete!")
	fmt.Println("\nNext: Create SUMMARY.md with your learnings")
}
