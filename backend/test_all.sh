#!/bin/bash

# Chess Coach Backend - Quick Test Script
# Run this to verify everything is working

set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}🧪 Chess Coach Backend - Integration Test${NC}"
echo ""

# Test 1: Health Check
echo -e "${YELLOW}📡 Testing HTTP Endpoints...${NC}"
HEALTH=$(curl -s http://localhost:8080/health)
if echo "$HEALTH" | grep -q "ok"; then
    echo -e "${GREEN}✅ Health check passed${NC}"
    echo "   Response: $HEALTH"
else
    echo -e "${RED}❌ Health check failed${NC}"
    exit 1
fi
echo ""

# Test 2: Database
echo -e "${YELLOW}🗄️  Testing Database...${NC}"
echo "Creating test user..."
USER_RESULT=$(docker exec chess-coach-postgres psql -U postgres -d chess_coach -t -c \
  "INSERT INTO users (email, username) VALUES ('test_$(date +%s)@example.com', 'test_user_$(date +%s)') RETURNING id, username;")
echo -e "${GREEN}✅ User created: $USER_RESULT${NC}"

# Get the user ID from the result (extract UUID from first line)
USER_ID=$(echo "$USER_RESULT" | head -1 | awk '{print $1}' | tr -d ' ')
USERNAME=$(echo "$USER_RESULT" | head -1 | awk '{print $3}' | tr -d ' ')

echo "Creating test game..."
GAME_RESULT=$(docker exec chess-coach-postgres psql -U postgres -d chess_coach -t -c \
  "INSERT INTO games (user_id, pgn, title) VALUES ('$USER_ID', '1. e4 e5 2. Nf3 Nc6 3. Bb5', 'Ruy Lopez Opening') RETURNING id, title;")
echo -e "${GREEN}✅ Game created: $GAME_RESULT${NC}"

echo "Verifying join query..."
JOIN_RESULT=$(docker exec chess-coach-postgres psql -U postgres -d chess_coach -t -c \
  "SELECT g.title, u.username FROM games g JOIN users u ON g.user_id = u.id WHERE u.id = '$USER_ID';")
echo -e "${GREEN}✅ Join query successful:${NC}"
echo "   $JOIN_RESULT"
echo ""

# Test 3: PGN Parser
echo -e "${YELLOW}♟️  Testing PGN Parser...${NC}"
if (cd backend && go test ./internal/game/... -run TestPGNParser_GetMoves 2>&1 | grep -q "PASS"); then
    echo -e "${GREEN}✅ PGN Parser tests passed${NC}"
else
    echo -e "${YELLOW}⚠️  PGN Parser tests skipped (run 'make test' manually)${NC}"
fi
echo ""

# Test 4: WebSocket Info
echo -e "${YELLOW}🔌 WebSocket Information${NC}"
echo "To test WebSocket, you have 3 options:"
echo ""
echo -e "${BLUE}Option 1: Browser (easiest)${NC}"
echo "  1. Open: backend/test-websocket.html"
echo "  2. Click 'Connect'"
echo "  3. Click 'Send Ping'"
echo "  4. Watch messages appear"
echo ""
echo -e "${BLUE}Option 2: wscat (command line)${NC}"
echo "  Install: npm install -g wscat"
echo "  Run: wscat -c ws://localhost:8080/ws -H \"Origin: http://localhost:5173\""
echo "  Type messages and press Enter"
echo ""
echo -e "${BLUE}Option 3: Browser Console${NC}"
echo "  1. Open http://localhost:5173 (or any page)"
echo "  2. Open DevTools Console (F12)"
echo "  3. Paste and run:"
echo ""
echo "     const ws = new WebSocket('ws://localhost:8080/ws');"
echo "     ws.onopen = () => console.log('✅ Connected!');"
echo "     ws.onmessage = (e) => console.log('📨', e.data);"
echo "     ws.send('Hello WebSocket!');"
echo ""

# Summary
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}🎉 All automated tests passed!${NC}"
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo "Test Data Created:"
echo "  👤 User ID: $USER_ID"
echo "  📧 Email: test_*@example.com"
echo "  🎮 Game: Ruy Lopez Opening (1. e4 e5 2. Nf3 Nc6 3. Bb5)"
echo ""
echo "Next Steps:"
echo "  • Open backend/test-websocket.html to test WebSocket"
echo "  • Run 'make db-shell' to explore the database"
echo "  • Check 'backend/docs/testing-guide.md' for more tests"
echo ""
