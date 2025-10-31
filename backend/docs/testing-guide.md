# Testing Guide - Chess Coach Backend

**Quick guide to test all components manually**

---

## 1. HTTP Endpoints Testing

### Health Check
```bash
curl http://localhost:8080/health
```

**Expected Response:**
```json
{"service":"chess-coach-backend","status":"ok"}
```

### API Info
```bash
curl http://localhost:8080/
```

**Expected Response:**
```json
{"message":"Chess Coach Backend API","version":"0.2.0"}
```

---

## 2. Database Testing

### Check Tables Exist
```bash
make db-shell
# Then in psql:
\dt
```

**Expected Output:**
```
 tablename
-----------
 users
 games
 moves
```

### Insert Test User
```sql
-- In psql (make db-shell)
INSERT INTO users (email, username)
VALUES ('test@example.com', 'testuser')
RETURNING *;
```

### Query Users
```sql
SELECT * FROM users;
```

### Insert Test Game
```sql
-- First, get a user ID
SELECT id FROM users LIMIT 1;

-- Then insert a game (replace YOUR_USER_ID)
INSERT INTO games (user_id, pgn, title)
VALUES ('YOUR_USER_ID', '1. e4 e5 2. Nf3 Nc6', 'Test Game')
RETURNING *;
```

### Verify Foreign Keys Work
```sql
-- This should show the game with user info
SELECT g.title, g.pgn, u.username
FROM games g
JOIN users u ON g.user_id = u.id;
```

### Exit psql
```sql
\q
```

---

## 3. WebSocket Testing (Browser)

### Option A: Use the Test HTML File

1. **Open the test file in browser:**
```bash
open backend/test-websocket.html
# Or manually open: file:///Users/ankit/code/learn/chess-coach/backend/test-websocket.html
```

2. **Click "Connect"** - You should see:
   - "Connected!" message
   - Backend logs showing: `new websocket connection`

3. **Click "Send Ping"** - You should see:
   - The ping message echoed back to all connected clients

4. **Open in multiple browser tabs** to test broadcasting

### Option B: Use Browser Console

1. **Open browser console** (Chrome DevTools)

2. **Run this JavaScript:**
```javascript
// Connect to WebSocket
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onopen = () => {
    console.log('✅ Connected!');

    // Send a ping
    ws.send(JSON.stringify({
        type: 'ping',
        payload: { timestamp: Date.now() }
    }));
};

ws.onmessage = (event) => {
    console.log('📨 Received:', event.data);
};

ws.onerror = (error) => {
    console.error('❌ Error:', error);
};

ws.onclose = () => {
    console.log('🔌 Disconnected');
};
```

3. **Send test messages:**
```javascript
// Send a move
ws.send(JSON.stringify({
    type: 'move',
    payload: {
        game_id: 'test-123',
        move: 'e4'
    }
}));

// Send custom message
ws.send('Hello from browser!');
```

### Option C: Use wscat (Command Line)

1. **Install wscat:**
```bash
npm install -g wscat
```

2. **Connect:**
```bash
wscat -c ws://localhost:8080/ws -H "Origin: http://localhost:5173"
```

3. **Send messages:**
```
> {"type":"ping","payload":{}}
> {"type":"move","payload":{"game_id":"test-123","move":"e4"}}
> Hello, WebSocket!
```

4. **Press Ctrl+C to disconnect**

---

## 4. PGN Parser Testing

### Run Unit Tests
```bash
make test
```

**Expected Output:**
```
=== RUN   TestPGNParser_Parse
=== RUN   TestPGNParser_Parse/valid_PGN_with_moves
--- PASS: TestPGNParser_Parse (0.00s)
    --- PASS: TestPGNParser_Parse/valid_PGN_with_moves (0.00s)
...
PASS
ok  	chess-coach/backend/internal/game	0.506s
```

### Test PGN Parser Manually (Go Playground)

Create a simple test file:

