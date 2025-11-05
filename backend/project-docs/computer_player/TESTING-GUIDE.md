# Testing Guide: Human vs Computer Chess

**Purpose**: Step-by-step instructions for testing the computer player functionality via WebSocket

**Last Updated**: 2025-11-05

---

## 🧪 Prerequisites

### 1. Install wscat

```bash
npm install -g wscat
```

### 2. Ensure Server is Running

```bash
# Start the server with Docker
docker compose up api

# Or run locally (requires Stockfish installed)
go run cmd/server/main.go
```

### 3. Verify Stockfish is Available

```bash
# In Docker
docker compose exec api which stockfish
# Should output: /usr/local/bin/stockfish

docker compose exec api stockfish
# Should show: Stockfish 17.1 by the Stockfish developers
```

---

## 🎮 Testing Human vs Computer

### Step 1: Connect to WebSocket

**Important**: Use a valid UUID for the user_id parameter!

```bash
wscat -c "ws://localhost:8080/ws?user_id=25d30da5-0cc4-4f5a-8c88-69d0f90b004c"
```

**Expected Output:**
```
Connected (press CTRL+C to quit)
>
```

**Common Error**: If you see `"invalid user ID"`, make sure:
- The user exists in the database, OR
- Use a UUID from an existing user in your database

---

### Step 2: Create a Game Against Computer

**Send this JSON message:**

```json
{"id":"req-1","type":"request","action":"createGame","data":{"mode":"human_vs_computer","difficulty":"easy","timeControl":"5+0"}}
```

**Expected Response:**

```json
{
  "id": "req-1",
  "type": "response",
  "data": {
    "gameId": "48e62cbf-a9db-4a8f-b918-2748cc9722b4",
    "status": "active",
    "side": "white",
    "mode": "human_vs_computer",
    "difficulty": "easy",
    "timeControl": "5+0",
    "fen": "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
  },
  "success": true
}
```

**Key Points:**
- ✅ Game is immediately `"active"` (no waiting for second player)
- ✅ You are `"white"` (playing first)
- ✅ Copy the `gameId` - you'll need it for making moves!

---

### Step 3: Make Your First Move

**Replace `<GAME_ID>` with your actual gameId from step 2:**

```json
{"id":"req-2","type":"request","action":"makeMove","data":{"gameId":"<GAME_ID>","move":"e2e4"}}
```

**Example with actual gameId:**

```json
{"id":"req-2","type":"request","action":"makeMove","data":{"gameId":"48e62cbf-a9db-4a8f-b918-2748cc9722b4","move":"e2e4"}}
```

**Expected Response (Your Move Confirmed):**

```json
{
  "id": "req-2",
  "type": "response",
  "data": {
    "moveNumber": 1,
    "san": "e4",
    "uci": "e2e4",
    "fen": "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
    "gameOver": false
  },
  "success": true
}
```

---

### Step 4: Receive Computer's Automatic Response

**Within 1 second, you'll receive:**

```json
{
  "type": "event",
  "event": "move",
  "data": {
    "gameId": "48e62cbf-a9db-4a8f-b918-2748cc9722b4",
    "moveSAN": "e5",
    "moveUCI": "e7e5",
    "fen": "rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq e6 0 2",
    "moveNumber": 2,
    "isGameOver": false,
    "result": null
  }
}
```

**🎉 The computer just played e7e5!**

---

### Step 5: Continue Playing

Keep making moves - the computer will respond to each one!

**Example moves:**

```json
{"id":"req-3","type":"request","action":"makeMove","data":{"gameId":"<GAME_ID>","move":"Nf3"}}
```

```json
{"id":"req-4","type":"request","action":"makeMove","data":{"gameId":"<GAME_ID>","move":"Bc4"}}
```

```json
{"id":"req-5","type":"request","action":"makeMove","data":{"gameId":"<GAME_ID>","move":"Nc3"}}
```

**The computer will automatically respond after each of your moves!**

---

## 🎯 Valid Move Formats

You can use either:

### UCI Format (Universal Chess Interface)
- Pawn moves: `e2e4`, `d2d4`, `c7c5`
- Knight moves: `g1f3`, `b1c3`
- Bishop moves: `f1c4`, `f8b4`
- Castling: `e1g1` (kingside), `e1c1` (queenside)
- Promotion: `e7e8q` (promote to queen)

### SAN Format (Standard Algebraic Notation)
- Pawn moves: `e4`, `d4`, `c5`
- Knight moves: `Nf3`, `Nc3`
- Bishop moves: `Bc4`, `Bb4`
- Castling: `O-O` (kingside), `O-O-O` (queenside)
- Captures: `Nxe5`, `dxe5`
- Check: `Qh5+`
- Checkmate: `Qf7#`

---

## 🎮 Testing Different Difficulty Levels

### Easy (100ms think time)

```json
{"id":"1","type":"request","action":"createGame","data":{"mode":"human_vs_computer","difficulty":"easy","timeControl":"5+0"}}
```

**Characteristics:**
- Quick responses
- Good for beginners
- May make some weak moves

### Medium (500ms think time)

```json
{"id":"1","type":"request","action":"createGame","data":{"mode":"human_vs_computer","difficulty":"medium","timeControl":"5+0"}}
```

**Characteristics:**
- Balanced gameplay
- Reasonable challenge
- Good for intermediate players

### Hard (2 seconds think time)

```json
{"id":"1","type":"request","action":"createGame","data":{"mode":"human_vs_computer","difficulty":"hard","timeControl":"5+0"}}
```

**Characteristics:**
- Strong play
- Longer responses (up to 2 seconds)
- Challenging for advanced players

---

## 🎯 Testing Other Game Modes

### Human vs Human

