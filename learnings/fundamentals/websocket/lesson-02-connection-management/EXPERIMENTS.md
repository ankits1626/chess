# Lesson 2: Connection Management - Experiments

**Goal**: Test and understand connection reliability, heartbeats, and reconnection behavior

**Time**: 30-45 minutes

---

## Before You Start

1. **Start the server**:
   ```bash
   cd /Users/ankit/code/learn/chess-coach/learnings/fundamentals/websocket/lesson-02-connection-management
   go run main.go
   ```

2. **Open the client**: http://localhost:8080/client.html

3. **Keep an eye on**:
   - Terminal (server logs) - you'll see 🏓 PING/PONG messages!
   - Browser console (F12) - for client-side logs
   - Client UI - status indicators and message log

---

## Experiment 1: Normal Connection & Heartbeat

**What you'll learn**: How ping/pong heartbeats work in practice

### Steps:

1. Click **Connect** in the browser
2. Watch the **server terminal** closely
3. Wait for **54 seconds** (you can keep sending messages while waiting)
4. Observe what happens

### Expected Behavior:

```
Server terminal:
✅ Client 127.0.0.1:xxxxx-timestamp registered. Total connections: 1
📤 Sent to client: Welcome! You are client...
[... wait 54 seconds ...]
🏓 Sent PING to client 127.0.0.1:xxxxx-timestamp
🏓 Received PONG from client 127.0.0.1:xxxxx-timestamp
[... 54 seconds later ...]
🏓 Sent PING to client 127.0.0.1:xxxxx-timestamp
🏓 Received PONG from client 127.0.0.1:xxxxx-timestamp
```

Browser shows:
- Connection time incrementing
- Messages being sent/received normally
- **No disconnections** despite 54 second gaps!

### Questions to Answer:

1. ✅ Did you see PING messages in the server logs every 54 seconds?
2. ✅ Did the client stay connected even with no user activity?
3. ✅ What happens between pings? (Can you still send messages?)

### Key Insight:

**Heartbeats run independently from your messages!** You can send messages anytime, and pings happen on schedule regardless.

---

## Experiment 2: Zombie Connection Detection

**What you'll learn**: How the server detects dead connections

### Steps:

