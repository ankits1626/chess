# Debugging Guide: Computer Player

**Date**: 2025-11-06

This guide helps you debug issues with the computer player functionality.

---

## 🔍 Console Logs to Watch

With the enhanced logging, you should see these messages in the browser console:

### 1. Starting a Game:

```
[GameStore] Starting computer game: {color: 'white', difficulty: 'medium', userId: '25d30da5-0cc4-4f5a-8c88-69d0f90b004c'}
[GameStore] Connecting to WebSocket...
[GameService] Connected to WebSocket
[GameStore] WebSocket connected
[GameStore] Creating game with mode: human_vs_computer, difficulty: medium
[GameService] Sending: {id: 'req-1', type: 'request', action: 'createGame', data: {...}}
[GameService] Received: {id: 'req-1', type: 'response', success: true, data: {...}}
[GameStore] Game created: {gameId: '...', status: 'active', ...}
[GameStore] Setting up event listeners
[GameStore] Computer game started successfully
```

### 2. Making a Move:

```
[GameStore] Move executed: {from: 'e2', to: 'e4', san: 'e4', uci: 'e2e4', opponentType: 'computer', gameId: '...'}
[GameStore] Sending move to backend: e4
[GameService] Sending: {id: 'req-2', type: 'request', action: 'makeMove', data: {gameId: '...', move: 'e4'}}
[GameStore] Move sent successfully, waiting for computer...
[GameService] Received: {id: 'req-2', type: 'response', success: true, data: {...}}
[GameService] Received: {type: 'event', event: 'move', data: {moveUCI: 'e7e5', ...}}
[GameStore] Move event received: {gameId: '...', moveSAN: 'e5', moveUCI: 'e7e5', ...}
[GameStore] Handling computer move: e7e5
[GameStore] Computer move parsed: {san: 'e5', from: 'e7', to: 'e5', ...}
[GameStore] Computer move applied successfully
```

---

## ❌ Common Issues & Solutions

### Issue 1: No Logs Appear

**Symptom**: Console is empty

**Solution**:
1. Open browser DevTools (F12)
2. Go to Console tab
3. Make sure no filters are active
4. Refresh the page

---

### Issue 2: WebSocket Connection Fails

**Symptom**:
```
[GameStore] Connecting to WebSocket...
Error: WebSocket connection failed
```

**Causes**:
1. Backend not running
2. Wrong WebSocket URL
3. CORS issues

**Solution**:
```bash
# Check backend is running
docker compose ps

# Check logs
docker compose logs -f api

# Verify WebSocket endpoint
# Should be: ws://localhost:8080/ws
```

---

### Issue 3: Invalid User ID

**Symptom**:
```
[GameService] Received: {id: 'req-1', type: 'response', error: 'invalid user ID'}
```

**Solution**:
```bash
# Create test user in database
docker compose exec db psql -U postgres -d chess_coach

INSERT INTO users (id, email, username, password_hash, created_at, updated_at)
VALUES (
  '25d30da5-0cc4-4f5a-8c88-69d0f90b004c',
  'test@example.com',
  'testplayer',
  'dummy_hash',
  NOW(),
  NOW()
);
```

---

### Issue 4: Move Not Sent

**Symptom**:
```
[GameStore] Move executed: {...}
// No "Sending move to backend" message
```

**Causes**:
1. `opponentType` is not 'computer'
2. `gameId` is null or undefined

**Debug**:
```javascript
// In browser console
const state = window.__ZUSTAND_DEV_STORE_STATE__
console.log('opponentType:', state.opponentType)
console.log('gameId:', state.gameId)
console.log('mode:', state.mode)
```

**Solution**: Verify game was started correctly with computer mode

---

### Issue 5: Invalid Move Format Error

**Symptom**:
```
[GameService] Received: {id: 'req-2', type: 'response', error: 'invalid move format'}
```

**Cause**: Move sent in wrong format

