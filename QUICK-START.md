# Quick Start Guide: Play Chess vs Computer

**Last Updated**: 2025-11-06

---

## 🚀 Get Started in 3 Minutes

### Step 1: Start the Backend
```bash
cd backend
docker compose up api
```

**Wait for**: `Server started on port 8080`

---

### Step 2: Start the Frontend
```bash
# Open new terminal
cd frontend/app
npm install    # First time only
npm run dev
```

**Wait for**: `Local: http://localhost:5173`

---

### Step 3: Play!

1. **Open Browser**: http://localhost:5173
2. **Click**: "Play vs Computer" (bottom-right button)
3. **Choose**:
   - Color: White or Black
   - Difficulty: Easy, Medium, or Hard
4. **Click**: "Start Game"
5. **Play**: Click pieces to move!

---

## 🎮 Controls

### Making Moves:
- Click a piece → See valid moves highlighted
- Click destination → Move executes
- Computer responds automatically

### Special Moves:
- **Castling**: Click king, then click 2 squares (e1-g1 or e1-c1)
- **Promotion**: Move pawn to last rank → Select piece in dialog
- **En Passant**: Just works automatically!

### Game Controls:
- **New Game**: Click "New Game" button in GameInfo panel
- **Resign**: Coming soon!

---

## ⚙️ Difficulty Levels

| Level  | Think Time | Strength          | Best For           |
|--------|------------|-------------------|--------------------|
| Easy   | ~100ms     | Makes weak moves  | Beginners          |
| Medium | ~500ms     | Solid play        | Intermediate       |
| Hard   | ~2 seconds | Very strong       | Advanced players   |

---

## 🐛 Troubleshooting

### "Connection failed" error?
```bash
# Check if backend is running
docker compose ps

# Restart backend
docker compose restart api
```

### "Invalid user ID" error?
```bash
# Create test user in database
docker compose exec db psql -U postgres -d chess_coach

# In psql:
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

### Frontend not loading?
```bash
# Clear cache and reinstall
cd frontend/app
rm -rf node_modules
npm install
npm run dev
```

### Computer not responding?
- Check backend logs: `docker compose logs -f api`
- Look for: `Computer (black) in game xxx plays: e7e5`
- Verify Stockfish: `docker compose exec api which stockfish`

---

## 📱 Browser Support

- ✅ Chrome 100+
- ✅ Firefox 100+
- ✅ Safari 15+
- ✅ Edge 100+

---

## 🎯 Tips for Playing

1. **Play as White first** - Easier to start
2. **Try Easy mode** - Get comfortable with the interface
3. **Watch the "thinking" indicator** - Don't click during computer's turn
4. **Make slow, deliberate moves** - Think before clicking
5. **Study computer's responses** - Learn from them!

---

## 📚 More Information

- **Full Documentation**: See [PROJECT-COMPLETE.md](PROJECT-COMPLETE.md)
- **Backend Testing**: See [backend/project-docs/computer_player/TESTING-GUIDE.md](backend/project-docs/computer_player/TESTING-GUIDE.md)
- **Frontend Testing**: See [frontend/project-docs/TESTING-CHECKLIST.md](frontend/project-docs/TESTING-CHECKLIST.md)

---

## 🆘 Getting Help

**Check Logs**:
```bash
# Backend logs
docker compose logs -f api

# Frontend console
# Open browser DevTools (F12) → Console tab
```

**Common Issues**:
- Backend not starting? → Check Docker is running
- Frontend errors? → Check Node.js version (18+)
- Database errors? → Run `docker compose up db` first

---

## ✅ Quick Verification

Everything working if you see:
- ✅ Backend: `Server started on port 8080`
- ✅ Frontend: `Local: http://localhost:5173`
- ✅ Browser: Chess board appears
- ✅ Button: "Play vs Computer" visible
- ✅ Modal: Opens when button clicked
- ✅ Game: Starts when "Start Game" clicked
- ✅ Computer: Responds to your moves

---

**Enjoy playing chess!** ♟️🎉