1. Connect from the browser
2. Wait for **first ping/pong** (54 seconds) - confirm it works
3. **Close the browser tab** (don't click Disconnect - just close the tab!)
4. Watch the **server terminal**
5. Wait up to **60 seconds**

### Expected Behavior:

```
Server terminal:
🏓 Sent PING to client...
[... no pong received ...]
[... after 60 seconds total without pong ...]
❌ Failed to send PING to client: websocket: close sent
👋 Client disconnected normally
👋 Client xxx unregistered. Total connections: 0
```

### What's Happening:

```
Timeline:
0s   - Browser tab closed
54s  - Server sends PING
55s  - Server waits for PONG... (none comes)
60s  - Read deadline expires (pongWait timeout)
60s  - Server detects connection is dead
60s  - Server cleans up resources
```

### Questions to Answer:

1. ✅ How long did it take for the server to detect the disconnection?
2. ✅ Why didn't it detect instantly when you closed the browser?
3. ✅ What would happen if we didn't have heartbeats?

### Key Insight:

**Without heartbeats, the server might never know the client disconnected!** TCP doesn't always notify immediately when a connection breaks.

---

## Experiment 3: Automatic Reconnection

**What you'll learn**: How exponential backoff works

### Steps:

1. Connect from the browser
2. In the **server terminal**, press **Ctrl+C** (stop the server)
3. Watch the **browser client**
4. Observe the reconnection attempts and delays
5. **Restart the server** after a few attempts
6. Watch what happens

### Expected Behavior:

```
Browser shows:
⚠️ Connection lost abnormally (code: 1006)
⏳ Reconnecting in 1s (attempt 1/10)
❌ Connection failed
⏳ Reconnecting in 2s (attempt 2/10)
❌ Connection failed
⏳ Reconnecting in 4s (attempt 3/10)
❌ Connection failed
⏳ Reconnecting in 8s (attempt 4/10)
[... you restart server here ...]
✅ Connected successfully!
📤 Sending X queued messages...
```

### Reconnection Pattern:

```
Attempt 1: Wait 1 second
Attempt 2: Wait 2 seconds
Attempt 3: Wait 4 seconds
Attempt 4: Wait 8 seconds
Attempt 5: Wait 16 seconds
Attempt 6+: Wait 30 seconds (max)
```

### Questions to Answer:

1. ✅ Did the delays double each time? (1s → 2s → 4s → 8s)
2. ✅ Why use exponential backoff instead of trying every second?
3. ✅ What happens to messages you send while disconnected?

### Key Insight:

**Exponential backoff prevents overwhelming the server** when it comes back online. If 1000 clients all reconnected every second, the server would crash again!

---

## Experiment 4: Message Queuing

**What you'll learn**: How messages are preserved during disconnection

### Steps:

1. Connect from the browser
2. Send a test message - confirm it works
3. Click **"Simulate Network Issue"** button
4. Connection will drop (code 1006 - abnormal closure)
5. **Quickly type and send 3-5 messages** while "Reconnecting" shows
6. Watch the queue count increase
7. Wait for reconnection
8. Observe what happens to queued messages

### Expected Behavior:

```
Browser shows:
✅ Connected
📤 Test message 1 (sent immediately)
⚠️ Connection lost
⏳ Message queued: Hello
⏳ Message queued: World
⏳ Message queued: Are you there?
[Queue count: 3]
🔄 Reconnecting...
✅ Connected successfully!
📤 Sending 3 queued messages...
📤 [Queued] Hello
📤 [Queued] World
📤 [Queued] Are you there?
📨 Server echo: Hello
📨 Server echo: World
📨 Server echo: Are you there?
```

### Questions to Answer:

1. ✅ Were all queued messages sent after reconnection?
2. ✅ Were they sent in the correct order?
3. ✅ What happens if you send 100 messages while disconnected?

### Key Insight:

**Message queuing ensures no data loss during temporary disconnections!** Critical for chess moves that happen during network blips.

---

## Experiment 5: Graceful Shutdown

**What you'll learn**: Difference between graceful and ungraceful shutdown

### Setup: Open TWO browser tabs with the client

**Tab 1**: Normal connection
**Tab 2**: Also connected

### Steps:

1. Both tabs connected - verify by checking "Total connections: 2" in server logs
2. In server terminal, press **Ctrl+C**
3. Watch **both browser tabs**
4. Look at the close codes and messages

### Expected Behavior:

```
Server terminal:
🛑 Shutdown signal received...
🛑 Shutting down... Closing 2 connections
✅ Server stopped gracefully

Browser Tab 1 & 2:
🔌 Disconnected (code: 1001, reason: Server shutting down)
🔄 Server going away - will reconnect shortly
⏳ Reconnecting in 1s (attempt 1/10)
```

**Close code 1001 = "Going Away"** - server told clients it's shutting down intentionally

### Compare with Experiment 3:

- **Ungraceful (server crash)**: Code 1006 (Abnormal Closure)
- **Graceful (Ctrl+C)**: Code 1001 (Going Away)

### Questions to Answer:

1. ✅ What close code did you see?
2. ✅ Did clients try to reconnect?
3. ✅ How is this different from a crash?

### Key Insight:

**Graceful shutdown provides context to clients** - they know the server is restarting, not crashing, and can adjust their reconnection strategy.

---

## Experiment 6: Timeout Detection

**What you'll learn**: How read/write deadlines protect the server

### Steps:

1. Connect from browser
2. Open browser **Developer Tools** (F12)
3. Go to **Console** tab
4. Paste this code to freeze the pong responses:

```javascript
// Override WebSocket to block pongs
const originalWS = reconnectingWS.ws;
originalWS.send = function(data) {
    console.log('BLOCKED: Attempted to send:', data);
    // Don't actually send - simulate frozen client
};
```

5. Wait for **60+ seconds**
6. Watch server terminal

### Expected Behavior:

```
Server terminal:
🏓 Sent PING to client...
[... no pong received ...]
[... 60 seconds pass ...]
❌ Client xxx: unexpected close error: ...
👋 Client xxx unregistered. Total connections: 0
```

### What's Happening:

```
Server side:
- Read deadline set: 60 seconds from last pong
- Sent PING at 54s
- Expected PONG within 60s total
- No PONG received
- Deadline expires
- Connection terminated
```

### Questions to Answer:

1. ✅ How long did the server wait before closing?
2. ✅ What would happen without this timeout?
3. ✅ Why is 60 seconds a good timeout value?

### Key Insight:

**Deadlines prevent zombie connections from consuming server resources!** Without them, dead connections would stay "open" forever.

---

## Experiment 7: Stress Test - Multiple Connections

**What you'll learn**: How the server handles concurrent connections

### Steps:

1. Open **4-5 browser tabs/windows** with the client
2. Connect all of them
3. Check server logs for "Total connections"
4. Send messages from **different tabs**
5. Watch the server terminal
6. Close one tab - observe connection count
7. Press Ctrl+C on server - watch all tabs reconnect

### Expected Behavior:

```
Server terminal:
✅ Client A registered. Total connections: 1
✅ Client B registered. Total connections: 2
✅ Client C registered. Total connections: 3
✅ Client D registered. Total connections: 4
📨 Client A sent TEXT: Hello from A
📨 Client B sent TEXT: Hello from B
🏓 Sent PING to client A
🏓 Sent PING to client B
🏓 Sent PING to client C
🏓 Sent PING to client D
[... close one tab ...]
👋 Client B disconnected normally
👋 Client B unregistered. Total connections: 3
```

### Questions to Answer:

1. ✅ Can the server handle multiple connections?
2. ✅ Are pings sent independently to each client?
3. ✅ What happens when you close one connection?

### Key Insight:

**Each client has its own goroutines** (readPump + writePump) running independently. This is why Go's concurrency model is perfect for WebSocket servers!

---

## Experiment 8: Burst Messages

**What you'll learn**: How send buffers handle rapid messages

### Steps:

1. Connect from browser
2. Click **"Send 10 Messages (Burst)"** button
3. Watch the logs - both browser and server
4. Observe timing and order

### Expected Behavior:

```
Browser (very fast):
📤 Burst message 1/10
📤 Burst message 2/10
📤 Burst message 3/10
...
📤 Burst message 10/10

Server (receives all):
📨 Client xxx sent TEXT: Burst message 1/10
📨 Client xxx sent TEXT: Burst message 2/10
...
📨 Client xxx sent TEXT: Burst message 10/10

Browser (receives echoes):
📨 Server echo: Burst message 1/10
📨 Server echo: Burst message 2/10
...
```

### Questions to Answer:

1. ✅ Were all 10 messages received?
2. ✅ Were they in order?
3. ✅ How fast were they sent?

### Key Insight:

**WebSocket maintains message order and handles bursts efficiently** through buffering. The `send` channel (buffer: 256) queues messages for the writePump goroutine.

---

## Challenge Experiments

### Challenge 1: Maximum Queue Size

**Hypothesis**: What happens if you queue 1000 messages?

1. Disconnect (stop server)
2. Send 1000 messages with burst
3. Reconnect
4. What happens?

**Think about**: Should there be a queue limit? Why?

---

### Challenge 2: Very Long Timeout

Change server code:
```go
pongWait = 10 * time.Minute  // 10 minutes
pingPeriod = 9 * time.Minute  // 9 minutes
```

Rebuild and test. What changes?

---

### Challenge 3: No Automatic Reconnection

In client, click **"Auto-Reconnect: OFF"**

Then simulate network issue. What happens?

---

## Summary Checklist

After completing all experiments, you should understand:

- ✅ How ping/pong heartbeats detect dead connections
- ✅ Why exponential backoff prevents server overload
- ✅ How message queuing prevents data loss
- ✅ Difference between graceful and ungraceful shutdown
- ✅ How timeouts protect server resources
- ✅ How concurrent connections work
- ✅ How WebSocket handles message bursts

---

## What's Next?

In the next lesson, you'll learn:
- **Message protocols** - structured communication (JSON, Protocol Buffers)
- **Authentication** - securing WebSocket connections
- **Room/Channel patterns** - broadcasting to groups (like chess games!)
- **State synchronization** - keeping multiple clients in sync

---

## Questions File

As you do these experiments, write down questions in `QUESTIONS.md`:
- Anything that surprised you
- Behavior you don't understand
- "What if..." scenarios

We'll discuss them together!

**Happy experimenting!** 🎉
