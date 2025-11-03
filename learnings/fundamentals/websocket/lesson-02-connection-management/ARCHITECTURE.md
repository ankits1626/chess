# Lesson 2: Architecture & Flow Diagrams

Visual guide to understanding connection management architecture.

---

## System Overview

```
┌─────────────────────────────────────────────────────-────────┐
│                         Browser Client                       │
│                                                              │
│  ┌────────────────────────────────────────────────────┐      │
│  │       ReconnectingWebSocket Class                  │      │
│  │                                                    │      │
│  │  • Auto-reconnection with exponential backoff      │      │
│  │  • Message queuing during disconnection            │      │
│  │  • Connection health monitoring                    │      │
│  └────────────────────────────────────────────────────┘      │
│                          │                                   │
└──────────────────────────┼───────────────────────────────────┘
                           │
                           │ WebSocket (ws://)
                           │
┌──────────────────────────┼───────────────────────────────────┐
│                          ▼                                   │
│                  HTTP/WebSocket Server                       │
│                                                              │
│  ┌────────────────────────────────────────────────────┐      │
│  │            Connection Registry                     │      │
│  │                                                    │      │
│  │  register chan ──┐                                 │      │
│  │  unregister chan ├─> Run() event loop              │      │
│  │  broadcast chan ─┘                                 │      │
│  │                                                    │      │
│  │  Tracks all active clients                         │      │
│  │  Handles graceful shutdown                         │      │
│  └────────────────────────────────────────────────────┘      │
│                          │                                   │
│                          │ manages                           │
│                          ▼                                   │
│  ┌────────────────────────────────────────────────────┐      │
│  │              Client (per connection)               │      │
│  │                                                    │      │
│  │  ┌──────────────┐         ┌───────────────┐        │      │
│  │  │  readPump()  │         │  writePump()  │        │      │
│  │  │  goroutine   │         │   goroutine   │        │      │
│  │  │              │         │               │        │      │
│  │  │ • Read msgs  │         │ • Write msgs  │        │      │
│  │  │ • Handle     │◄───────►│ • Send pings  │        │      │
│  │  │   pongs      │  sync   │ • Monitor     │        │      │
│  │  │ • Set        │         │   timeouts    │        │      │
│  │  │   deadlines  │         │               │        │      │
│  │  └──────────────┘         └───────────────┘        │      │
│  │         │                         ▲                │      │
│  │         │                         │                │      │
│  │         └────► send channel ──────┘                │      │
│  │                (buffered: 256)                     │      │
│  └────────────────────────────────────────────────────┘      │
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

---

## Connection Lifecycle

```
┌─────────────────────────────────────────────────────────────┐
│                    Connection Lifecycle                     │
└─────────────────────────────────────────────────────────────┘

1. CONNECTION ESTABLISHMENT
   ═══════════════════════

   Client                          Server
     │                               │
     │  HTTP GET /ws                 │
     │  Upgrade: websocket ────────> │
     │                               │
     │                               ├─ upgrader.Upgrade()
     │                               ├─ Create Client struct
     │                               ├─ client.send = make(chan []byte, 256)
     │                               ├─ registry.register <- client
     │                               │
     │                               ├─ go client.writePump()
     │                               │    └─ Start ping ticker (54s)
     │                               │
     │                               └─ go client.readPump()
     │                                    ├─ SetReadDeadline(now + 60s)
     │                                    └─ SetPongHandler(reset deadline)
     │                               │
     │  <─────────────────────────   │
     │  101 Switching Protocols      │
     │                               │
     │  <─────────────────────────   │
     │  Welcome message              │
     │                               │
   [Connected]                   [Connected]


