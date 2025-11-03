# Lesson 1 - Experiments

**Goal**: Get hands-on experience with WebSocket basics

**Before starting**: Make sure the server is running (`go run main.go`)

---

## Experiment 1: Connect and Send Messages

### Steps:
1. Open `http://localhost:8080/client.html` in your browser
2. Click **Connect**
3. Watch the server terminal - see the connection log?
4. Type "Hello World" and click **Send**
5. Observe both the client and server logs

### Questions to Answer:
- What happens in the server terminal when you connect?
- What message does the server send first?
- When you send "Hello World", what do you receive back?
- How fast is the response? (instant or delayed?)

---

## Experiment 2: Multiple Messages

### Steps:
1. Stay connected
2. Send these messages rapidly one after another:
   - "Message 1"
   - "Message 2"
   - "Message 3"
   - "Message 4"
   - "Message 5"

### Questions to Answer:
- Do all messages arrive?
- Are they in order?
- Does the connection stay open between messages?
- Check the server logs - what do you see for each message?

---

## Experiment 3: Multiple Clients

### Steps:
1. Keep the first client connected
2. Open a **second browser tab**
3. Open `http://localhost:8080/client.html` again
4. Connect the second client
5. Send messages from both clients

### Questions to Answer:
- Can both clients connect simultaneously?
- Does one client see messages from the other? (Why/why not?)
- Check server logs - can you identify each client?
- What happens if you disconnect one client? Does the other stay connected?

---

## Experiment 4: Connection Lifecycle

### Steps:
1. Connect to the server
2. Watch the server terminal carefully
3. Click **Disconnect**
4. Observe the logs

