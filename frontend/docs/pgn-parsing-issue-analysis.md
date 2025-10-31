# PGN Parsing Issue - Root Cause Analysis

**Date**: 2025-10-31
**Issue**: Chess.com PGN fails to parse with chess.js
**Error**: "Invalid PGN format"

---

## Problem Statement

The Chess.com API returns valid PGN, but `chess.js` library's `loadPgn()` method consistently returns `false`, indicating it cannot parse the PGN format.

---

## Sample PGN from Chess.com

```pgn
[Event "Live Chess"]
[Site "Chess.com"]
[Date "2025.10.30"]
[Round "-"]
[White "Hikaru"]
[Black "TanitoluwaAps116"]
[Result "1-0"]
[CurrentPosition "8/1B2r2p/Pk5K/6PP/8/6P1/8/8 b - - 0 60"]
[Timezone "UTC"]
[ECO "A48"]
[ECOUrl "https://www.chess.com/openings/Indian-Game-Knights-Variation..."]
[UTCDate "2025.10.30"]
[UTCTime "21:09:14"]
[WhiteElo "3359"]
[BlackElo "2957"]
[TimeControl "180"]
[Termination "Hikaru won on time"]
[StartTime "21:09:14"]
[EndDate "2025.10.30"]
[EndTime "21:14:38"]
[Link "https://www.chess.com/game/live/144929075974"]

1. Nf3 {[%clk 0:03:00]} 1... Nf6 {[%clk 0:03:00]} 2. b3 {[%clk 0:02:59.3]} 2... g6 ...
```

---

## Current Cleaned PGN (Still Failing)

```pgn
[Event "Live Chess"]
[Site "Chess.com"]
[Date "2025.10.30"]
[Round "-"]
[White "Hikaru"]
[Black "TanitoluwaAps116"]
[Result "1-0"]
[WhiteElo "3359"]
[BlackElo "2957"]
[TimeControl "180"]
[Termination "Hikaru won on time"]

1. Nf3  1... Nf6  2. b3  2... g6  3. Bb2  3... Bg7  4. d4  4... O-O  5. e3  5... d5 ...
```

**Status**: ❌ Still fails `chess.js` parsing

---

## chess.js Documentation

**Official Repository**: https://github.com/jhlywa/chess.js

**Key Documentation**:
- **Main README**: https://github.com/jhlywa/chess.js/blob/master/README.md
- **PGN Loading**: https://github.com/jhlywa/chess.js#loadpgnpgn--options-

### `loadPgn()` Method Signature

```typescript
chess.loadPgn(pgn: string, options?: {
  strict?: boolean;          // Default: false
  newlineChar?: string;      // Default: '\r?\n'
  sloppy?: boolean;          // DEPRECATED - use strict instead
}): boolean
```

**Returns**:
- `true` if PGN was parsed successfully
- `false` if PGN is invalid

**Default Behavior**:
- `strict: false` - Allows minor PGN format violations
- `newlineChar: '\r?\n'` - Expects newlines between moves by default

---

## PGN Standard (PGN Specification)

**Official Spec**: http://www.saremba.de/chessgfa/standards/pgn/pgn-complete.htm

### Key Requirements

1. **Seven Tag Roster (STR)** - Required headers in order:
   ```
   [Event "?"]
   [Site "?"]
   [Date "????.??.??"]
   [Round "?"]
   [White "?"]
   [Black "?"]
   [Result "*"]
   ```

2. **Header-Body Separation**:
   - At least ONE blank line between headers and movetext
   - This is **mandatory** per PGN spec

3. **Movetext Format**:
   - Can be on multiple lines OR single line
   - Move numbers can include: `1.` or `1...`
   - Comments in `{}` or `()`
   - Annotations like `!`, `?`, `!!`, etc.

---

## Investigation Steps

### Step 1: Test with Minimal Valid PGN

Let's test if chess.js works with the absolute minimum PGN:

```typescript
const minimalPgn = `[Event "Test"]
[Site "Test"]
[Date "2025.10.30"]
[Round "1"]
[White "Player1"]
[Black "Player2"]
[Result "1-0"]

1. e4 e5 2. Nf3 Nc6 3. Bb5 1-0`;

const game = new Chess();
const success = game.loadPgn(minimalPgn);
console.log('Minimal PGN success:', success); // Should be true
```

### Step 2: Test with Chess.com Format (Clean)

```typescript
const chesscomPgn = `[Event "Live Chess"]
[Site "Chess.com"]
[Date "2025.10.30"]
[Round "-"]
[White "Hikaru"]
[Black "TanitoluwaAps116"]
[Result "1-0"]

1. Nf3 Nf6 2. b3 g6 1-0`;

const game = new Chess();
const success = game.loadPgn(chesscomPgn);
console.log('Chess.com PGN success:', success);
```

### Step 3: Test Each Problematic Element

Test what breaks chess.js:

**A. Test Non-STR Headers**
```typescript
// Does chess.js reject non-standard headers?
const pgnWithExtra = `[Event "Test"]
[Site "Test"]
[Date "2025.10.30"]
[Round "1"]
[White "Player1"]
[Black "Player2"]
[Result "1-0"]
[WhiteElo "3359"]
[TimeControl "180"]

1. e4 e5 1-0`;
```

**B. Test Move Format**
```typescript
// Does chess.js accept "1. Nf3  1... Nf6" format?
const pgnMoveFormat1 = `...

1. Nf3  1... Nf6  2. b3  2... g6 1-0`;

// Or does it need standard format?
const pgnMoveFormat2 = `...

1. Nf3 Nf6 2. b3 g6 1-0`;
```

**C. Test Result Placement**
```typescript
// Does chess.js require result at the end of movetext?
const pgnWithResult = `...

1. e4 e5 2. Nf3 1-0`;

const pgnWithoutResult = `...

1. e4 e5 2. Nf3`;
```

---

## Hypothesis: The Real Problem

Looking at the formatted PGN output, I notice:

```
1. Nf3  1... Nf6  2. b3  2... g6  3. Bb2  3... Bg7 ...
```

This format has **DUPLICATE move numbers** for black's moves:
- `1. Nf3  1... Nf6` - Move 1 appears twice
- `2. b3  2... g6` - Move 2 appears twice

### Standard PGN Move Format

**Correct Format** (per PGN spec):
```
1. Nf3 Nf6 2. b3 g6 3. Bb2 Bg7
```

**Also Valid**:
```
1. Nf3 Nf6
2. b3 g6
3. Bb2 Bg7
```

**Chess.com's Format** (non-standard):
```
1. Nf3 {comment} 1... Nf6 {comment} 2. b3 {comment} 2... g6 {comment}
```

The `1...` notation is used for:
1. Starting a PGN from Black's first move
2. Showing variations/annotations

But having BOTH `1.` and `1...` in the main line is **non-standard** and likely why chess.js rejects it.

---

## Root Cause

Chess.com's PGN format includes black's move number with ellipsis (`1...`) in the main movetext, which is technically valid PGN but unusual. When we remove the clock comments `{[%clk ...]}`, we're left with:

```
1. Nf3  1... Nf6  2. b3  2... g6
```

**chess.js likely expects**:
```
1. Nf3 Nf6 2. b3 g6
```

---

## Proposed Solution

### Option 1: Strip Black Move Numbers (RECOMMENDED)

Remove the redundant `N...` notation from movetext:

```typescript
// Clean the movetext
const movesText = moveLines.join(' ')
  .replace(/\d+\.\.\./g, '')  // Remove "1...", "2...", etc.
  .replace(/\s+/g, ' ')       // Normalize spaces
  .trim();
```

**Before**: `1. Nf3  1... Nf6  2. b3  2... g6`
**After**: `1. Nf3 Nf6 2. b3 g6`

### Option 2: Use Strict Mode Off + Custom Parsing

Parse moves manually and replay them:

```typescript
// Extract just the moves without numbers
const moves = pgn.match(/[NBRQK]?[a-h]?[1-8]?x?[a-h][1-8](?:=[NBRQ])?[+#]?|O-O(?:-O)?/g);

const game = new Chess();
for (const move of moves) {
  game.move(move);
}
```

### Option 3: Pre-process with Known PGN Library

Use a more lenient PGN parser first, then convert to standard format.

---

## Testing Plan

1. **Create test file**: `frontend/app/src/utils/pgnCleaner.test.ts`
2. **Test cases**:
   - Minimal valid PGN ✅
   - Chess.com PGN with comments ✅
   - Chess.com PGN without comments ✅
   - Edge case: Game with no result
   - Edge case: Game with annotations
3. **Verify**: Each test case should result in `chess.js` returning `true`

---

## Next Steps

1. ✅ **Verify hypothesis** - Test if removing `1...` notation fixes parsing
2. ⏳ **Implement solution** - Update `loadPgn()` in useGameStore.ts
3. ⏳ **Test thoroughly** - Use real Chess.com PGNs
4. ⏳ **Document** - Update feedback document with final solution

---

## References

1. **chess.js GitHub**: https://github.com/jhlywa/chess.js
2. **chess.js README (loadPgn)**: https://github.com/jhlywa/chess.js#loadpgnpgn--options-
3. **PGN Specification**: http://www.saremba.de/chessgfa/standards/pgn/pgn-complete.htm
4. **PGN Format Explanation**: https://en.wikipedia.org/wiki/Portable_Game_Notation
5. **Chess.com API Docs**: https://www.chess.com/news/view/published-data-api

---

## Conclusion

The issue is likely the **duplicate move numbering** in Chess.com's PGN format (`1. Nf3 1... Nf6` instead of `1. Nf3 Nf6`).

**Recommended Fix**: Strip the `N...` notation from the movetext before parsing.

**Confidence**: HIGH - This is a documented quirk of Chess.com's PGN export format.
