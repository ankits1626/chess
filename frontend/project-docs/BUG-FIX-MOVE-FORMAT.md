# Bug Fix: Move Format Issue

**Date**: 2025-11-06

**Issue**: Moves not working when playing against computer

---

## 🐛 Problem

When a user made a move against the computer, the frontend was sending moves in **UCI format** (e.g., "e2e4") but the backend expected **SAN format** (e.g., "e4").

### Symptoms:
- User clicks to make a move
- Move appears to execute on the board
- But backend rejects the move
- Error: "invalid move format"
- Computer never responds

---

## 🔍 Root Cause

**Frontend Code** ([useGameStore.ts:327](../app/src/store/useGameStore.ts#L327)):
```typescript
// BEFORE (wrong):
const moveStr = `${selectedSquare}${square}`; // "e2e4" (UCI)
gameService.makeMove(gameId, moveStr);
```

**Backend Expected** ([make_move.go:40](../../backend/internal/websocket/handlers/make_move.go#L40)):
```go
moveSAN, err := websocket.RequireString(msg.Data, "move")
// Expects SAN format: "e4", "Nf3", "O-O", etc.
```

**Mismatch**: UCI vs SAN format

---

## ✅ Solution

Use the `move.san` property from chess.js, which contains the SAN format:

### Regular Moves:
```typescript
// AFTER (correct):
const move = game.move({ from: selectedSquare, to: square });
if (move) {
  // move.san contains "e4", "Nf3", etc.
  gameService.makeMove(gameId, move.san);
}
```

### Promotion Moves:
```typescript
// AFTER (correct):
const move = game.move({
  from: pendingMove.from,
  to: pendingMove.to,
  promotion: piece
});
if (move) {
  // move.san contains "e8=Q", "a1=N", etc.
  gameService.makeMove(gameId, move.san);
}
```

---

## 📝 Changes Made

**File**: `frontend/app/src/store/useGameStore.ts`

### Change 1: Regular Moves (Line ~327)
```diff
- const moveStr = `${selectedSquare}${square}`;
- gameService.makeMove(gameId, moveStr).then(() => {
+ gameService.makeMove(gameId, move.san).then(() => {
```

### Change 2: Promotion Moves (Line ~357)
```diff
  handlePromotion: (piece: PieceType) => {
    // ...
-   game.move({ from: pendingMove.from, to: pendingMove.to, promotion: piece });
+   const move = game.move({ from: pendingMove.from, to: pendingMove.to, promotion: piece });
    // ...
-   const moveStr = `${pendingMove.from}${pendingMove.to}${piece}`;
-   gameService.makeMove(gameId, moveStr).then(() => {
+   if (opponentType === 'computer' && gameId && move) {
+     gameService.makeMove(gameId, move.san).then(() => {
```

---

## 🧪 Testing

### Before Fix:
```
User makes move e2-e4
Frontend sends: "e2e4"
Backend rejects: "invalid move format"
No computer response
```

### After Fix:
```
User makes move e2-e4
Frontend sends: "e4"
Backend accepts: "e4"
Computer responds: "e5"
```

### Test Cases:
- [x] Regular pawn move (e2-e4 → "e4")
- [x] Knight move (g1-f3 → "Nf3")
- [x] Capture (d4xc5 → "dxc5")
- [x] Castling kingside (e1-g1 → "O-O")
- [x] Castling queenside (e1-c1 → "O-O-O")
- [x] Pawn promotion (e7-e8=Q → "e8=Q")
- [x] Check (Qf7+ → "Qf7+")
- [x] Checkmate (Qf7# → "Qf7#")

---

## 📚 Chess Notation Formats

### SAN (Standard Algebraic Notation):
**Used by**: Chess books, games, human-readable format

**Examples**:
- Pawn moves: `e4`, `d4`, `c5`
- Piece moves: `Nf3`, `Bc4`, `Qh5`
- Captures: `Nxe5`, `dxe5`
- Castling: `O-O` (kingside), `O-O-O` (queenside)
- Promotion: `e8=Q`, `a1=N`
- Check: `Qf7+`
- Checkmate: `Qf7#`

### UCI (Universal Chess Interface):
**Used by**: Chess engines, computer protocols

**Examples**:
- Pawn moves: `e2e4`, `d2d4`, `c7c5`
- Piece moves: `g1f3`, `f1c4`, `d8h5`
- Captures: `g1f3` (same as regular move)
- Castling: `e1g1` (kingside), `e1c1` (queenside)
- Promotion: `e7e8q`, `a2a1n`
- Check/Checkmate: No special notation

---

## 💡 Key Takeaway

**Frontend** uses chess.js which provides both formats:
- `move.san` - Standard Algebraic Notation (human-readable)
- `move.from` + `move.to` - UCI-like format

**Backend** expects SAN format for validation and storage.

**Always use** `move.san` when sending moves to the backend!

---

## 🎯 Impact

**Before**: Computer player functionality completely broken

**After**: Full functionality restored
- ✅ Moves work correctly
- ✅ Computer responds to each move
- ✅ Promotions work
- ✅ Special moves (castling, en passant) work

---

## 🔜 Future Improvements

Consider adding validation to prevent this type of issue:

1. **Frontend validation**:
   ```typescript
   // Validate move format before sending
   if (!move.san || move.san.length === 0) {
     throw new Error('Invalid move format');
   }
   ```

2. **Better error messages**:
   ```typescript
   .catch(error => {
     const message = error.message.includes('invalid move format')
       ? 'Invalid move. Please try again.'
       : error.message;
     set({ gameError: message });
   });
   ```

3. **Type safety**:
   ```typescript
   // Add JSDoc to document expected format
   /**
    * Send a move to the backend
    * @param gameId - The game ID
    * @param moveSAN - Move in SAN format (e.g., "e4", "Nf3")
    */
   async makeMove(gameId: string, moveSAN: string): Promise<MoveResult>
   ```

---

**Status**: ✅ Fixed and tested

**Verified**: 2025-11-06
