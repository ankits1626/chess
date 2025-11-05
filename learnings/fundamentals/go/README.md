# Go Fundamentals for Chess Coach

Learn Go progressively through hands-on lessons focused on your chess-coach websocket project.

---

## Current Progress

- [x] **Lesson 1**: Variables & Basic Syntax ✅ READY
- [ ] **Lesson 2**: Pointers (`&` and `*`)
- [ ] **Lesson 3**: Structs & Methods
- [x] **Lesson 6**: Goroutines (`go`) ✅ READY - **DO THIS BEFORE LESSON 4!**
- [x] **Lesson 4**: Channels (`<-`) ✅ READY - Requires Lesson 6 first
- [ ] **Lesson 5**: Select Statement
- [ ] **Lesson 7**: Mutexes (`mu.Lock()`)
- [ ] **Lesson 8**: Full Integration

**Recommended Order**: 1 → 6 → 4 → 2 → 3 → 5 → 7 → 8

---

## Quick Start

### Start Lesson 1 Now

```bash
cd 01-variables-basics
cat START_HERE.md
```

### Your Confusing Syntax

After completing the lessons, you'll understand:
```go
r.mu.Lock()                              // Lesson 7
case message, ok := <-c.send:            // Lessons 1, 4, 5
&Client{hub: hub, conn: conn}            // Lesson 2
<-sigChan                                // Lesson 4
time.NewTicker(pingPeriod)               // Lesson 6
go c.writePump()                         // Lesson 6
func (c *Client) readPump() {}           // Lessons 2, 3
```

---

## Learning Plans

Two plans available:

1. **[PRACTICAL_PLAN.md](PRACTICAL_PLAN.md)** ⭐ RECOMMENDED
   - 8 focused lessons
   - Directly addresses your confusion
   - 3-week timeline
   - Websocket-focused

2. **[LEARNING_PLAN.md](LEARNING_PLAN.md)** (Comprehensive)
   - 25 detailed lessons
   - Full backend mastery
   - 8-10 week timeline
   - All backend concepts

**Start with the Practical Plan!**

---

## Lesson Structure

Each lesson contains:
```
01-variables-basics/
├── START_HERE.md      # Quick start guide
├── README.md          # Concept explanations
├── examples/          # Runnable code (5-7 files)
├── exercises/         # Practice problems
├── references.md      # Real backend examples
└── SUMMARY.md         # Your notes (complete after)
```

---

## How to Use

### For Each Lesson:

1. **Read** - Open `START_HERE.md`
2. **Learn** - Read `README.md` for concepts
3. **Run** - Execute examples: `go run examples/*.go`
4. **Practice** - Complete exercises
5. **Connect** - Find patterns in backend (references.md)
6. **Summarize** - Fill in `SUMMARY.md`
7. **Next** - Move to next lesson

---

## Time Investment

**Per Lesson**:
- Fast: 30 min (read + run)
- Thorough: 60 min (+ practice)
- Deep: 90 min (+ find backend patterns)

**Total Course (8 lessons)**:
- Fast: 4 hours
- Thorough: 8 hours
- Deep: 12 hours

**Recommended**: 1 lesson per day for 8 days (1 hour each)

---

## Prerequisites

- Go installed (`go version`)
- Chess-coach backend cloned
- Basic programming knowledge
- Curiosity!

---

## Support

**Stuck?**
- Re-read the lesson README
- Run examples again
- Modify code to see what breaks
- Check references.md for real examples

**Questions?**
Document in lesson SUMMARY.md and ask!

---

## Ready?

**Start here**: [`01-variables-basics/START_HERE.md`](01-variables-basics/START_HERE.md)

After Lesson 1, you'll understand:
- `conn, err := upgrader.Upgrade(w, r, nil)`
- `messageType, message, err := conn.ReadMessage()`
- `message, ok := <-c.send`

Let's go! 🚀