**File: `backend/test/pgn_manual_test.go`**
```go
package main

import (
	"fmt"
	"chess-coach/backend/internal/game"
)

func main() {
	parser := game.NewPGNParser()

	// Test 1: Valid PGN
	pgn1 := `1. e4 e5 2. Nf3 Nc6 3. Bb5`
	gameObj, err := parser.Parse(pgn1)
	if err != nil {
		fmt.Printf("❌ Parse error: %v\n", err)
		return
	}
	fmt.Printf("✅ Parsed successfully\n")

	moves := parser.GetMoves(gameObj)
	fmt.Printf("📝 Moves: %v\n", moves)

	fen := parser.GetFEN(gameObj)
	fmt.Printf("🎯 Final FEN: %s\n", fen)

	// Test 2: Invalid PGN
	pgn2 := `1. e4 e5 2. InvalidMove`
	err = parser.Validate(pgn2)
	if err != nil {
		fmt.Printf("❌ Validation failed (expected): %v\n", err)
	} else {
		fmt.Printf("⚠️  Invalid PGN passed validation (library is lenient)\n")
	}
}
```

**Run it:**
```bash
cd backend
go run test/pgn_manual_test.go
```

---

## 5. End-to-End Integration Test

### Full Workflow Test

```bash
# 1. Start services
make up

# 2. Wait for healthy
sleep 10

# 3. Check health
curl -s http://localhost:8080/health | jq

# 4. Create test user
docker exec chess-coach-postgres psql -U postgres -d chess_coach -c \
  "INSERT INTO users (email, username) VALUES ('player1@test.com', 'player1') RETURNING id;"

# 5. Get user ID (save it)
USER_ID=$(docker exec chess-coach-postgres psql -U postgres -d chess_coach -t -c \
  "SELECT id FROM users WHERE email='player1@test.com';")

echo "User ID: $USER_ID"

# 6. Create game
docker exec chess-coach-postgres psql -U postgres -d chess_coach -c \
  "INSERT INTO games (user_id, pgn, title) VALUES ('$USER_ID', '1. e4 e5 2. Nf3', 'Test Game') RETURNING id, title;"

# 7. Verify game exists
docker exec chess-coach-postgres psql -U postgres -d chess_coach -c \
  "SELECT g.title, u.username FROM games g JOIN users u ON g.user_id = u.id;"

# 8. Test WebSocket (in another terminal)
wscat -c ws://localhost:8080/ws -H "Origin: http://localhost:5173"
```

---

## 6. Automated Test Script

Create a comprehensive test script:

**File: `backend/test/integration_test.sh`**
```bash
#!/bin/bash

echo "🧪 Running Integration Tests..."
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Test counter
PASSED=0
FAILED=0

# Test function
test_endpoint() {
    local name=$1
    local url=$2
    local expected=$3

    response=$(curl -s "$url")

    if echo "$response" | grep -q "$expected"; then
        echo -e "${GREEN}✅ $name${NC}"
        ((PASSED++))
    else
        echo -e "${RED}❌ $name${NC}"
        echo "   Expected: $expected"
        echo "   Got: $response"
        ((FAILED++))
    fi
}

# Run tests
echo "📡 Testing HTTP Endpoints..."
test_endpoint "Health Check" "http://localhost:8080/health" "ok"
test_endpoint "API Info" "http://localhost:8080/" "Chess Coach Backend API"

echo ""
echo "🗄️  Testing Database..."
TABLES=$(docker exec chess-coach-postgres psql -U postgres -d chess_coach -t -c "SELECT COUNT(*) FROM pg_tables WHERE schemaname='public';")
if [ "$TABLES" -eq 3 ]; then
    echo -e "${GREEN}✅ Database has 3 tables${NC}"
    ((PASSED++))
else
    echo -e "${RED}❌ Database tables count incorrect${NC}"
    ((FAILED++))
fi

echo ""
echo "🧪 Testing PGN Parser..."
cd backend
if go test ./internal/game/... -run TestPGNParser_GetMoves > /dev/null 2>&1; then
    echo -e "${GREEN}✅ PGN Parser tests pass${NC}"
    ((PASSED++))
else
    echo -e "${RED}❌ PGN Parser tests fail${NC}"
    ((FAILED++))
fi

echo ""
echo "📊 Results:"
echo "   Passed: $PASSED"
echo "   Failed: $FAILED"

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}🎉 All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}❌ Some tests failed${NC}"
    exit 1
fi
```

**Make it executable and run:**
```bash
chmod +x backend/test/integration_test.sh
./backend/test/integration_test.sh
```

---

## 7. Load Testing (Optional)

