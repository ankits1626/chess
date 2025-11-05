# Lesson 1 Summary

**Complete this after finishing the lesson!**

## What I Learned

### Key Concepts
- [ ] Difference between `var` and `:=`
- [ ] When to use each declaration style
- [ ] Multiple return values
- [ ] Zero values
- [ ] Error handling pattern

### Confusing Syntax - Now Clear!

**Before Lesson**: `message, ok := <-c.send` ← What??

**After Lesson**:
- `:=` creates two variables: `message` and `ok`
- `<-c.send` receives from a channel (we'll learn this in Lesson 4)
- `message` gets the value from the channel
- `ok` is true if channel is open, false if closed

## My Notes

Write your observations here:

1. What was most confusing?
2. What clicked for you?
3. Any questions remaining?

## Examples I Ran

- [ ] 01-basic.go
- [ ] 02-short-declaration.go
- [ ] 03-multiple.go
- [ ] 04-zero-values.go
- [ ] 05-websocket-pattern.go

## Exercises Completed

- [ ] Exercise 1: Fix errors
- [ ] Exercise 2: Multiple returns
- [ ] Exercise 3: Error handling
- [ ] Exercise 4: Zero values
- [ ] Exercise 5: Websocket pattern
- [ ] Bonus: Custom function

## Real Code I Now Understand

Find examples in the chess-coach backend where you now understand the variable declarations:

1. File: `_______________________`
   Line: `_____`
   Code:
   ```go

   ```

2. File: `_______________________`
   Line: `_____`
   Code:
   ```go

   ```

## Questions for Review

Write any questions you still have:

1.
2.
3.

---

**Ready for Lesson 2?** Check this when ready: [ ]

Lesson 2 will cover: **Pointers (`&` and `*`)** - Understanding `&Client`, `*Hub`
