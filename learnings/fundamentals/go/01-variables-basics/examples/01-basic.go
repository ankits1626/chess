package main

import "fmt"

func main() {
	fmt.Println("=== Example 1: Basic Variable Declarations ===\n")

	// Method 1: Declare variable, then assign
	var username string
	username = "alice"
	fmt.Printf("Username: %s\n", username)

	// Method 2: Declare with initial value
	var age int = 25
	fmt.Printf("Age: %d\n", age)

	// Method 3: Let Go infer the type
	var score = 100 // Go knows this is int
	fmt.Printf("Score: %d\n", score)

	// Multiple variables of same type
	var x, y, z int = 1, 2, 3
	fmt.Printf("Coordinates: x=%d, y=%d, z=%d\n", x, y, z)

	// Zero values (default when not initialized)
	var count int      // 0
	var name string    // ""
	var isActive bool  // false
	fmt.Printf("\nZero values:\n")
	fmt.Printf("  count: %d\n", count)
	fmt.Printf("  name: '%s'\n", name)
	fmt.Printf("  isActive: %t\n", isActive)
}