2. HEARTBEAT (PING/PONG)
   ══════════════════════

   Time: 0s
   ┌────────────────────────────────────────────────────────┐
   │ Connection established                                 │
   │ ReadDeadline set: now + 60s = 60s                      │
   └────────────────────────────────────────────────────────┘

   Time: 54s (pingPeriod)
   ┌────────────────────────────────────────────────────-────┐
   │ writePump: Ticker fires                                 │
   │                                                         │
   │ Client                          Server                  │
   │   │                               │                     │
   │   │                               ├─ ticker.C received  │
   │   │                               ├─ SetWriteDeadline   │
   │   │  <──────── PING ─────────────┤   (now + 10s)        │
   │   │                               │                     │
   │   ├─ Auto-respond with PONG       │                     │
   │   │                               │                     │
   │   │  ─────── PONG ──────────────> │                     │
   │   │                               │                     │
   │   │                               ├─ readPump receives  │
   │   │                               ├─ PongHandler called │
   │   │                               └─ SetReadDeadline    │
   │   │                                  (now + 60s = 114s) │
   └────────────────────────────────────────────────────────┘

   Time: 108s (next ping)
   ┌────────────────────────────────────────────────────────┐
   │ Repeat ping/pong cycle                                 │
   │ ReadDeadline updated: now + 60s = 168s                 │
   └────────────────────────────────────────────────────────┘


3. MESSAGE EXCHANGE
   ═════════════════

   Client                          Server
     │                               │
     │  "Hello" ──────────────────>  │
     │                               │
     │                               ├─ readPump: ReadMessage()
     │                               ├─ Receive: "Hello"
     │                               ├─ Echo: "Server echo: Hello"
     │                               └─ client.send <- response
     │                               │
     │                               ├─ writePump: <-client.send
     │  <─────── "Server echo" ─────┤   WriteMessage()
     │                               │


4. ABNORMAL DISCONNECTION (Network Failure)
   ═════════════════════════════════════════

   Time: 30s
   ┌───────────────────────────────────────────────────-─────┐
   │ Client network fails (WiFi disconnect, etc.)            │
   │                                                         │
   │ Client                          Server                  │
   │   💥                              │                     │
   │  [Dead]                           │                     │
   │                                   │                     │
   │                              [Still thinks alive]       │
   │                              [Waiting to send ping...]  │
   └───────────────────────────────────────────────────────-─┘

   Time: 54s
   ┌─────────────────────────────────────────────────────-───┐
   │ Server tries to send PING                               │
   │                                                         │
   │   💥                              │                     │
   │  [Dead]                           │                     │
   │                                   │                     │
   │   (no response)     PING ────X   ├─ Sent PING           │
   │                                   │                     │
   │                                   └─ Waiting for PONG...│
   └───────────────────────────────────────────────────────-─┘

   Time: 60s (ReadDeadline expires)
   ┌───────────────────────────────────────────────────-─────┐
   │ Server detects timeout                                  │
   │                                                         │
   │   💥                              │                     │
   │  [Dead]                           │                     │
   │                                   │                     │
   │                                   ├─ ⏰ ReadDeadline!   │
   │                                   ├─ ReadMessage() err  │
   │                                   ├─ Close connection   │
   │                                   ├─ registry.unregister│
   │                                   └─ Clean up resources │
   └───────────────────────────────────────────────────────-─┘


5. GRACEFUL DISCONNECTION (User-Initiated)
   ════════════════════════════════════════

   Client                          Server
     │                               │
     │  User clicks "Disconnect"     │
     │                               │
     │  Close(1000, "User closed") > │
     │                               │
     │                               ├─ ReadMessage() receives
     │                               │   CloseMessage
     │                               ├─ Normal closure detected
     │                               ├─ registry.unregister
     │  <─────── Close ACK ──────────┤   conn.Close()
     │                               │
   [Closed]                      [Closed]


6. GRACEFUL SHUTDOWN (Server-Initiated)
   ════════════════════════════════════

   Server receives SIGINT/SIGTERM
     │
     ├─ registry.Shutdown() called
     │
     ├─ For each client:
     │    ├─ FormatCloseMessage(1001, "Server shutting down")
     │    ├─ WriteControl(CloseMessage)
     │    └─ conn.Close()
     │
     └─ server.Shutdown(ctx)

   Client                          Server
     │                               │
     │  <──── CloseMessage ──────────┤ Code: 1001
     │     (Server shutting down)    │ Reason: "Server shutting..."
     │                               │
     ├─ onclose event fires          │
     ├─ event.code === 1001          │
     ├─ "Server going away"          │
     └─ scheduleReconnect()          │
          (with appropriate delay)   │
