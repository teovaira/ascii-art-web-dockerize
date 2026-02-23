# Middleware Stack - Complete Documentation

## Overview

The server uses a layered middleware architecture that processes every HTTP request through multiple security and monitoring layers.

## Middleware Chain

```
Request → Security → Recovery → Logging → Routes → Response
```

### Execution Order

**Incoming Request:**
1. Security Headers - Adds security headers
2. Recovery - Sets up panic handler
3. Logging - Records start time
4. Routes - Executes handler
5. Logging - Records duration
6. Recovery - Catches any panics
7. Security - Headers already set
8. Response sent to client

## Individual Middleware

### 1. Security Headers (Outermost)
**File:** `internal/middleware/security.go`

**Purpose:** Protect against common web attacks

**Headers Set:**
- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`
- `X-XSS-Protection: 1; mode=block`
- `Content-Security-Policy: default-src 'self'`

**Example:**
```go
handler := middleware.SecurityHeaders(nextHandler)
```

---

### 2. Recovery (Middle Layer)
**File:** `internal/middleware/recovery.go`

**Purpose:** Catch panics and prevent server crashes

**Behavior:**
- Catches all panics
- Logs panic details
- Returns 500 Internal Server Error
- Server continues running

**Example:**
```go
handler := middleware.Recovery(logger)(nextHandler)
```

---

### 3. Logging (Inner Layer)
**File:** `internal/middleware/logging.go`

**Purpose:** Track all HTTP requests

**Log Format:**
```
[2026-02-23 20:51:19] GET / 200 0ms
[YYYY-MM-DD HH:MM:SS] METHOD PATH STATUS DURATIONms
```

**Example:**
```go
handler := middleware.Logging(logger)(nextHandler)
```

---

## Complete Integration

**In `internal/server/server.go`:**
```go
logger := log.New(os.Stdout, "", log.LstdFlags)
handler := middleware.SecurityHeaders(
    middleware.Recovery(logger)(
        middleware.Logging(logger)(mux)))
```

## Test Coverage

| Package | Coverage | Tests |
|---------|----------|-------|
| middleware | 100.0% | 24 tests |
| server | 90.0% | 5 tests |

## Benefits

### Security
- ✅ Protection against XSS attacks
- ✅ Protection against clickjacking
- ✅ Protection against MIME sniffing
- ✅ Content Security Policy enforcement

### Reliability
- ✅ Server never crashes from panics
- ✅ Graceful error handling
- ✅ All errors logged

### Monitoring
- ✅ Complete request tracking
- ✅ Performance monitoring (duration)
- ✅ Status code tracking
- ✅ Timestamp for every request

### Maintainability
- ✅ 100% test coverage
- ✅ Table-driven tests
- ✅ Easy to add new middleware
- ✅ Clear separation of concerns

## Adding New Middleware

To add new middleware to the chain:

1. Create `internal/middleware/newmiddleware.go`
2. Create `internal/middleware/newmiddleware_test.go`
3. Update `internal/server/server.go`:
   ```go
   handler := middleware.NewMiddleware(
       middleware.SecurityHeaders(
           middleware.Recovery(logger)(
               middleware.Logging(logger)(mux))))
   ```

**Remember:** Order matters! Outermost middleware executes first.

## Team Integration

### For Vasiliki (Frontend/Templates)
- All your template routes automatically get:
  - Security headers
  - Panic recovery
  - Request logging

### For Theo (Integration Tests)
- All endpoints return proper status codes
- Security headers present in all responses
- Server never crashes during tests

### For You (Backend)
- Add new routes to `mux`
- Middleware automatically applies
- No additional configuration needed
