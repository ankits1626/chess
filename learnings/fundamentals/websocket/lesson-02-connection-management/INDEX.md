# Lesson 2: Connection Management & Reliability - Complete Guide

Welcome to Lesson 2! This directory contains everything you need to master production-ready WebSocket connection management.

---

## 📚 Documentation Guide

### 🚀 Getting Started (Read These First!)

1. **[QUICKSTART.md](QUICKSTART.md)** - Get running in 2 minutes
   - How to start the server
   - How to connect the client
   - First tests to try
   - What to watch for

2. **[README.md](README.md)** - Theory & Concepts (16KB)
   - Why connection management matters
   - How ping/pong heartbeats work
   - Connection lifecycle states
   - Close event codes
   - Production patterns

### 🧪 Hands-On Learning

3. **[EXPERIMENTS.md](EXPERIMENTS.md)** - 8 Interactive Experiments (12KB)
   - Experiment 1: Normal Connection & Heartbeat
   - Experiment 2: Zombie Connection Detection
   - Experiment 3: Automatic Reconnection
   - Experiment 4: Message Queuing
   - Experiment 5: Graceful Shutdown
   - Experiment 6: Timeout Detection
   - Experiment 7: Multiple Connections
   - Experiment 8: Burst Messages
   - **Time**: 30-45 minutes

4. **[QUESTIONS.md](QUESTIONS.md)** - Your Learning Journal
   - Write questions as you experiment
   - Template provided
   - Discussion starters

### 📖 Deep Dive

5. **[ARCHITECTURE.md](ARCHITECTURE.md)** - Visual Architecture Guide (20KB)
   - System overview diagrams
   - Connection lifecycle flow
   - Goroutine communication
   - Reconnection state machine
   - Message flow diagrams
   - Timeout management
   - Connection registry internals

6. **[SUMMARY.md](SUMMARY.md)** - Implementation Summary (15KB)
   - What we built
   - Code architecture
   - Key constants & configuration
   - Production patterns explained
   - Before/after comparison
   - Key learnings

---

## 📁 Source Files

### Server Implementation

- **[main.go](main.go)** (10KB) - Production-ready WebSocket server
  - Ping/pong heartbeats (54s interval)
  - Read/write deadlines (60s/10s)
  - Connection registry
  - Graceful shutdown
  - Concurrent connection handling

### Client Implementation

- **[client.html](client.html)** (25KB) - Interactive browser client
  - Automatic reconnection with exponential backoff
  - Message queuing during disconnection
  - Connection health monitoring
  - Test utilities
  - Beautiful UI

### Dependencies

- **go.mod** - Go module definition
- **go.sum** - Dependency checksums

---

## 🎯 Learning Path

### For Complete Beginners

```
Day 1:
├─ Read: QUICKSTART.md (5 min)
├─ Run: Server + Client (5 min)
├─ Try: Experiments 1-3 (20 min)
└─ Read: README.md sections 1-3 (20 min)

Day 2:
├─ Try: Experiments 4-8 (30 min)
├─ Read: ARCHITECTURE.md (20 min)
└─ Review: Code in main.go (20 min)

Day 3:
├─ Read: SUMMARY.md (20 min)
├─ Write: Questions in QUESTIONS.md (10 min)
└─ Experiment: Create your own tests (30 min)
```

### For Experienced Developers

```
Session 1 (1 hour):
├─ Skim: README.md (10 min)
├─ Review: main.go + client.html (20 min)
├─ Run: All experiments (20 min)
└─ Read: ARCHITECTURE.md (10 min)

Session 2 (30 min):
├─ Read: SUMMARY.md (15 min)
├─ Modify: Change timeouts and test (10 min)
└─ Plan: How to apply to your project (5 min)
```

---

## 🎓 What You'll Learn

By completing this lesson, you'll understand:

### Core Concepts

- ✅ **Heartbeats (Ping/Pong)**
  - Why they're essential
  - How they work in practice
  - When to send them
  - How to handle timeouts

