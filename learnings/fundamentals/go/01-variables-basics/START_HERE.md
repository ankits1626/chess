# Lesson 1: Variables & Basic Syntax - START HERE

## What You'll Learn (30 minutes)

By the end of this lesson, confusing syntax like this will make sense:
```go
message, ok := <-c.send  // You'll understand this!
conn, err := upgrader.Upgrade(w, r, nil)  // And this!
messageType, message, err := conn.ReadMessage()  // And this!
```

---

## Learning Path

### Step 1: Read the Concepts (10 min)
Open and read: [README.md](README.md)

Focus on:
- The two ways to declare variables (`var` vs `:=`)
- When to use each
- Multiple return values
- Zero values

### Step 2: Run Examples (15 min)

```bash
cd examples

# Run each example in order
go run 01-basic.go
go run 02-short-declaration.go
go run 03-multiple.go
go run 04-zero-values.go
go run 05-websocket-pattern.go  # Most important for websockets!
```

**Tip**: After running each, modify the code and run again to see what happens!

### Step 3: Practice (30 min)

```bash
cd ../exercises
go run practice.go
```

The exercises have TODO comments - complete them!

All exercises already have solutions, but try to understand WHY each works.

### Step 4: Connect to Real Code (15 min)

Open: [references.md](references.md)

Find the patterns in your actual chess-coach backend code!

This reinforces learning by seeing real-world usage.

### Step 5: Create Your Summary (10 min)

Open: [SUMMARY.md](SUMMARY.md)

Fill in:
- What you learned
- What clicked for you
- Examples from the backend you now understand
- Any remaining questions

---

## Quick Reference

| You'll See | It Means | Example |
|------------|----------|---------|
| `x := 5` | Create variable `x` with value `5` | `conn := getConnection()` |
| `var x int` | Create variable `x` of type `int` (value is 0) | `var count int` |
| `x, y := 1, 2` | Create TWO variables | `conn, err := upgrade()` |
| `a, b := fn()` | Function returns 2 values | `user, err := getUser()` |

---

## Expected Outcomes

After completing:
- ✅ Understand `var` vs `:=`
- ✅ Know why `conn, err := ...` creates two variables
- ✅ Recognize multiple return value pattern
- ✅ Understand error handling: `if err != nil`
- ✅ See zero values in action

---

## Time Estimate

- **Fast track**: 30 minutes (read + run examples)
- **Thorough**: 60 minutes (read + run + practice + references)
- **Deep dive**: 90 minutes (everything + modify examples + find more patterns)

---

## Need Help?

**Stuck on a concept?**
- Re-read that section in README.md
- Run the corresponding example again
- Try modifying the code to see what breaks

**Example not working?**
- Make sure you're in the right directory
- Check you have Go installed: `go version`
- All examples are tested and working!

**Ready for more?**
After completing this lesson and creating your SUMMARY.md, you're ready for:
**Lesson 2: Pointers** - Understanding `&Client`, `*Hub`, and pointer receivers

---

## Let's Begin!

Start with: [README.md](README.md)

Good luck! 🚀