**Debug**: Check what format was sent:
```
[GameStore] Sending move to backend: e2e4  // WRONG (UCI)
[GameStore] Sending move to backend: e4    // CORRECT (SAN)
```

**Solution**: Verify `move.san` is used, not UCI format

---

### Issue 6: Computer Doesn't Respond

**Symptom**:
```
[GameStore] Move sent successfully, waiting for computer...
// No move event received
```

**Causes**:
1. Backend Stockfish not working
2. Event listener not set up
3. WebSocket disconnected

**Debug**:

**Check backend logs**:
```bash
docker compose logs -f api | grep Computer
# Should see: "Computer (black) in game xxx plays: e7e5"
```

**Check event listener**:
```
[GameStore] Setting up event listeners  // Should appear on game start
```

**Check WebSocket connection**:
```javascript
// In browser console
window.__WS_CONNECTION_STATUS__
// Or check Network tab → WS connection should be open
```

---

### Issue 7: Computer Move Invalid

**Symptom**:
```
[GameStore] Invalid computer move: e7e5
```

**Cause**: Computer sent move in wrong format or illegal move

**Debug**:
```
[GameService] Received: {type: 'event', event: 'move', data: {moveUCI: '...', ...}}
```

**Check**: What's in `moveUCI`? Should be UCI format like "e7e5"

---

## 🧪 Step-by-Step Testing

### Test 1: Verify Backend

```bash
# Terminal 1: Start backend
cd backend
docker compose up api

# Terminal 2: Test with wscat
npm install -g wscat
wscat -c "ws://localhost:8080/ws?user_id=25d30da5-0cc4-4f5a-8c88-69d0f90b004c"

# Send create game
{"id":"1","type":"request","action":"createGame","data":{"mode":"human_vs_computer","difficulty":"easy","timeControl":"5+0"}}

# Should get success response with gameId
```

If this works, backend is fine. Issue is in frontend.

---

### Test 2: Verify Frontend Connection

```bash
# Start frontend
cd frontend/app
npm run dev

# Open browser to http://localhost:5173
# Open DevTools Console
# Click "Play vs Computer"
# Look for:
[GameStore] Starting computer game: ...
[GameStore] Connecting to WebSocket...
[GameService] Connected to WebSocket  // ← Should see this
```

If you don't see "Connected to WebSocket", check:
- Backend is running
- No CORS errors
- WebSocket URL is correct

---

### Test 3: Verify Game Creation

After successful connection, look for:
```
[GameStore] Creating game with mode: human_vs_computer, difficulty: easy
[GameService] Sending: {action: 'createGame', ...}
[GameService] Received: {success: true, data: {gameId: '...'}}  // ← Should see this
[GameStore] Game created: {gameId: '...'}
```

If you don't see success response:
- Check backend logs for errors
- Verify user ID exists in database

---

### Test 4: Verify Move Sending

Make a move (e2-e4) and look for:
```
[GameStore] Move executed: {san: 'e4', ...}
[GameStore] Sending move to backend: e4  // ← Should be SAN, not e2e4
[GameService] Sending: {action: 'makeMove', data: {move: 'e4'}}
[GameService] Received: {success: true, ...}
[GameStore] Move sent successfully
```

If you don't see these logs, check selectSquare function is being called.

---

### Test 5: Verify Computer Response

After sending move, look for:
```
[GameService] Received: {type: 'event', event: 'move', data: {...}}
[GameStore] Move event received: {moveUCI: 'e7e5', ...}
[GameStore] Handling computer move: e7e5
[GameStore] Computer move applied successfully
```

If you don't see event:
- Check backend logs: `docker compose logs -f api | grep Computer`
- Verify Stockfish is running: `docker compose exec api which stockfish`

---

## 🔧 Advanced Debugging

### Check Zustand State:

```javascript
// In browser console
const state = useGameStore.getState()
console.log('Current state:', {
  mode: state.mode,
  opponentType: state.opponentType,
  gameId: state.gameId,
  playerColor: state.playerColor,
  isComputerThinking: state.isComputerThinking,
  connectionStatus: state.connectionStatus,
  gameError: state.gameError,
})
```