- ✅ **Connection Lifecycle**
  - States: CONNECTING, OPEN, CLOSING, CLOSED
  - Graceful vs ungraceful disconnection
  - Close event codes and their meanings
  - Connection cleanup

- ✅ **Automatic Reconnection**
  - Why exponential backoff matters
  - How to prevent server overload
  - When to reconnect vs when to stop
  - User experience considerations

- ✅ **Message Queuing**
  - How to preserve messages during disconnection
  - Queue limits and overflow handling
  - Order preservation
  - Flush strategies

### Production Patterns

- ✅ **Timeouts & Deadlines**
  - Read deadlines (detecting stalled reads)
  - Write deadlines (detecting stalled writes)
  - Pong deadlines (detecting dead connections)
  - Balancing responsiveness vs false positives

- ✅ **Concurrent Architecture**
  - Two goroutines per connection (readPump + writePump)
  - Channel-based communication
  - Why separation matters
  - Thread-safety without excessive locking

- ✅ **Connection Registry**
  - Tracking active connections
  - Broadcast capabilities
  - Graceful shutdown
  - Connection limits

- ✅ **Graceful Shutdown**
  - Catching signals (SIGINT/SIGTERM)
  - Closing connections with reason
  - Resource cleanup
  - Client behavior on shutdown

---

## 🛠️ How to Use This Lesson

### Quick Start (15 minutes)

1. Read [QUICKSTART.md](QUICKSTART.md)
2. Start server: `go run main.go`
3. Open client: http://localhost:8080/client.html
4. Try Experiments 1 & 2

### Full Experience (2-3 hours)

1. **Understand the theory** (30 min)
   - Read [README.md](README.md)
   - Focus on diagrams and examples

2. **See it in action** (45 min)
   - Do all 8 experiments in [EXPERIMENTS.md](EXPERIMENTS.md)
   - Watch server logs carefully
   - Take notes of observations

3. **Understand the code** (45 min)
   - Read [ARCHITECTURE.md](ARCHITECTURE.md) with diagrams
   - Review [main.go](main.go) with comments
   - Review [client.html](client.html) ReconnectingWebSocket class

4. **Synthesize learning** (30 min)
   - Read [SUMMARY.md](SUMMARY.md)
   - Write questions in [QUESTIONS.md](QUESTIONS.md)
   - Plan how to apply to your project

### For Your Chess-Coach Project

After this lesson, you'll be ready to:

- ✅ Build reliable game connections
- ✅ Handle player disconnections gracefully
- ✅ Preserve moves during network blips
- ✅ Detect zombie connections
- ✅ Implement graceful server restarts

---

## 📊 Code Statistics

```
Total Lines: ~1,200
Server: ~350 lines (Go)
Client: ~650 lines (HTML/JavaScript)
Documentation: ~3,000 lines (Markdown)

Features Implemented:
├─ Ping/pong heartbeats
├─ Read/write deadlines
├─ Automatic reconnection
├─ Exponential backoff
├─ Message queuing
├─ Connection registry
├─ Graceful shutdown
├─ Concurrent handling
└─ Error recovery

Patterns Covered:
├─ Two-goroutine per connection
├─ Channel-based communication
├─ Event loop (connection registry)
├─ State machine (client reconnection)
├─ Deadline management
└─ Resource cleanup
```

---

## 🔍 Quick Reference

### Server Configuration

```go
writeWait = 10 * time.Second      // Write timeout
pongWait  = 60 * time.Second      // Pong timeout
pingPeriod = (pongWait * 9) / 10  // 54 seconds
maxMessageSize = 1 * 1024 * 1024  // 1 MB
```

### Client Configuration

```javascript
maxReconnectAttempts = 10         // Stop after 10 tries
reconnectDelay = 5000             // Start with 5s (gives time to queue messages)
// Exponential backoff: 5s → 10s → 20s → 30s (max)
```

