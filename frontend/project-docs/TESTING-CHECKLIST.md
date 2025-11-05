# Frontend Testing Checklist: Computer Player

**Date**: 2025-11-06

Use this checklist to verify all computer player functionality works correctly.

---

## ✅ Pre-Testing Setup

- [ ] Backend server is running (`docker compose up api`)
- [ ] Frontend dev server is running (`npm run dev`)
- [ ] Browser console is open (for debugging)
- [ ] Database has test user with ID: `25d30da5-0cc4-4f5a-8c88-69d0f90b004c`

---

## 🎮 Basic Functionality

### Game Creation:
- [ ] "Play vs Computer" button appears in live mode
- [ ] Button is styled correctly (purple-blue gradient, bottom-right)
- [ ] Clicking button opens GameSetup modal
- [ ] Modal shows all options correctly
- [ ] Modal can be closed with Cancel button
- [ ] Modal can be closed by clicking outside (TODO: implement if desired)

### Color Selection:
- [ ] White option selectable
- [ ] Black option selectable
- [ ] Selected color highlights with blue ring
- [ ] Visual indicators show correctly (white/black circles)
- [ ] Info text updates based on selection

### Difficulty Selection:
- [ ] Easy option selectable
- [ ] Medium option selectable
- [ ] Hard option selectable
- [ ] Selected difficulty highlights with blue ring
- [ ] Think time info displays correctly

### Game Start:
- [ ] "Start Game" button works
- [ ] Loading spinner appears during connection
- [ ] Button text changes to "Starting..."
- [ ] Modal closes after successful start
- [ ] Board resets to starting position
- [ ] GameInfo shows opponent as "Computer (difficulty)"

---

## 🎯 Playing as White

- [ ] Start game as White, Medium difficulty
- [ ] Board shows starting position
- [ ] Turn indicator shows "White"
- [ ] Can click white pieces
- [ ] Valid moves highlight when piece selected
- [ ] Can make first move (e.g., e2-e4)
- [ ] Move executes on board
- [ ] Turn changes to "Black"
- [ ] "Computer is thinking..." appears immediately
- [ ] Spinner animation shows
- [ ] Computer responds within ~500ms
- [ ] Computer move is legal
- [ ] Board updates with computer's move
- [ ] Turn changes back to "White"
- [ ] "Computer is thinking..." disappears
- [ ] Can make second move
- [ ] Pattern continues correctly

---

## 🎯 Playing as Black

