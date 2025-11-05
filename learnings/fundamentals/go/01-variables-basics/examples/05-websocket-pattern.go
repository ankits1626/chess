package main

import "fmt"

// Simulating websocket patterns you'll see in chess-coach

// Mock types (we don't have real websocket here)
type Connection struct {
	id string
}

type Message struct {
	Type int
	Data string
}

// Simulates websocket upgrade
func upgradeConnection(clientID string) (*Connection, error) {
	if clientID == "" {
		return nil, fmt.Errorf("client ID required")
	}
	return &Connection{id: clientID}, nil
}

// Simulates reading a message
func readMessage(conn *Connection) (int, []byte, error) {
	// Returns: messageType, message, error
	return 1, []byte("hello from " + conn.id), nil
}

// Simulates receiving from a channel
func receiveFromChannel(messages chan string) (string, bool) {
	// In real code: message, ok := <-c.send
	// Returns: message, ok (ok=false if channel closed)
	select {
	case msg := <-messages:
		return msg, true
	default:
		return "", false
	}
}

func main() {
	fmt.Println("=== Example 5: Websocket Patterns ===\n")

	// Pattern 1: Upgrade connection (2 return values)
	// Real: conn, err := upgrader.Upgrade(w, r, nil)
	conn, err := upgradeConnection("player1")
	if err != nil {
		fmt.Printf("Connection failed: %v\n", err)
		return
	}
	fmt.Printf("✅ Connection established: %s\n\n", conn.id)

	// Pattern 2: Read message (3 return values)
	// Real: messageType, message, err := conn.ReadMessage()
	messageType, message, err := readMessage(conn)
	if err != nil {
		fmt.Printf("Read failed: %v\n", err)
		return
	}
	fmt.Printf("📨 Received message:\n")
	fmt.Printf("   Type: %d\n", messageType)
	fmt.Printf("   Data: %s\n\n", string(message))

	// Pattern 3: Channel receive (2 return values)
	// Real: message, ok := <-c.send
	messages := make(chan string, 2)
	messages <- "move: e2e4"
	messages <- "move: e7e5"

	msg, ok := receiveFromChannel(messages)
	if ok {
		fmt.Printf("📬 From channel: %s\n", msg)
	} else {
		fmt.Println("Channel is empty or closed")
	}

	// Another receive
	msg, ok = receiveFromChannel(messages)
	if ok {
		fmt.Printf("📬 From channel: %s\n", msg)
	}

	// Pattern 4: Reusing 'err' variable
	// Very common in Go - declare once, reuse with =
	conn2, err := upgradeConnection("player2") // err is reused from line 38
	if err != nil {
		fmt.Printf("\n❌ Error: %v\n", err)
	} else {
		fmt.Printf("\n✅ Second connection: %s\n", conn2.id)
	}

	fmt.Println("\n✅ All patterns demonstrated!")
	fmt.Println("\nKey takeaways:")
	fmt.Println("  • := creates new variables")
	fmt.Println("  • Multiple return values are common")
	fmt.Println("  • err is almost always the last return value")
	fmt.Println("  • You can reuse variables if at least one is new")
}