Now:
5. Connect again
6. Close the browser tab (don't click disconnect)
7. Check server terminal

### Questions to Answer:
- What's the difference between "Disconnect" button and closing the tab?
- Does the server detect both types of disconnection?
- How long does it take to detect a closed tab?

---

## Experiment 5: What Happens When Server Restarts?

### Steps:
1. Connect a client
2. Send a message - works fine
3. In the server terminal, press **Ctrl+C** to stop the server
4. Try to send another message from the client
5. What error appears?
6. Restart the server: `go run main.go`
7. Try to send from the same client (don't refresh browser)
8. Now refresh the browser and reconnect

### Questions to Answer:
- Does the client know immediately when the server stops?
- Can the client auto-reconnect? (In this simple implementation)
- What does the browser console show? (Open Developer Tools → Console)

---

## Experiment 6: Modify the Code

### Task 1: Change the Welcome Message

**Find this in main.go (line ~40):**
```go
welcomeMsg := fmt.Sprintf("Welcome! Connected at %s", time.Now().Format("15:04:05"))
```

**Change to:**
```go
welcomeMsg := "🎮 Chess Coach WebSocket - Ready for real-time games!"
```

**Steps:**
1. Save the file
2. Restart the server
3. Reconnect the client
4. What's different?

---

### Task 2: Add Message Counter

**Add this after line 40 (after welcomeMsg):**
```go
messageCount := 0
```

**Then inside the for loop (around line 56), after receiving a message, add:**
```go
messageCount++
log.Printf("📊 Total messages from this client: %d", messageCount)
```

**Steps:**
1. Save and restart server
2. Connect and send several messages
3. Watch the server count messages
4. Does the counter reset when you disconnect and reconnect? (Why?)

---

### Task 3: Make the Echo Uppercase

**Find line ~61:**
```go
responseMsg := fmt.Sprintf("Echo: %s (received at %s)", string(message), time.Now().Format("15:04:05"))
```

**Change to:**
```go
import "strings" // Add at the top
...
responseMsg := fmt.Sprintf("ECHO: %s (received at %s)", strings.ToUpper(string(message)), time.Now().Format("15:04:05"))
```

**Steps:**
1. Save and restart
2. Send "hello websocket"
3. What do you receive back?

---

## Experiment 7: Break Things (Learn from Errors!)

### Test 1: What if CheckOrigin returns false?

**In main.go, find line ~14:**
```go
CheckOrigin: func(r *http.Request) bool {
    return true
},
```

**Change to:**
```go
CheckOrigin: func(r *http.Request) bool {
    return false  // Reject all connections!
},
```

**Try connecting** - What error do you see?

This is CORS protection. In production, you'd check the origin:
```go
CheckOrigin: func(r *http.Request) bool {
    origin := r.Header.Get("Origin")
    return origin == "https://yourfrontend.com"
},
```

**Remember to change it back to `true` when done!**

---

### Test 2: What if you never call conn.Close()?

**In main.go, comment out line ~27:**
```go
// defer conn.Close()
```

**What happens?** Resources leak! Always close connections.

**Remember to uncomment it!**

---

### Test 3: Send a really large message

**In the client:**
1. Connect
2. Open browser console (F12)
3. Run:
```javascript
const largeMsg = 'x'.repeat(1000000); // 1 MB of 'x'
ws.send(largeMsg);
```

**What happens?** Check server logs. Our buffer is 1024 bytes!

This shows why we need `maxMessageSize` limits (we'll add in future lessons).

---

## Experiment 8: Browser Developer Tools

### Steps:
1. Connect the client
2. Open Developer Tools (F12)
3. Go to **Network** tab
4. Filter by **WS** (WebSocket)
5. Click on the WebSocket connection
6. Click **Messages** sub-tab

### What You See:
- All messages sent/received
- Timestamps
- Message content
- Frame types

### Try:
1. Send a message
2. Watch it appear in real-time in the Messages tab
3. See the green (sent) and white (received) arrows?

---

## Experiment 9: Browser Console WebSocket

### Manual WebSocket from Console:

Open browser console and paste:

```javascript
// Connect
const ws = new WebSocket('ws://localhost:8080/ws');

// Log connection
ws.onopen = () => console.log('✅ Connected!');

// Log messages
ws.onmessage = (e) => console.log('📨 Received:', e.data);

// Log errors
ws.onerror = (e) => console.error('❌ Error:', e);

// Log close
ws.onclose = () => console.log('🔌 Disconnected');

// Send a message (wait for connection first!)
setTimeout(() => {
    ws.send('Hello from console!');
}, 1000);
```

**Watch the console and server terminal!**

---

## Experiment 10: Compare with HTTP

### HTTP Version:

Open browser console and run:
```javascript
// HTTP Request
fetch('http://localhost:8080/')
  .then(r => r.text())
  .then(data => console.log('HTTP Response length:', data.length));
```

**Note:**
- Opens connection
- Gets response
- **Closes immediately**

### WebSocket Version:

```javascript
const ws = new WebSocket('ws://localhost:8080/ws');
ws.onopen = () => {
    console.log('WebSocket open');
    ws.send('Message 1');
    setTimeout(() => ws.send('Message 2'), 1000);
    setTimeout(() => ws.send('Message 3'), 2000);
};
```

**Note:**
- Opens once
- **Stays open**
- Can send multiple messages

**See the difference?**

---

## Challenge Experiments

### Challenge 1: Message Rate
How many messages can you send per second? Try:
```javascript
const ws = new WebSocket('ws://localhost:8080/ws');
ws.onopen = () => {
    let count = 0;
    const interval = setInterval(() => {
        ws.send(`Message ${++count}`);
        if (count >= 100) clearInterval(interval);
    }, 10); // Every 10ms = 100 msg/sec
};
```

Does it handle them all?

---

### Challenge 2: Long-Running Connection
Connect and leave it open for 5 minutes. Does it stay connected?
(It should - but we'll add ping/pong in Lesson 7 to guarantee this)

---

### Challenge 3: Binary Messages
Can this server handle binary data? Try:
```javascript
const ws = new WebSocket('ws://localhost:8080/ws');
ws.onopen = () => {
    const buffer = new Uint8Array([72, 101, 108, 108, 111]); // "Hello"
    ws.send(buffer);
};
```

Check server logs - what message type does it receive?

---

## Reflection Questions

After doing these experiments, answer:

1. **When is the WebSocket connection actually established?**
   - During `new WebSocket()`?
   - Or during `ws.onopen`?

2. **Can you send messages before `onopen` fires?**
   - Try it! What happens?

3. **What's the lifetime of a WebSocket connection?**
   - Until you close it?
   - Until the page refreshes?
   - Forever?

4. **How is this different from HTTP?**
   - List 3 key differences you observed

5. **For your chess game, which scenarios need WebSocket?**
   - Move updates?
   - Loading game history?
   - User profile?

---

## Document Your Findings

Write your observations in `SUMMARY.md`:
- What surprised you?
- What clicked?
- Any "aha!" moments?
- Questions that came up?

---

**Next**: When you've finished experimenting, say **"I'm done with experiments"** and we'll discuss your findings and create your summary!
