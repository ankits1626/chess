# Step 4 Overview: Refactor GameManager

**Status**: Ready to implement

**Time**: 2 hours

---

## 📊 What Changes

### Before (Human vs Human only):
```
GameManager
├── pendingGames: map[string]*PendingGame
│   └── WhitePlayer: *websocket.Client ❌
└── activeGames: map[string]*ActiveGame
    ├── WhitePlayer: *websocket.Client ❌
    └── BlackPlayer: *websocket.Client ❌
```

### After (Any game mode):
```
GameManager
├── aiService: player.AIService ✅
├── playerFactory: *player.PlayerFactory ✅
├── pendingGames: map[string]*PendingGame
│   ├── Mode: player.GameMode ✅
│   ├── WhitePlayer: player.Player ✅
│   └── Difficulty: string ✅
└── activeGames: map[string]*ActiveGame
    ├── Mode: player.GameMode ✅
    ├── WhitePlayer: player.Player ✅
    └── BlackPlayer: player.Player ✅
```

---

## 🎯 Key Changes Summary

1. **GameManager gets AI capabilities**:
   - Add `aiService` field
   - Add `playerFactory` field
   - Update constructor to accept AIService

2. **Games track their mode**:
   - Add `Mode` field to PendingGame and ActiveGame
   - Add `Difficulty` field to PendingGame

3. **Players become interfaces**:
   - Change `*websocket.Client` to `player.Player`
   - Works for both humans and computers

4. **Smart player creation**:
   - `CreatePendingGame` creates player based on mode
   - `ActivateGame` creates second player based on mode

5. **Computer move handling**:
   - New `HandleComputerMove()` method
   - Automatically makes move when it's computer's turn
   - Notifies both players

---

## 📝 Files to Create

None - only modify existing files.

---

## 📝 Files to Modify

1. **game_manager.go** - Entire refactor

---

## 🔧 Breaking Changes

These files will have compilation errors after Step 4 (fixed in Steps 5 & 6):

- `create_game.go` - CreatePendingGame signature changed
- `join_game.go` - May need updates
- `make_move.go` - Needs to call HandleComputerMove
- Server initialization - NewGameManager needs AIService

**Don't worry!** This is expected. We fix these systematically in Steps 5 & 6.

---

## 📚 Implementation Resources

- **Full Guide**: [04-refactor-game-manager.md](./04-refactor-game-manager.md)
- **Quick Checklist**: [STEP-4-CHECKLIST.md](./STEP-4-CHECKLIST.md)

---

## 🎯 Success Criteria

After Step 4:
- ✅ GameManager has AIService and PlayerFactory
- ✅ PendingGame and ActiveGame use Player interface
- ✅ CreatePendingGame accepts mode and difficulty
- ✅ ActivateGame creates players by mode
- ✅ HandleComputerMove method exists
- ✅ game_manager.go compiles
- ⚠️ Other files have errors (expected!)

---

## 🚀 Next Steps

After completing Step 4:
- **Step 5**: Update handlers (create_game, join_game)
- **Step 6**: Update make_move to trigger computer moves

Total remaining: ~2.5 hours

---

**Ready to start?** Go to: [STEP-4-CHECKLIST.md](./STEP-4-CHECKLIST.md)
