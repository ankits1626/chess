package main

import "fmt"

func main() {
	fmt.Println("=== Example 2: Short Declaration with := ===\n")

	// Short declaration - Go figures out the type
	username := "bob"      // string
	age := 30              // int
	height := 5.9          // float64
	isPlayer := true       // bool

	fmt.Printf("Player Info:\n")
	fmt.Printf("  Username: %s (type: string)\n", username)
	fmt.Printf("  Age: %d (type: int)\n", age)
	fmt.Printf("  Height: %.1f (type: float64)\n", height)
	fmt.Printf("  Is Player: %t (type: bool)\n", isPlayer)

	// You can reassign with = (not :=)
	age = 31 // ✅ OK - assignment
	// age := 32 // ❌ Would be ERROR - age already declared

	fmt.Printf("\nAfter birthday:\n")
	fmt.Printf("  Age: %d\n", age)

	// Multiple short declarations
	x, y := 10, 20
	fmt.Printf("\nCoordinates: x=%d, y=%d\n", x, y)

	// Swapping is easy!
	x, y = y, x
	fmt.Printf("After swap: x=%d, y=%d\n", x, y)
}