```

---

## Goroutine Communication

```
┌──────────────────────────────────────────────────────────────┐
│                    Per-Client Goroutines                     │
└──────────────────────────────────────────────────────────────┘

                    ┌─────────────────┐
                    │   WebSocket     │
                    │   Connection    │
                    └────────┬────────┘
                             │
                    ┌────────┴──────-──┐
                    │                  │
              ┌─────▼─────┐      ┌────▼──────┐
              │ readPump  │      │writePump  │
              │ goroutine │      │ goroutine │
              └───────────┘      └───────────┘
                    │                  ▲
                    │                  │
                    │            ┌─────┴─────┐
                    │            │   send    │
                    └───────────>│  channel  │
                      writes to  │(buf: 256) │
                                 └───────────┘
                                       ▲
                                       │
                                 ┌─────┴──────┐
                                 │   Other     │
                                 │ goroutines  │
                                 │   (e.g.,    │
                                 │ broadcast)  │
                                 └────────────┘


Details:

readPump:
  • Reads from WebSocket connection
  • Processes incoming messages
  • Handles pongs (resets deadline)
  • Writes responses to send channel
  • Runs until error or close

writePump:
  • Reads from send channel
  • Writes to WebSocket connection
  • Sends periodic pings (ticker)
  • Monitors write deadlines
  • Runs until error or close

Communication:
  • readPump → send channel → writePump
  • Thread-safe through Go channels
  • No mutexes needed for message passing
  • Buffered channel (256) for flow control
```

---

## Reconnection State Machine

```
┌──────────────────────────────────────────────────────────────┐
│              Client Reconnection State Machine               │
└──────────────────────────────────────────────────────────────┘

                    ┌──────────────┐
                    │ DISCONNECTED │◄────────-──┐
                    └──────┬───────┘            │
                           │                    │
                           │ connect()          │
                           │                    │
                           ▼                    │
                    ┌──────────────┐            │
              ┌────►│  CONNECTING  │            │
              │     └──────┬───────┘            │
              │            │                    │
              │            │                    │
    retry     │      ┌─────┴─────┐              │
              │      │           │              │
              │      │ Success   │ Failed       │
              │      ▼           ▼              │
              │   ┌─────────┐ ┌──────────┐      │
              │   │CONNECTED│ │  ERROR   │      │
              │   └────┬────┘ └────┬─────┘      │
              │        │           │            │
              │        │           │            │
              │        │           ▼            │
              │        │    ┌──────────────┐    │
              │        │    │RECONNECTING  │    │
              │        │    │              │    │
              │        │    │ Exponential  │    │
              │        │    │  Backoff:    │    │
              │        │    │ 1s→2s→4s→8s  │─---┘
              │        │    └──────────────┘
              │        │
              │        │ onclose
              │        ▼
              └─── (check code)
                        │
                 ┌──────┴──────┐
                 │             │
            1000 │             │ 1006
          (normal)            (abnormal)
                 │             │
                 ▼             ▼
         Don't reconnect   Reconnect


State Transitions:

DISCONNECTED → CONNECTING
  • User clicks "Connect"
  • Or auto-reconnect triggered

CONNECTING → CONNECTED
  • WebSocket handshake succeeds
  • onopen event fires
  • Flush message queue
  • Reset retry counter

CONNECTING → ERROR
  • Connection refused
  • Network unreachable
  • Handshake failed

ERROR → RECONNECTING
  • If auto-reconnect enabled
  • If attempts < max attempts
  • Calculate exponential backoff delay

RECONNECTING → CONNECTING
  • After backoff delay expires
  • Increment retry counter
  • Double backoff delay (cap at 30s)

CONNECTED → DISCONNECTED
  • Code 1000 (normal closure)
  • User requested disconnect
  • Don't auto-reconnect

CONNECTED → RECONNECTING
  • Code 1001 (server going away)
  • Code 1006 (abnormal closure)
  • Auto-reconnect enabled
```

---

## Message Flow

```
┌──────────────────────────────────────────────────────────────┐
│                   Message Send Flow                          │
└──────────────────────────────────────────────────────────────┘