### WebSocket Close Codes

| Code | Meaning | Reconnect? |
|------|---------|-----------|
| 1000 | Normal closure | ❌ No |
| 1001 | Going away (server restart) | ✅ Yes |
| 1006 | Abnormal closure (network) | ✅ Yes |
| 1008 | Policy violation | ❌ No |

### Endpoints

- Home: http://localhost:8080/
- WebSocket: ws://localhost:8080/ws
- Client: http://localhost:8080/client.html

---

## 🐛 Troubleshooting

### Server won't start
- **"address already in use"** → Port 8080 is taken
  - Solution: `lsof -i :8080` then `kill -9 <PID>`
  - Or change port in main.go

### Client won't connect
- Check if server is running
- Check browser console (F12) for errors
- Verify URL: ws://localhost:8080/ws

### No ping/pong messages
- Wait at least 54 seconds!
- Check server logs (not client)
- Should see "🏓 Sent PING" every 54s

### Build fails
- Missing dependencies: `go mod download`
- Wrong directory: `cd lesson-02-connection-management`

---

## 📝 Files at a Glance

| File | Size | Purpose | Read When |
|------|------|---------|-----------|
| QUICKSTART.md | 6KB | Get started fast | First! |
| README.md | 16KB | Theory & concepts | After quick start |
| EXPERIMENTS.md | 12KB | Hands-on practice | After theory |
| ARCHITECTURE.md | 20KB | Deep dive | After experiments |
| SUMMARY.md | 15KB | Implementation review | After code |
| QUESTIONS.md | 1KB | Your notes | As you go |
| INDEX.md | 6KB | This file | For navigation |
| main.go | 10KB | Server code | Read with architecture |
| client.html | 25KB | Client code | Read with architecture |

---

## ✅ Completion Checklist

Mark these off as you complete them:

### Understanding
- [ ] I understand why ping/pong heartbeats are necessary
- [ ] I know how exponential backoff prevents server overload
- [ ] I can explain the two-goroutine pattern
- [ ] I understand read/write deadlines
- [ ] I know the difference between close codes 1000, 1001, 1006

### Practical
- [ ] I've run the server successfully
- [ ] I've connected the client
- [ ] I've completed all 8 experiments
- [ ] I've watched ping/pong messages in logs
- [ ] I've tested automatic reconnection
- [ ] I've seen message queuing work
- [ ] I've tested graceful shutdown

### Code
- [ ] I've read main.go with understanding
- [ ] I've read client.html ReconnectingWebSocket class
- [ ] I can modify the timeout values
- [ ] I understand the connection registry pattern

### Application
- [ ] I can explain how this applies to chess-coach
- [ ] I have ideas for extending this code
- [ ] I've written my questions in QUESTIONS.md

---

## 🎉 What's Next?

After completing Lesson 2, you're ready for:

### Lesson 3: Message Protocols & State Management
- Structured communication (JSON, Protocol Buffers)
- Request/Response patterns
- Event-based architecture
- Message validation
- State synchronization

### Lesson 4: Authentication & Authorization
- WebSocket authentication
- Token-based auth
- Session management
- Permission systems

### Lesson 5: Scaling & Production
- Horizontal scaling
- Load balancing
- Redis for state sharing
- Monitoring & metrics
- Performance optimization

---

## 💬 Need Help?

If you have questions:

1. Write them in [QUESTIONS.md](QUESTIONS.md)
2. Say "I have questions" to discuss
3. Review [ARCHITECTURE.md](ARCHITECTURE.md) diagrams
4. Check troubleshooting section above

---

## 📚 Additional Resources

From [README.md](README.md):
- WebSocket RFC 6455 specification
- Gorilla WebSocket documentation
- Go concurrency patterns
- Exponential backoff algorithms

---

**Ready to start?** Go to [QUICKSTART.md](QUICKSTART.md)! 🚀

**Good luck and have fun learning!** 🎓
