# Lesson 2: Quick Start Guide

Get up and running in 2 minutes! ⚡

---

## 1. Start the Server

```bash
cd /Users/ankit/code/learn/chess-coach/learnings/fundamentals/websocket/lesson-02-connection-management
go run main.go
```

You should see:

```
╔════════════════════════════════════════════════╗
║   Lesson 2: Connection Management             ║
║   Production-Ready WebSocket Server           ║
╚════════════════════════════════════════════════╝

⚙️  Configuration:
   • Ping interval:    54s
   • Pong timeout:     1m0s
   • Write timeout:    10s
   • Max message size: 1 MB

🚀 Server starting on http://localhost:8080
```

**Keep this terminal open** - you'll watch ping/pong messages here!

---

## 2. Open the Client

Open your browser to: **http://localhost:8080/client.html**

You'll see a beautiful client with:
- Connection status indicator
- Message statistics
- Test buttons
- Message log

---

## 3. Connect

Click the big green **"Connect"** button.

You should see:
- Status changes to "Connected" (green dot)
- Welcome message appears
- Connection time starts counting

In your **server terminal**, you'll see:
```
📞 New connection request from 127.0.0.1:xxxxx
✅ Client 127.0.0.1:xxxxx-timestamp registered. Total connections: 1
📤 Sent to client: Welcome!
```

---

## 4. Send a Message

Type "Hello WebSocket!" and click **Send** (or press Enter).

**Browser shows**:
```
📤 Hello WebSocket!
📨 Server echo: Hello WebSocket! (at 18:45:23)
```

**Server terminal shows**:
```
📨 Client xxx sent TEXT: Hello WebSocket!
📤 Sent to client xxx: Server echo: Hello WebSocket! (at 18:45:23)
```

---

## 5. Watch the Heartbeat

**Wait 54 seconds** (you can send more messages while waiting - heartbeat is independent!)

After 54 seconds, in your **server terminal**, you'll see:

```
🏓 Sent PING to client xxx
🏓 Received PONG from client xxx
```

This happens **automatically every 54 seconds** to keep the connection healthy!

---

## 6. Test Automatic Reconnection

Click the purple button: **"Simulate Network Issue"**

Watch what happens:
1. Connection drops (red dot)
2. Status shows "Reconnecting..."
3. Reconnection attempts: 1, 2, 3...
4. Delays increase: 1s → 2s → 4s...
5. Connection restored (green dot)

This is **exponential backoff** in action!

---

## 7. Test Message Queuing

1. Click "Simulate Network Issue" again
2. **Quickly** type and send 3-4 messages while "Reconnecting" shows
3. Watch the "Queued" counter increase
4. Wait for reconnection
5. See all queued messages get sent automatically!

**No message loss!** 🎉

---

## 8. Test Graceful Shutdown

1. Make sure you're connected
2. Go to your **server terminal**
3. Press **Ctrl+C**

In the **browser**, you'll see:
```
🔌 Disconnected (code: 1001, reason: Server shutting down)
🔄 Server going away - will reconnect shortly
⏳ Reconnecting in 1s (attempt 1/10)
```

The server told the client **why** it's closing (code 1001 = "Going Away")!

---

## Quick Reference

### What to Watch

**Server Terminal**:
- 📞 New connections
- 🏓 PING/PONG every 54 seconds
- 📨 Received messages
- 📤 Sent messages
- 👋 Disconnections
- 🛑 Shutdown messages

**Browser Client**:
- Connection status (colored dot)
- Message statistics
- Queued messages
- Reconnection attempts
- All messages sent/received

### Test Buttons

| Button | What It Does |
|--------|--------------|
| **Send Test Message** | Sends timestamped test message |
| **Simulate Network Issue** | Forces disconnection (code 1006) |
| **Send 10 Messages (Burst)** | Tests rapid message handling |
| **Auto-Reconnect: ON/OFF** | Toggle automatic reconnection |

### Keyboard Shortcuts

- **Enter** in message input = Send message
- **Ctrl+C** in server terminal = Graceful shutdown

---

## Experiments

Ready for hands-on experiments? Check out:

👉 **[EXPERIMENTS.md](EXPERIMENTS.md)** - 8 detailed experiments

Each experiment teaches you something specific about connection management:
1. Normal Connection & Heartbeat
2. Zombie Connection Detection
3. Automatic Reconnection
4. Message Queuing
5. Graceful Shutdown
6. Timeout Detection
7. Stress Test - Multiple Connections
8. Burst Messages

---

## Key Concepts to Observe

As you experiment, notice:

1. **Heartbeats are independent from messages**
   - Pings fire every 54s even if you're actively sending messages
   - They run in the background to keep connection healthy

2. **Reconnection uses exponential backoff**
   - First attempt: 1s
   - Second attempt: 2s
   - Third attempt: 4s
   - Prevents overwhelming server!

3. **Messages are queued during disconnection**
   - Type messages while disconnected
   - They're automatically sent when reconnected
   - No data loss!

4. **Different close codes mean different things**
   - 1000 = Normal closure (don't reconnect)
   - 1001 = Server going away (reconnect soon)
   - 1006 = Abnormal closure (network issue, reconnect)

---

## Troubleshooting

### Server won't start

**Error**: "address already in use"

**Solution**: Another program is using port 8080
```bash
# Find what's using port 8080
lsof -i :8080

# Kill it (use PID from above command)
kill -9 <PID>

# Or change port in main.go line 299
server.Addr = ":8081"  // Use different port
```

### Client won't connect

**Check**:
1. Is server running? (Check terminal)
2. Is URL correct? (ws://localhost:8080/ws)
3. Check browser console (F12) for errors

### No ping/pong messages

**Wait longer!** First ping happens at **54 seconds** after connection.

Set a timer and wait - it will come! ⏱️

---

## What's Next?

1. **Try all experiments** in [EXPERIMENTS.md](EXPERIMENTS.md)
2. **Write down questions** in [QUESTIONS.md](QUESTIONS.md)
3. **Read the theory** in [README.md](README.md)
4. **Review implementation** in [SUMMARY.md](SUMMARY.md)

When you're ready for the next lesson, say: **"Move to Lesson 3"**

---

**Have fun experimenting!** 🚀