User types message
      │
      ▼
┌─────────────┐
│  sendMessage│
│  function   │
└──────┬──────┘
       │
       │ Check connection state
       │
   ┌───┴────┐
   │        │
   │    OPEN?
   │        │
   └┬──────┬┘
    │      │
  Yes│     │No
    │      │
    ▼      ▼
┌───────┐ ┌──────────┐
│ Send  │ │  Queue   │
│ Now   │ │ Message  │
└───┬───┘ └────┬─────┘
    │          │
    │          ├─ messageQueue.push()
    │          ├─ Show "Queued" indicator
    │          └─ Wait for reconnection
    │
    ├─ ws.send(message)
    ├─ Add to sent messages
    └─ Update statistics


┌──────────────────────────────────────────────────────────────┐
│                  Message Receive Flow                        │
└──────────────────────────────────────────────────────────────┘

Server                           Client
  │                               │
  ├─ Send message                 │
  │                               │
  │  ───────────────────────────> │
  │                               │
  │                               ├─ onmessage event
  │                               │
  │                               ├─ Parse message
  │                               │   (Text? JSON? Binary?)
  │                               │
  │                               ├─ Update statistics
  │                               │   (receivedCount++)
  │                               │
  │                               └─ Display in UI
  │                                   └─ addMessage()


┌──────────────────────────────────────────────────────────────┐
│                  Message Queue Flush                         │
└──────────────────────────────────────────────────────────────┘

Reconnection succeeds
      │
      ▼
flushMessageQueue()
      │
      ├─ messageQueue.length > 0?
      │
      ├─ Yes: Loop through queue
      │   │
      │   ├─ For each message:
      │   │   ├─ ws.send(message)
      │   │   ├─ messageQueue.shift()
      │   │   └─ Update UI
      │   │
      │   └─ Clear queue
      │       └─ Hide "Queued" indicator
      │
      └─ No: Continue normal operation
```

---

## Timeout Management

```
┌──────────────────────────────────────────────────────────────┐
│                  Read Deadline Management                    │
└──────────────────────────────────────────────────────────────┘

Time: 0s
┌────────────────────────────────────────────────────┐
│ Connection established                              │
│                                                      │
│ conn.SetReadDeadline(time.Now().Add(pongWait))     │
│                                                      │
│ Deadline: 0s + 60s = 60s                           │
│                                                      │
│ ┌────────────────────────────────────────────┐    │
│ │ Deadline Timer                              │    │
│ │ [●●●●●●●●●●●●●●●●●●●●●●●●●●●●●] 60s       │    │
│ └────────────────────────────────────────────┘    │
└────────────────────────────────────────────────────┘

Time: 54s (Ping sent)
┌────────────────────────────────────────────────────┐
│ Server sends PING                                   │
│                                                      │
│ Deadline still: 60s                                 │
│                                                      │
│ ┌────────────────────────────────────────────┐    │
│ │ Deadline Timer                              │    │
│ │ [●●●●●●······················] 6s left     │    │
│ └────────────────────────────────────────────┘    │
└────────────────────────────────────────────────────┘

Time: 55s (Pong received)
┌────────────────────────────────────────────────────┐
│ Client sends PONG                                   │
│                                                      │
│ PongHandler called:                                 │
│ conn.SetReadDeadline(time.Now().Add(pongWait))     │
│                                                      │
│ New deadline: 55s + 60s = 115s                     │
│                                                      │
│ ┌────────────────────────────────────────────┐    │
│ │ Deadline Timer (RESET!)                     │    │
│ │ [●●●●●●●●●●●●●●●●●●●●●●●●●●●●●] 60s       │    │
│ └────────────────────────────────────────────┘    │
└────────────────────────────────────────────────────┘


What if no pong?
═══════════════════

