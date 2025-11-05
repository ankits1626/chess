# WebSocket Implementation - Step-by-Step Guide

**Goal**: Implement production-ready WebSocket infrastructure for real-time chess gameplay

**Approach**: Break down into small, testable increments

---

## 📁 Documentation Structure

```
websockets/
├── README.md                           # This file - overview and navigation
├── 00-overview-and-fixes.md           # Analysis and required fixes
├── 01-foundation.md                   # Phase 1: Core WebSocket infrastructure
├── 02-integration.md                  # Phase 2: Connect to existing app
├── 03-utilities.md                    # Phase 3: Helper functions and validation
├── 04-message-handlers.md             # Phase 4: Business logic handlers
├── 05-authentication.md               # Phase 5: Security and auth
├── 06-testing.md                      # Phase 6: Unit and integration tests
├── 07-deployment.md                   # Phase 7: Production readiness
└── implementation-checklist.md        # Step-by-step checklist
```

---

## 🎯 Implementation Phases

### Phase 1: Foundation (4-5 hours)
**Files**: [`01-foundation.md`](./01-foundation.md) ✅

Create core WebSocket components:
- [x] Message protocol definitions
- [x] Client wrapper with pumps
- [x] Hub with connection registry
- [x] Room management

**Deliverable**: Compiling WebSocket package (no integration yet)

**Status**: ✅ Documentation complete - Ready to implement

---

### Phase 2: Integration (2-3 hours)
**Files**: [`02-integration.md`](./02-integration.md) ✅

Connect WebSocket to existing architecture:
- [ ] Create `upgrade.go` - HTTP → WebSocket handler
- [ ] Modify `app.go` to start hub
- [ ] Modify `server.go` to accept hub
- [ ] Modify `router.go` to register `/ws` endpoint
- [ ] Update `main.go` to wire dependencies

**Deliverable**: WebSocket endpoint accessible, connects to app

**Status**: ✅ Documentation complete - Ready to implement

---

### Phase 3: Utilities (1-2 hours)
**Files**: `03-utilities.md`

Helper functions for type conversions and validation:
- [ ] UUID conversion utilities
- [ ] Context handling helpers
- [ ] Error response helpers
- [ ] Message validation utilities

**Deliverable**: Type-safe utilities for handlers

---

### Phase 4: Message Handlers (6-8 hours)
**Files**: `04-message-handlers.md`

Implement business logic:
- [ ] `createGame` handler
- [ ] `joinGame` handler
- [ ] `makeMove` handler (with validation)
- [ ] `leaveGame` handler
- [ ] Error handling and edge cases

**Deliverable**: Full game flow working end-to-end

---

### Phase 5: Authentication (2-3 hours)
**Files**: `05-authentication.md`

Secure WebSocket connections:
- [ ] JWT token validation
- [ ] CORS configuration
- [ ] User context extraction
- [ ] Authorization checks

**Deliverable**: Secure WebSocket endpoint

---

### Phase 6: Testing (4-6 hours)
**Files**: `06-testing.md`

Comprehensive test coverage:
- [ ] Unit tests (message, client, hub)
- [ ] Integration tests (game flow)
- [ ] Manual testing with wscat
- [ ] Load testing basics

**Deliverable**: 80%+ test coverage

---

### Phase 7: Production (2-3 hours)
**Files**: `07-deployment.md`

Polish and deploy:
- [ ] Configuration management
- [ ] Monitoring and logging
- [ ] Graceful shutdown
- [ ] Performance optimization

**Deliverable**: Production-ready WebSocket server

---

## 📊 Progress Tracking

