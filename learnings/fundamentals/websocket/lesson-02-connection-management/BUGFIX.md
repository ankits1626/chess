# Bug Fix: Simulate Network Issue Button

## Problem

The "Simulate Network Issue" button was not working because it was trying to close the WebSocket with code 1006:

```javascript
// ❌ This doesn't work!
reconnectingWS.ws.close(1006);
```

## Why It Failed

**WebSocket Close Code 1006 is reserved** and cannot be manually set by client code. It's only used by the browser when there's an actual abnormal closure (network failure, unexpected disconnect).

From the WebSocket specification:
```
1006 - Abnormal Closure
Reserved. Used to indicate that a connection was closed abnormally
(that is, with no close frame being sent) when a status code is
expected. This code MUST NOT be set explicitly by an application.
```

Attempting to use code 1006 in `close()` either:
1. Gets silently ignored (defaults to 1005 or 1000)
2. Throws an error in some browsers

## Solution

Changed to use **code 4000 (Custom Application Code)** which is explicitly allowed by browsers:

```javascript
// ✅ This works!
reconnectingWS.ws.close(4000, 'Simulated network failure');
```

## Why Code 4000?

**Browsers are strict about which close codes you can use!** Even though codes 1001-1011 are defined in the WebSocket spec, browsers only allow:
- **1000** (Normal Closure)
- **3000-4999** (Custom application codes)

Attempting to use other codes (like 1001) results in:
```
InvalidAccessError: The close code must be either 1000, or between 3000 and 4999.
```

Code 4000 is in the allowed custom range and our `onclose` handler handles it:

```javascript
if (event.code === 4000) {
    addMessage('🧪 Simulated disconnect - reconnecting...', 'info');
    this.scheduleReconnect();  // ✅ Reconnection triggered!
}
```

## Allowed Close Codes

You **CAN** manually set these codes:

| Code | Name | When to Use |
|------|------|------------|
| 1000 | Normal Closure | Clean disconnect |
| 1001 | Going Away | Simulating server shutdown |
| 1002 | Protocol Error | Invalid WebSocket message |
| 1003 | Unsupported Data | Can't handle data type |
| 1007 | Invalid Frame Payload | Bad UTF-8 or similar |
| 1008 | Policy Violation | Message violates policy |
| 1009 | Message Too Big | Message exceeded limit |
| 1010 | Mandatory Extension | Client needs extension |
| 1011 | Internal Error | Server error |
| 3000-3999 | Custom | Your application codes |
| 4000-4999 | Custom | Your application codes |

You **CANNOT** manually set these:

| Code | Name | Why |
|------|------|-----|
| 1004 | Reserved | Do not use |
| 1005 | No Status Received | Set by browser only |
| 1006 | Abnormal Closure | **Set by browser only** |
| 1015 | TLS Handshake | Set by browser only |

## Testing

After the fix, the button should now:

1. Show message: "🧪 Simulating network disconnect..."
2. Close connection with code 1001
3. Show: "🔄 Server going away - will reconnect shortly"
4. Start reconnection with exponential backoff
5. Reconnect to server
6. Flush any queued messages

## To Test

1. **Restart your server** (if it's running):
   ```bash
   # Stop with Ctrl+C, then restart:
   go run main.go
   ```

2. **Refresh the browser** (important - reload the page to get the fixed JavaScript)

3. Click **"Connect"**

4. Click **"Simulate Network Issue"**

5. You should now see:
   ```
   🧪 Simulating network disconnect...
   🔌 Disconnected (code: 1001, reason: Simulated network failure)
   🔄 Server going away - will reconnect shortly
   ⏳ Reconnecting in 1s (attempt 1/10)
   ✅ Connected successfully!
   ```

## Alternative Ways to Test Real Network Failures

If you want to test **actual** code 1006 behavior:

### Option 1: Kill the server while connected
```bash
# In server terminal: Ctrl+C (this triggers graceful shutdown with 1001)
# Or: kill -9 <PID> (this causes abnormal closure - real 1006!)
```

### Option 2: Use browser DevTools
```javascript
// In browser console, break the connection at TCP level:
// (This is more advanced)
```

### Option 3: Close browser tab
Just close the browser tab - the server will detect it after the ping timeout (60 seconds).

### Option 4: Disconnect WiFi
Turn off your WiFi/network while connected - this causes real abnormal closure!

## Key Takeaway

**Code 1006 cannot be simulated programmatically** - it only happens with real network failures. For testing, use code 1001 which has similar reconnection behavior but is allowed to be set manually.

---

**Fixed**: November 2, 2025