Time: 60s (Deadline expires, no pong)
┌────────────────────────────────────────────────────┐
│ Deadline reached!                                   │
│                                                      │
│ ┌────────────────────────────────────────────┐    │
│ │ Deadline Timer                              │    │
│ │ [······························] EXPIRED!  │    │
│ └────────────────────────────────────────────┘    │
│                                                      │
│ conn.ReadMessage() returns error                   │
│   └─> "i/o timeout" or similar                     │
│                                                      │
│ readPump detects error                             │
│   ├─ Break read loop                               │
│   ├─ defer: registry.unregister                    │
│   └─ defer: conn.Close()                           │
└────────────────────────────────────────────────────┘


┌──────────────────────────────────────────────────────────────┐
│                  Write Deadline Management                    │
└──────────────────────────────────────────────────────────────┘

Each write operation:

conn.SetWriteDeadline(time.Now().Add(writeWait))  // 10s
conn.WriteMessage(messageType, data)

If write doesn't complete in 10s:
  └─> Error returned
      ├─ writePump detects error
      ├─ Break write loop
      └─ Connection closed

Why 10 seconds?
  • Slow clients (mobile, poor network)
  • Large messages
  • Network congestion
  • Balance: too short = false positives
            too long = resource waste
```

---

## Connection Registry Event Loop

```
┌──────────────────────────────────────────────────────────────┐
│              Connection Registry Architecture                 │
└──────────────────────────────────────────────────────────────┘

┌────────────────────────────────────────────────────┐
│  ConnectionRegistry                                 │
│                                                      │
│  • clients:    map[*Client]bool                    │
│  • register:   chan *Client                        │
│  • unregister: chan *Client                        │
│  • broadcast:  chan []byte                         │
│  • mu:         sync.RWMutex                        │
└─────────────┬──────────────────────────────────────┘
              │
              │ go Run()
              ▼
     ┌────────────────┐
     │  Event Loop    │
     │  (goroutine)   │
     └────────┬───────┘
              │
              │ select { ... }
              │
   ┌──────────┼──────────┬─────────────┐
   │          │          │             │
   ▼          ▼          ▼             ▼
┌────────┐ ┌──────┐ ┌─────────┐ ┌──────────┐
│register│ │unregis│ │broadcast│ │ (more)  │
│ case   │ │ter case│ │  case   │ │         │
└────┬───┘ └───┬───┘ └────┬────┘ └──────────┘
     │         │          │
     │         │          │
     ▼         ▼          ▼

Register:
  ├─ Lock mutex
  ├─ clients[client] = true
  ├─ Unlock mutex
  └─ Log connection count

Unregister:
  ├─ Lock mutex
  ├─ if client exists:
  │    ├─ delete from map
  │    └─ close(client.send)
  ├─ Unlock mutex
  └─ Log connection count

Broadcast:
  ├─ RLock mutex (read lock)
  ├─ For each client:
  │    ├─ Try send to client.send
  │    ├─ If buffer full:
  │    │    └─ Close that client
  │    └─ Continue to next
  └─ RUnlock mutex


Why channels + event loop?
═══════════════════════════

Instead of:
  ❌ Multiple goroutines locking mutex frequently
  ❌ Lock contention
  ❌ Potential deadlocks

We have:
  ✅ Single goroutine manages map
  ✅ Other goroutines send to channels
  ✅ No lock contention
  ✅ Go's channel synchronization
```

---

## Summary

Key architectural patterns:

1. **Two Goroutines Per Connection**
   - readPump: Blocking reads + deadline management
   - writePump: Ticker for pings + send channel for messages

2. **Channel-Based Communication**
   - send channel connects readPump → writePump
   - No mutex needed for message passing

3. **Connection Registry**
   - Central event loop manages all connections
   - Channels for register/unregister/broadcast
   - Thread-safe with minimal locking

4. **Deadline Management**
   - Read deadline: Reset on each pong
   - Write deadline: Set before each write
   - Automatic timeout detection

5. **Client Reconnection**
   - State machine (Disconnected → Connecting → Connected)
   - Exponential backoff prevents server overload
   - Message queuing prevents data loss

These patterns enable:
- ✅ Scalability (1000s of concurrent connections)
- ✅ Reliability (automatic recovery)
- ✅ Resource protection (timeouts)
- ✅ Clean shutdown (graceful close)

---

**Ready to see it in action?** Check [QUICKSTART.md](QUICKSTART.md)!