```json
{"id":"1","type":"request","action":"createGame","data":{"mode":"human_vs_human","timeControl":"5+0"}}
```

**Expected:**
- Status: `"waiting"` (needs second player to join)
- Another player must join via `joinGame` action

**Note**: Computer modes are currently the main focus. Human vs Human is the legacy mode.

---

## 🐛 Troubleshooting

### Issue: "Connection refused"

**Solution:**
```bash
# Check if server is running
docker compose ps

# Check server logs
docker compose logs -f api
```

### Issue: "invalid user ID"

**Cause**: User UUID doesn't exist in database or wrong format

**Solution 1** - Use an existing user:
```bash
# Get existing users from database
docker compose exec db psql -U postgres -d chess_coach -c "SELECT id FROM users LIMIT 5;"

# Use one of the returned UUIDs
wscat -c "ws://localhost:8080/ws?user_id=<EXISTING_UUID>"
```

**Solution 2** - Create a test user:
```bash
docker compose exec db psql -U postgres -d chess_coach
```

```sql
INSERT INTO users (id, email, username, password_hash, created_at, updated_at)
VALUES (
  '25d30da5-0cc4-4f5a-8c88-69d0f90b004c',
  'test@example.com',
  'testplayer',
  'dummy_hash',
  NOW(),
  NOW()
)
ON CONFLICT (id) DO NOTHING;
```

### Issue: "invalid game mode"

**Cause**: Wrong mode string format

**Solution**: Use underscores, not camelCase:
- ✅ `"human_vs_computer"`
- ❌ `"humanVsComputer"`

**Valid modes:**
- `"human_vs_computer"`
- `"human_vs_human"`
- `"computer_vs_computer"` (not fully implemented)
- `"human_vs_llm"` (not implemented)

### Issue: "Invalid move"

**Cause**: Move is illegal in current position

**Solution:**
- Check the current FEN position
- Ensure it's your turn (white moves on even move numbers)
- Verify the move is legal in chess

### Issue: No computer response

**Check server logs:**
```bash
docker compose logs -f api | grep Computer
```

**Expected log output:**
```
Computer (black) in game xxx plays: e7e5
```

**If missing:**
- Check if Stockfish is installed: `docker compose exec api which stockfish`
- Check for errors in logs: `docker compose logs -f api`

### Issue: Stockfish not found

**Verify installation:**
```bash
docker compose exec api which stockfish
# Expected: /usr/local/bin/stockfish

docker compose exec api stockfish
# Expected: Stockfish 17.1 by the Stockfish developers
```

**If missing, rebuild container:**
```bash
docker compose build api
docker compose up api
```

---

## 📊 Complete Test Session Example

Here's a complete session playing a few moves:

```bash
# 1. Connect
wscat -c "ws://localhost:8080/ws?user_id=25d30da5-0cc4-4f5a-8c88-69d0f90b004c"

# 2. Create game
> {"id":"1","type":"request","action":"createGame","data":{"mode":"human_vs_computer","difficulty":"medium","timeControl":"5+0"}}
< [Game created response with gameId]

# 3. Play e4 (King's Pawn Opening)
> {"id":"2","type":"request","action":"makeMove","data":{"gameId":"<GAME_ID>","move":"e2e4"}}
< [Your move confirmed]
< [Computer plays e7e5 automatically]

# 4. Play Nf3 (Knight develops)
> {"id":"3","type":"request","action":"makeMove","data":{"gameId":"<GAME_ID>","move":"Nf3"}}
< [Your move confirmed]
< [Computer responds - maybe Nc6]

# 5. Play Bc4 (Italian Opening)
> {"id":"4","type":"request","action":"makeMove","data":{"gameId":"<GAME_ID>","move":"Bc4"}}
< [Your move confirmed]
< [Computer responds - maybe Bc5]

# 6. Castle kingside
> {"id":"5","type":"request","action":"makeMove","data":{"gameId":"<GAME_ID>","move":"O-O"}}
< [Your move confirmed]
< [Computer responds]

# Continue playing until checkmate or resignation!
```

---

## ✅ Success Criteria

Your test is successful when you see:

1. ✅ Game creates with `status: "active"`
2. ✅ Your move returns success response
3. ✅ Computer automatically sends move event
4. ✅ Computer's move is legal and makes sense
5. ✅ You can continue making moves
6. ✅ Computer responds to each move

---

## 🎓 What to Test

### Basic Functionality
- [x] Create game with easy difficulty
- [x] Create game with medium difficulty
- [x] Create game with hard difficulty
- [x] Make valid moves
- [x] Computer responds automatically
- [x] Continue game for multiple moves

### Edge Cases
- [ ] Invalid move attempt
- [ ] Game to checkmate
- [ ] Game to stalemate
- [ ] Resignation
- [ ] Disconnection handling

### Performance
- [ ] Computer response time (should be < difficulty time)
- [ ] Easy: ~100ms
- [ ] Medium: ~500ms
- [ ] Hard: ~2000ms

---

## 📝 Next Steps

Once basic testing works:

1. **Test full game** - Play to checkmate/draw
2. **Test resignation** - `{"id":"X","type":"request","action":"resign","data":{"gameId":"<ID>"}}`
3. **Test different openings** - See how computer responds to various strategies
4. **Test all difficulty levels** - Compare computer strength

---

## 🎉 Expected Outcome

After following this guide, you should be able to:

✅ Create games against computer
✅ Make moves and see computer responses
✅ Play complete games
✅ Test different difficulty levels
✅ Verify all functionality works

**Congratulations on having a working chess AI opponent!** 🏆♟️

---

**Questions or Issues?** Check the server logs for detailed error messages:
```bash
docker compose logs -f api
```