### Check WebSocket Messages:

```javascript
// In browser console - intercept WebSocket
const originalSend = WebSocket.prototype.send
WebSocket.prototype.send = function(data) {
  console.log('WS Send:', JSON.parse(data))
  originalSend.call(this, data)
}
```

### Enable Verbose Logging:

Already enabled! All logs are prefixed:
- `[GameStore]` - State management logs
- `[GameService]` - WebSocket service logs

### Network Tab:

1. Open DevTools → Network tab
2. Filter by "WS" (WebSocket)
3. Click on the connection
4. View "Messages" tab
5. See all sent/received messages

---

## 📋 Debug Checklist

When something doesn't work, check in this order:

- [ ] Backend is running (`docker compose ps`)
- [ ] Frontend dev server is running (`npm run dev`)
- [ ] Console shows "Connected to WebSocket"
- [ ] User exists in database
- [ ] Game created successfully (check for gameId in logs)
- [ ] Move format is SAN (e4) not UCI (e2e4)
- [ ] Event listener set up (check logs)
- [ ] Stockfish is working (check backend logs)
- [ ] WebSocket connection is still open (check Network tab)

---

## 🎯 Expected Log Flow (Complete Game)

```
// === GAME START ===
[GameStore] Starting computer game: {color: 'white', difficulty: 'medium', userId: '...'}
[GameStore] Connecting to WebSocket...
[GameService] Connected to WebSocket
[GameStore] WebSocket connected
[GameStore] Creating game with mode: human_vs_computer, difficulty: medium
[GameService] Sending: {id: 'req-1', type: 'request', action: 'createGame', ...}
[GameService] Received: {id: 'req-1', type: 'response', success: true, data: {gameId: '...'}}
[GameStore] Game created: {gameId: '...', status: 'active', side: 'white', ...}
[GameStore] Setting up event listeners
[GameStore] Computer game started successfully

// === MOVE 1 (Player) ===
[GameStore] Move executed: {from: 'e2', to: 'e4', san: 'e4', opponentType: 'computer', gameId: '...'}
[GameStore] Sending move to backend: e4
[GameService] Sending: {id: 'req-2', type: 'request', action: 'makeMove', data: {gameId: '...', move: 'e4'}}
[GameStore] Move sent successfully, waiting for computer...
[GameService] Received: {id: 'req-2', type: 'response', success: true, data: {moveNumber: 1, san: 'e4', ...}}

// === MOVE 2 (Computer) ===
[GameService] Received: {type: 'event', event: 'move', data: {moveSAN: 'e5', moveUCI: 'e7e5', ...}}
[GameStore] Move event received: {gameId: '...', moveSAN: 'e5', moveUCI: 'e7e5', ...}
[GameStore] Handling computer move: e7e5
[GameStore] Computer move parsed: {san: 'e5', from: 'e7', to: 'e5', color: 'b', ...}
[GameStore] Computer move applied successfully

// === MOVE 3 (Player) ===
[GameStore] Move executed: {from: 'g1', to: 'f3', san: 'Nf3', ...}
[GameStore] Sending move to backend: Nf3
...
```

---

## 🆘 Still Not Working?

1. **Copy all console logs** and paste them in a file
2. **Check backend logs**: `docker compose logs api > backend-logs.txt`
3. **Check database**: Verify user exists, games table, moves table
4. **Network tab**: Export HAR file with WebSocket messages
5. **Compare with expected flow** above

---

## ✅ Success Indicators

You know it's working when you see:

✅ `[GameService] Connected to WebSocket`
✅ `[GameStore] Game created: {gameId: '...'}`
✅ `[GameStore] Sending move to backend: e4` (SAN format)
✅ `[GameService] Received: {type: 'event', event: 'move', ...}`
✅ `[GameStore] Computer move applied successfully`

If you see all of these, the computer player is working! 🎉

---

**Happy Debugging!** 🔍
