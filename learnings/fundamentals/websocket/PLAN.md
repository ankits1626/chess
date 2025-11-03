# WebSocket Fundamentals Learning Plan (Go)

**Goal**: Master WebSocket concepts to understand and implement the chess-coach real-time architecture

**Target Project**: Chess Coach - Computer play, multiplayer, AI coaching via WebSocket

---

## Learning Structure

Each lesson follows this pattern:
1. **Concept Introduction** - Theory and explanation
2. **Code Example** - Runnable Go program
3. **Experiments** - Things to try and modify
4. **Q&A Session** - Ask questions about the lesson
5. **Summary Document** - Key takeaways and notes
6. **Next Steps** - What we'll build on

---

## Lesson Roadmap

### Phase 1: Core WebSocket Fundamentals (Lessons 1-4)
Build foundation for understanding WebSocket protocol and Go implementation

### Phase 2: Production Patterns (Lessons 5-7)
Learn patterns used in real applications like chess-coach

### Phase 3: Chess-Coach Architecture (Lessons 8-10)
Apply everything to understand your project's implementation

---

## Lessons Overview

### **Lesson 1: What is WebSocket?**
**Folder**: `lesson-01-basics/`
**Concepts**:
- HTTP vs WebSocket (why real-time matters)
- WebSocket handshake and upgrade
- Full-duplex bidirectional communication
- When to use WebSocket vs REST/SSE

**Deliverable**: Simple Go server that accepts WebSocket connections
**Chess-Coach Connection**: Why your project needs WebSocket for moves

---

### **Lesson 2: Your First WebSocket Connection**
**Folder**: `lesson-02-first-connection/`
**Concepts**:
- Gorilla WebSocket package
- Upgrading HTTP to WebSocket
- Reading and writing messages
- Connection lifecycle (open, message, close)

**Deliverable**: Echo server - sends back whatever you send
**Chess-Coach Connection**: Basic client-server communication pattern

---

### **Lesson 3: Message Formats and JSON**
**Folder**: `lesson-03-messages/`
**Concepts**:
- Text vs Binary messages
- JSON message structure
- Message typing and routing
- Error handling

**Deliverable**: Server that routes different message types
**Chess-Coach Connection**: `WSMessage` structure with type and payload

---

### **Lesson 4: Concurrent Connections**
**Folder**: `lesson-04-concurrency/`
**Concepts**:
- Go goroutines for WebSocket
- Read and Write pumps
- Channels for communication
- Graceful connection cleanup

**Deliverable**: Server handling multiple clients concurrently
**Chess-Coach Connection**: `ReadPump` and `WritePump` in your Client

---

### **Lesson 5: The Hub Pattern**
**Folder**: `lesson-05-hub-pattern/`
**Concepts**:
- Central connection manager
- Client registration/unregistration
- Broadcasting messages
- Thread-safe operations with mutex

**Deliverable**: Chat server with Hub managing all connections
**Chess-Coach Connection**: Your `Hub` struct managing game connections

---

### **Lesson 6: Room and Channel Management**
**Folder**: `lesson-06-rooms/`
**Concepts**:
- Organizing clients into rooms/channels
- Joining and leaving rooms
- Broadcasting to specific rooms
- Excluding sender from broadcasts

**Deliverable**: Multi-room chat (like Discord channels)
**Chess-Coach Connection**: Game rooms where players connect

---

### **Lesson 7: Ping/Pong and Connection Health**
**Folder**: `lesson-07-health/`
**Concepts**:
- Detecting dead connections
- Ping/Pong protocol
- Timeouts and deadlines
- Automatic reconnection strategies

**Deliverable**: Server with connection health monitoring
**Chess-Coach Connection**: `pingPeriod`, `pongWait` in your Client

---

### **Lesson 8: Request-Response Pattern**
**Folder**: `lesson-08-request-response/`
**Concepts**:
- Async request-response over WebSocket
- Correlation IDs for matching responses
- Handling slow operations
- Timeouts and error responses

**Deliverable**: Server that processes requests and sends responses
**Chess-Coach Connection**: Move request → Computer move response

---

### **Lesson 9: Integration with Services**
**Folder**: `lesson-09-services/`
**Concepts**:
- Connecting WebSocket to business logic
- Context propagation
- Service layer integration
- Database operations from WebSocket handlers

**Deliverable**: WebSocket server integrated with mock game service
**Chess-Coach Connection**: `EngineService` integration for move generation

---