| Phase | Status | Duration | Completion |
|-------|--------|----------|------------|
| Phase 1: Foundation | 🟢 Complete | 4-5h | 100% |
| Phase 2: Integration | 🟢 Complete | 2-3h | 100% |
| Phase 3: Utilities | 🟢 Complete | 1-2h | 100% |
| Phase 4: Handlers | 🔴 Not Started | 6-8h | 0% |
| Phase 5: Authentication | 🔴 Not Started | 2-3h | 0% |
| Phase 6: Testing | 🔴 Not Started | 4-6h | 0% |
| Phase 7: Production | 🔴 Not Started | 2-3h | 0% |
| **TOTAL** | **🟡** | **21-30h** | **42%** |

**Legend**: 🔴 Not Started | 🟡 In Progress | 🟢 Complete

---

## 🚀 Quick Start

### Option A: Follow in Order (Recommended)
1. Start with `00-overview-and-fixes.md` - Understand the issues
2. Follow `01-foundation.md` - Build core components
3. Continue through phases sequentially

### Option B: Jump to Specific Phase
- Each document is self-contained
- Prerequisites listed at top of each file
- Can implement phases independently (with caveats)

---

## 🎯 Success Criteria

After completing all phases, you should have:

✅ **Working WebSocket server**
- Accepts connections at `/ws`
- Handles authentication
- Manages multiple concurrent games

✅ **Complete game flow**
- Create game via WebSocket
- Join game (multiplayer)
- Make moves with validation
- Broadcast to opponents
- Persist to database

✅ **Production ready**
- Graceful shutdown
- Error handling
- Security hardened
- Test coverage 80%+
- Performance optimized

✅ **Well documented**
- All functions documented
- Integration points clear
- Testing instructions included

---

## 🛠️ Prerequisites

Before starting, ensure you have:
- [x] Go 1.25.3 installed
- [x] Existing chess-coach backend running
- [x] PostgreSQL database connected
- [x] Understanding of goroutines and channels
- [ ] `gorilla/websocket` package installed

---

## 📚 Reference Materials

### Your Existing Learning Materials
- `learnings/fundamentals/websocket/LESSON-1.md` - WebSocket basics
- `learnings/fundamentals/websocket/LESSON-2.md` - Connection management
- `learnings/fundamentals/websocket/LESSON-3.md` - Message protocols
- `learnings/fundamentals/websocket/HUB-AND-PUMP-EXPLAINED.md` - Core patterns

### External Resources
- [Gorilla WebSocket Docs](https://pkg.go.dev/github.com/gorilla/websocket)
- [WebSocket Protocol RFC](https://tools.ietf.org/html/rfc6455)

---

## 🔄 Iterative Approach

Each phase follows this pattern:

1. **📖 Read** - Understand the requirements
2. **📝 Plan** - Review code examples
3. **⌨️ Code** - Implement the phase
4. **✅ Test** - Verify it works
5. **📊 Review** - Check against success criteria
6. **➡️ Next** - Move to next phase

**After each phase, you should have working, testable code!**

---

## 💡 Key Principles

### 1. **Incremental Development**
- Each phase builds on the previous
- Test after every change
- Commit frequently

### 2. **Type Safety**
- Match existing database types (`pgtype.UUID`)
- Proper context propagation
- No `interface{}` without type assertions

### 3. **Integration First**
- Don't build in isolation
- Use existing repositories
- Follow existing patterns

### 4. **Security By Default**
- Validate all inputs
- Authenticate all connections
- Authorize all actions

---

## 🤝 Getting Help

If stuck on a phase:
1. Review the prerequisites section
2. Check the reference materials
3. Look at similar patterns in existing codebase
4. Break the phase into smaller steps

---

## 📝 Notes

- **Estimated total time**: 21-30 hours over 4-5 days
- **Can pause between phases** - Each phase is a natural checkpoint
- **Phase 4 is the longest** - Plan accordingly
- **Testing is crucial** - Don't skip Phase 6

---

**Ready to start? Begin with:** [`00-overview-and-fixes.md`](./00-overview-and-fixes.md)

Or jump straight to implementation: [`01-foundation.md`](./01-foundation.md)

---

**Last Updated**: 2025-11-05
**Status**: Planning Complete, Ready for Implementation