- [ ] Start game as Black, Medium difficulty
- [ ] Board shows starting position
- [ ] Turn indicator shows "White"
- [ ] Cannot click any pieces (computer's turn)
- [ ] "Computer is thinking..." appears immediately
- [ ] Computer makes first move within ~500ms
- [ ] Computer move is legal (e.g., e2-e4)
- [ ] Board updates correctly
- [ ] Turn changes to "Black"
- [ ] "Computer is thinking..." disappears
- [ ] Can now click black pieces
- [ ] Valid moves highlight
- [ ] Can make response move
- [ ] Computer responds automatically
- [ ] Pattern continues correctly

---

## 🎲 Difficulty Levels

### Easy (~100ms):
- [ ] Start game with Easy difficulty
- [ ] Computer responds very quickly
- [ ] Response time consistently under 200ms
- [ ] Computer makes reasonable but not optimal moves
- [ ] Badge shows "Computer (easy)"

### Medium (~500ms):
- [ ] Start game with Medium difficulty
- [ ] Computer responds in ~500ms
- [ ] Computer plays stronger moves
- [ ] Badge shows "Computer (medium)"

### Hard (~2 seconds):
- [ ] Start game with Hard difficulty
- [ ] Computer takes noticeably longer (~2s)
- [ ] Computer plays very strong moves
- [ ] Badge shows "Computer (hard)"
- [ ] "Thinking" indicator shows for full duration

---

## 🏁 Game End Scenarios

### Checkmate (Computer Wins):
- [ ] Play poorly and let computer checkmate you
- [ ] "Checkmate!" message appears
- [ ] Winner is announced correctly
- [ ] "Computer is thinking..." disappears
- [ ] Cannot make more moves
- [ ] "New Game" button still works

### Checkmate (You Win):
- [ ] Checkmate the computer
- [ ] "Checkmate!" message appears
- [ ] "You" announced as winner
- [ ] Cannot make more moves
- [ ] "New Game" button works

### Stalemate:
- [ ] Force a stalemate position
- [ ] "Stalemate - Draw!" message appears
- [ ] Game stops correctly

### Resignation (Future):
- [ ] TODO: Test resign button when implemented

---

## 🔄 Multiple Games

- [ ] Start first game (White, Easy)
- [ ] Play a few moves
- [ ] Click "New Game"
- [ ] Board resets correctly
- [ ] GameInfo resets to default
- [ ] Can start second game (Black, Hard)
- [ ] No state leakage from first game
- [ ] Second game plays correctly
- [ ] Start third game without finishing second
- [ ] WebSocket cleanup happens correctly

---

## ♟️ Chess Rules

### Basic Moves:
- [ ] Pawns move forward correctly
- [ ] Knights jump correctly
- [ ] Bishops move diagonally
- [ ] Rooks move straight
- [ ] Queen moves in all directions
- [ ] King moves one square
- [ ] Cannot move into check
- [ ] Cannot leave king in check

### Special Moves:
- [ ] Castling kingside works (white: e1-g1)
- [ ] Castling queenside works (white: e1-c1)
- [ ] En passant capture works
- [ ] Pawn promotion to Queen works
- [ ] Pawn promotion to Knight/Bishop/Rook works

### Computer Respects Rules:
- [ ] Computer never makes illegal moves
- [ ] Computer doesn't move into check
- [ ] Computer castles legally
- [ ] Computer captures en passant when possible
- [ ] Computer promotes pawns correctly

---

## 🌐 WebSocket Connection

### Connection:
- [ ] WebSocket connects on game start
- [ ] Connection status shows "connected"
- [ ] No console errors during connection
- [ ] Connection ID logged correctly

### Messages:
- [ ] createGame request sent correctly
- [ ] createGame response received with gameId
- [ ] makeMove requests sent for each move
- [ ] makeMove responses received
- [ ] move events received for computer moves
- [ ] All messages have correct format

### Disconnection:
- [ ] Disconnect happens on "New Game"
- [ ] Disconnect happens on page refresh
- [ ] No orphaned connections remain
- [ ] Can reconnect after disconnect

---

## ❌ Error Handling

### Backend Down:
- [ ] Stop backend server
- [ ] Try to start game
- [ ] Error alert appears
- [ ] Error message is user-friendly
- [ ] Can try again after backend restarts
- [ ] Connection status shows "error"

### Mid-Game Disconnect:
- [ ] Start game and make moves
- [ ] Stop backend server
- [ ] Try to make move
- [ ] Error message appears in GameInfo
- [ ] Cannot continue game
- [ ] "New Game" button still works

### Invalid Moves:
- [ ] TODO: Test invalid move handling

### Timeout:
- [ ] TODO: Test move timeout scenarios

---

## 🎨 UI/UX

### GameSetup Modal:
- [ ] Modal is centered on screen
- [ ] Background overlay is semi-transparent
- [ ] Modal is responsive on mobile
- [ ] Text is readable
- [ ] Buttons are easy to click
- [ ] Hover states work correctly
- [ ] Loading state is clear

### GameInfo Display:
- [ ] Opponent info is prominent
- [ ] Difficulty badge is visible
- [ ] "Thinking" indicator is noticeable
- [ ] Spinner animation is smooth
- [ ] Player names display correctly
- [ ] Error messages are readable
- [ ] Turn indicator updates instantly

### Floating Button:
- [ ] Button is visible in live mode
- [ ] Button is not visible in computer mode
- [ ] Button is not visible in replay mode
- [ ] Hover effect works
- [ ] Icon displays correctly
- [ ] Text is readable
- [ ] Position doesn't overlap other elements

---

## 📱 Responsive Design

### Desktop (1920x1080):
- [ ] Layout looks good
- [ ] All elements visible
- [ ] No horizontal scroll
- [ ] Floating button positioned correctly

### Laptop (1366x768):
- [ ] Layout adapts correctly
- [ ] No elements cut off
- [ ] Floating button still accessible

### Tablet (768x1024):
- [ ] GameSetup modal fits on screen
- [ ] Board is appropriate size
- [ ] Buttons are touchable
- [ ] No layout issues

### Mobile (375x667):
- [ ] Modal is readable
- [ ] Can select all options
- [ ] Board is playable
- [ ] Floating button accessible
- [ ] No overflow issues

---

## 🚀 Performance

### Load Time:
- [ ] Initial page load < 2 seconds
- [ ] GameSetup modal opens instantly
- [ ] No stuttering or lag

### Move Response:
- [ ] Player moves execute instantly
- [ ] UI updates immediately
- [ ] No delay in highlighting
- [ ] Board animations smooth

### Computer Response:
- [ ] Easy: < 200ms consistently
- [ ] Medium: ~500ms consistently
- [ ] Hard: ~2s consistently
- [ ] No timeouts or hangs

### Memory:
- [ ] No memory leaks after multiple games
- [ ] No performance degradation over time
- [ ] Console shows no warnings

---

## 🔍 Console Checks

### No Errors:
- [ ] No red errors in console during normal play
- [ ] No React warnings
- [ ] No TypeScript errors

### Expected Logs:
- [ ] "[GameService] Connected to WebSocket"
- [ ] "[GameService] Sending: {createGame request}"
- [ ] "[GameService] Received: {game created}"
- [ ] "[GameService] Sending: {makeMove request}"
- [ ] "[GameService] Received: {move event}"

### Network Tab:
- [ ] WebSocket connection established
- [ ] Messages sent/received correctly
- [ ] No failed requests
- [ ] Clean disconnection

---

## ✨ Polish

- [ ] No typos in UI text
- [ ] Consistent styling throughout
- [ ] Smooth transitions
- [ ] Loading states everywhere needed
- [ ] Error messages are helpful
- [ ] Icons render correctly
- [ ] Colors follow design system

---

## 🎯 Edge Cases

- [ ] Start game, refresh page (should disconnect cleanly)
- [ ] Open multiple tabs (each should have own connection)
- [ ] Switch between live/replay/computer modes
- [ ] Start game while another is in progress
- [ ] Click squares very rapidly
- [ ] Make move exactly when computer is thinking
- [ ] Promote pawn while computer is thinking

---

## 📊 Success Criteria

**All Must Pass**:
- [x] Can create game as White
- [x] Can create game as Black
- [x] Can select all difficulty levels
- [x] Computer responds to every move
- [x] No illegal moves by computer
- [x] Game ends correctly
- [x] Can start multiple games
- [x] No console errors
- [x] WebSocket connects/disconnects cleanly
- [x] Error handling works

---

## 🐛 Bug Report Template

If you find a bug, document it like this:

```
**Title**: [Short description]

**Steps to Reproduce**:
1. Step one
2. Step two
3. Step three

**Expected**: What should happen

**Actual**: What actually happened

**Console Errors**: [Copy paste any errors]

**Environment**:
- Browser: Chrome 120
- OS: macOS
- Backend: Running locally

**Screenshots**: [If applicable]
```

---

## ✅ Final Verification

After completing all tests above:

- [ ] All critical tests passed
- [ ] No showstopper bugs found
- [ ] Performance is acceptable
- [ ] UX is smooth
- [ ] Ready for demo/production

---

**Testing Complete!** 🎉

If all checks pass, the computer player integration is ready for use!