### **Lesson 10: Chess-Coach Architecture Deep Dive**
**Folder**: `lesson-10-chess-coach/`
**Concepts**:
- Complete architecture review
- Message flow diagrams
- Hub → Client → Handler → Service chain
- Computer play implementation
- Multiplayer readiness
- AI coach streaming (preview)

**Deliverable**: Annotated walkthrough of chess-coach WebSocket code
**Chess-Coach Connection**: Full understanding of your project

---

## Lesson Structure Template

Each lesson folder contains:

```
lesson-XX-name/
├── README.md           # Concept explanation and theory
├── main.go            # Complete runnable example
├── client.html        # Web client for testing (if needed)
├── EXPERIMENTS.md     # Things to try and modify
├── SUMMARY.md         # Your notes after completing lesson
└── QUESTIONS.md       # Your questions during learning
```

---

## How to Use This Plan

### Step-by-Step Process:

1. **Start with Lesson 1**
   - Read `README.md` thoroughly
   - Understand the concept before coding

2. **Run the Example**
   - `cd lesson-01-basics/`
   - `go run main.go`
   - Test with provided client

3. **Do Experiments**
   - Follow `EXPERIMENTS.md`
   - Modify code to test understanding
   - Break things and fix them

4. **Ask Questions**
   - Write questions in `QUESTIONS.md`
   - We discuss and clarify

5. **Create Summary**
   - I help you write `SUMMARY.md`
   - Key concepts in your own words
   - Code snippets that clicked

6. **Move to Next Lesson**
   - Each builds on previous
   - Don't skip ahead

---

## Progress Tracking

| Lesson | Status | Time Spent | Key Insight |
|--------|--------|------------|-------------|
| 01 - Basics | ⬜ Not Started | - | - |
| 02 - First Connection | ⬜ Not Started | - | - |
| 03 - Messages | ⬜ Not Started | - | - |
| 04 - Concurrency | ⬜ Not Started | - | - |
| 05 - Hub Pattern | ⬜ Not Started | - | - |
| 06 - Rooms | ⬜ Not Started | - | - |
| 07 - Health | ⬜ Not Started | - | - |
| 08 - Request-Response | ⬜ Not Started | - | - |
| 09 - Services | ⬜ Not Started | - | - |
| 10 - Chess-Coach | ⬜ Not Started | - | - |

**Update this table as you complete lessons!**

---

## Estimated Timeline

- **Fast Track**: 1 lesson per day = 10 days
- **Normal Pace**: 2 lessons per week = 5 weeks
- **Deep Dive**: 1 lesson per week = 10 weeks

**Recommendation**: Start with 1-2 lessons, then adjust pace based on comfort

---

## Prerequisites

Before starting:
- ✅ Go basics (you have this)
- ✅ HTTP fundamentals
- ✅ JSON marshaling/unmarshaling
- ✅ Basic goroutines and channels

If rusty on any, we'll review during lessons.

---

## Tools Needed

```bash
# Install dependencies
go get github.com/gorilla/websocket
go get github.com/google/uuid

# Optional but recommended
go install github.com/cosmtrek/air@latest  # Hot reload
```

---

## Success Criteria

By the end, you will:
- ✅ Understand WebSocket protocol deeply
- ✅ Know when to use WebSocket vs REST
- ✅ Implement production-ready WebSocket servers
- ✅ Master the Hub pattern for connection management
- ✅ Handle concurrent connections safely
- ✅ Understand chess-coach architecture completely
- ✅ Be able to extend chess-coach with new features
- ✅ Debug WebSocket issues confidently

---

## Quick Reference

### Key Patterns in Chess-Coach

1. **Hub Pattern**: Central manager for all connections
2. **Client Struct**: Individual connection wrapper
3. **Read/Write Pumps**: Concurrent goroutines for I/O
4. **Message Routing**: Type-based handler dispatch
5. **Room Management**: Game-specific client groups
6. **Service Integration**: WebSocket → Business Logic

### Common WebSocket Operations

```go
// Upgrade connection
conn, _ := upgrader.Upgrade(w, r, nil)

// Read message
_, msg, _ := conn.ReadMessage()

// Write message
conn.WriteMessage(websocket.TextMessage, []byte("hello"))

// Close connection
conn.Close()
```

---

## Next Steps

When ready to begin:
```bash
cd /Users/ankit/code/learn/chess-coach/learnings/fundamentals/websocket
```

Say: **"Let's start Lesson 1"**

I'll create the lesson folder and walk you through it!

---

**Last Updated**: 2025-11-02
**Created By**: SuperClaude Learning System
**For**: Chess Coach WebSocket Implementation