### Test WebSocket Under Load

**Install hey (HTTP load tester):**
```bash
brew install hey
```

**Test HTTP endpoints:**
```bash
# 100 requests, 10 concurrent
hey -n 100 -c 10 http://localhost:8080/health
```

### Test WebSocket Connections

Create a simple load test:

**File: `backend/test/ws_load_test.js`**
```javascript
// Run with: node backend/test/ws_load_test.js
const WebSocket = require('ws');

const NUM_CLIENTS = 10;
const MESSAGES_PER_CLIENT = 5;

console.log(`🚀 Starting ${NUM_CLIENTS} WebSocket clients...`);

let connectedClients = 0;
let totalMessagesSent = 0;
let totalMessagesReceived = 0;

for (let i = 0; i < NUM_CLIENTS; i++) {
    const ws = new WebSocket('ws://localhost:8080/ws', {
        headers: { Origin: 'http://localhost:5173' }
    });

    ws.on('open', () => {
        connectedClients++;
        console.log(`✅ Client ${i+1} connected (${connectedClients}/${NUM_CLIENTS})`);

        // Send test messages
        for (let j = 0; j < MESSAGES_PER_CLIENT; j++) {
            ws.send(JSON.stringify({
                type: 'test',
                payload: {
                    client: i+1,
                    message: j+1
                }
            }));
            totalMessagesSent++;
        }
    });

    ws.on('message', (data) => {
        totalMessagesReceived++;
    });

    ws.on('error', (error) => {
        console.error(`❌ Client ${i+1} error:`, error.message);
    });
}

// Summary after 5 seconds
setTimeout(() => {
    console.log('\n📊 Load Test Summary:');
    console.log(`   Connected clients: ${connectedClients}`);
    console.log(`   Messages sent: ${totalMessagesSent}`);
    console.log(`   Messages received: ${totalMessagesReceived}`);
    process.exit(0);
}, 5000);
```

**Run it:**
```bash
npm install ws  # If not already installed
node backend/test/ws_load_test.js
```

---

## 8. Monitoring Logs

### Watch Backend Logs
```bash
make backend-logs
```

### Watch Database Logs
```bash
docker logs -f chess-coach-postgres
```

### Watch All Logs
```bash
make logs
```

### Filter for Specific Events
```bash
# Watch for WebSocket connections
docker logs -f chess-coach-backend | grep "websocket connection"

# Watch for database queries
docker logs -f chess-coach-backend | grep "database"

# Watch for errors
docker logs -f chess-coach-backend | grep "error"
```

---

## 9. Quick Test Checklist

Use this for daily verification:

```bash
# ✅ Services running
docker ps | grep chess-coach

# ✅ Health check
curl -s localhost:8080/health | jq .status

# ✅ Database accessible
docker exec chess-coach-postgres pg_isready -U postgres

# ✅ Tables exist
docker exec chess-coach-postgres psql -U postgres -d chess_coach -c "\dt" | grep "3 rows"

# ✅ Tests pass
cd backend && go test ./internal/game/... -run TestPGNParser_GetMoves

# ✅ WebSocket accepts connections
# (Open backend/test-websocket.html in browser and click Connect)
```

---

## 10. Troubleshooting

### Backend won't start
```bash
# Check logs
make backend-logs

# Common issues:
# 1. Port 8080 in use
lsof -ti:8080 | xargs kill -9

# 2. Go module issues
cd backend && go mod tidy

# 3. Rebuild from scratch
make down && make up
```

### Database connection fails
```bash
# Check postgres is healthy
docker ps | grep postgres

# Check connection string
echo $DATABASE_URL

# Verify from backend container
docker exec chess-coach-backend ping -c 1 postgres
```

### WebSocket won't connect
```bash
# Check CORS headers
curl -I http://localhost:8080/ws

# Try with correct origin
wscat -c ws://localhost:8080/ws -H "Origin: http://localhost:5173"

# Check backend logs for errors
make backend-logs | grep websocket
```

---

## Next Steps

Once you've verified everything works:

1. **Add API endpoints** for games (POST, GET, PUT)
2. **Implement authentication** (JWT)
3. **Integrate Claude AI** for move analysis
4. **Build the frontend** to connect to these APIs

Happy testing! 🚀
