# Issue 1.2: Global Redis Singleton Anti-Pattern - Detailed Analysis

## What is a Singleton Pattern?

A singleton pattern ensures that only **one instance** of a class/object exists throughout the application's lifetime. In this codebase, it's being used to ensure only one Redis client connection exists.

## Current Implementation

```go
// internal/database/redis.go

var redisClient *redis.Client  // ← Global variable (package-level)

func ConnectRedisClient(cfg *config.RedisConfig) *redis.Client {
	if redisClient != nil {  // ← Check if already exists
		return redisClient   // ← Return existing instance
	}
	// ← Create new instance only if nil
	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	// ... ping and validation ...
	return redisClient
}

func GetRedisClient() *redis.Client {
	if redisClient == nil {
		log.Fatalln("Redis client not initialized. Call ConnectRedis() first")
		return nil
	}
	return redisClient
}
```

## Why This Pattern Was Likely Chosen

**Good intentions:**
- ✅ Ensure only one Redis connection pool exists (efficient)
- ✅ Avoid creating multiple connections (resource management)
- ✅ Easy access from anywhere in the codebase

**However, these intentions are flawed in Go's concurrent environment.**

---

## Problem 1: Race Condition (CRITICAL) 🔴

### The Problem

Go is a **concurrent language**. Multiple goroutines can call `ConnectRedisClient()` simultaneously. The current code has a **data race**:

```go
func ConnectRedisClient(cfg *config.RedisConfig) *redis.Client {
	if redisClient != nil {  // ← Thread 1 checks: nil
		return redisClient
	}
	// ← Thread 2 also checks: nil (at the same time!)
	redisClient = redis.NewClient(...)  // ← Both threads create clients!
	// ...
}
```

### Race Condition Scenario

**Timeline:**
```
Time    Thread 1                    Thread 2
─────────────────────────────────────────────────
T1      if redisClient != nil       (waiting)
T2      (checking...)               if redisClient != nil
T3      redisClient = NewClient()   redisClient = NewClient()
T4      return redisClient          return redisClient
```

**Result:** Two Redis clients are created! The second one overwrites the first, potentially causing:
- Memory leak (first client not closed)
- Connection leak (first client's connections not released)
- Unpredictable behavior (which client is actually used?)

### Proof of Race Condition

You can detect this with Go's race detector:

```bash
go run -race cmd/server/main.go
```

This will likely show:
```
WARNING: DATA RACE
Read at 0x... by goroutine X:
  internal/database/redis.go:15: if redisClient != nil

Write at 0x... by goroutine Y:
  internal/database/redis.go:18: redisClient = redis.NewClient(...)
```

### Why This Matters

In production with high concurrency:
- Multiple HTTP requests arrive simultaneously
- Each might trigger `ConnectRedisClient()` 
- Race condition occurs → multiple clients created
- Resource exhaustion → application crashes

---

## Problem 2: Testing Nightmare 🟡

### The Problem

Global state makes unit testing **extremely difficult**:

```go
// Test 1
func TestSomething(t *testing.T) {
    // Setup: ConnectRedisClient() called
    redisClient = redis.NewClient(...)  // Global state modified
    
    // Test runs...
    
    // Cleanup: How do we reset?
    redisClient = nil  // ← Breaks other tests running in parallel!
}
```

### Issues:

1. **Tests can't run in parallel** - They share global state
2. **Tests interfere with each other** - One test's setup affects another
3. **Can't mock Redis** - Hard to inject test doubles
4. **Cleanup is impossible** - Can't reset state without affecting other tests

### Example Test Failure

```go
// Test 1 (runs first)
func TestAuthService_Login(t *testing.T) {
    cfg := &config.RedisConfig{Host: "localhost", Port: 6379}
    client := database.ConnectRedisClient(cfg)  // Creates real client
    service := auth.NewService(..., client)
    // Test passes
}

// Test 2 (runs in parallel - FAILS!)
func TestAuthService_Logout(t *testing.T) {
    // Wants to use mock Redis, but...
    mockRedis := &MockRedisClient{}
    // Can't inject it because ConnectRedisClient() returns global instance!
    // Test fails or uses wrong client
}
```

---

## Problem 3: No Lifecycle Management 🟡

### The Problem

The singleton pattern makes it **impossible to properly manage the Redis connection lifecycle**:

```go
// How do you:
// 1. Close the connection gracefully?
// 2. Reconnect on failure?
// 3. Handle connection timeouts?
// 4. Shutdown cleanly?

// Current code: NONE of these are possible!
```

### Missing Features:

1. **Graceful Shutdown** - Can't close connection on app exit
2. **Connection Pooling** - No control over pool size
3. **Health Checks** - No periodic health checks
4. **Reconnection Logic** - If connection drops, how to reconnect?
5. **Connection Limits** - No way to limit concurrent connections

### Example: Graceful Shutdown Failure

```go
// main.go
func main() {
    redisClient := database.ConnectRedisClient(cfg)
    // ... server runs ...
    
    // On shutdown signal (SIGTERM):
    // How do we close redisClient?
    // It's a global variable, not accessible here!
    // Result: Connection leaks, Redis server sees hanging connections
}
```

---

## Problem 4: Hidden Dependencies 🟡

### The Problem

Global variables create **hidden dependencies** that are hard to track:

```go
// Somewhere deep in your code:
func someFunction() {
    client := database.GetRedisClient()  // ← Hidden dependency!
    // Where did this come from? When was it initialized?
    // What if it's nil? What if it's the wrong client?
}
```

### Issues:

1. **Dependency Graph is Unclear** - Hard to see what depends on Redis
2. **Initialization Order Matters** - Must call `ConnectRedisClient()` before `GetRedisClient()`
3. **No Compile-Time Safety** - Compiler can't catch missing initialization
4. **Hard to Refactor** - Can't easily swap implementations

---

## Problem 5: Configuration Can't Change 🟡

### The Problem

Once initialized, the Redis client is **locked to the first configuration**:

```go
// First call (in production):
redisClient := ConnectRedisClient(&config.RedisConfig{
    Host: "prod-redis.example.com",
    Port: 6379,
})

// Later, test wants different config:
redisClient := ConnectRedisClient(&config.RedisConfig{
    Host: "localhost",
    Port: 6380,  // ← IGNORED! Returns first client
})
```

**Result:** Tests can't use different Redis instances, can't test with different configurations.

---

## Problem 6: No Error Recovery 🟡

### The Problem

If the Redis connection fails after initialization, there's **no way to recover**:

```go
// Initial connection succeeds
redisClient := ConnectRedisClient(cfg)  // ✅ Connected

// Later, network issue:
// Redis connection drops
// Application continues using broken client
// All Redis operations fail silently or crash

// How to reconnect? Can't! redisClient is already set.
```

---

## The Correct Solution: Dependency Injection

### Principle

Instead of global state, **pass the Redis client as a dependency** through function parameters and struct fields.

### Refactored Code

```go
// internal/database/redis.go

// Remove global variable entirely!
// var redisClient *redis.Client  ← DELETE THIS

// Simple factory function (no singleton)
func NewRedisClient(cfg *config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     10,              // ← Configurable pool size
		MinIdleConns: 5,               // ← Minimum idle connections
		MaxRetries:   3,               // ← Retry configuration
		DialTimeout:  5 * time.Second, // ← Connection timeout
		ReadTimeout:  3 * time.Second, // ← Read timeout
		WriteTimeout: 3 * time.Second, // ← Write timeout
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return client, nil
}

// Optional: Helper to close client gracefully
func CloseRedisClient(client *redis.Client) error {
	if client != nil {
		return client.Close()
	}
	return nil
}
```

### Updated Router Setup

```go
// internal/router/router.go

func SetupRouter(cfg *config.Config) *gin.Engine {
	// Infrastructure layer - Database
	db := database.Connect(&cfg.Database)
	
	// Infrastructure layer - Redis (no singleton!)
	redisClient, err := database.NewRedisClient(&cfg.Redis)
	if err != nil {
		log.Fatalln("Failed to connect to Redis:", err)
	}
	defer func() {
		// Note: In real app, close on graceful shutdown, not here
		// This is just for demonstration
	}()

	// Data layer - Repositories
	userRepo := user.NewRepository(db)

	// Domain layer - Services (Redis passed as dependency)
	userService := user.NewService(userRepo)
	authService := auth.NewService(userService, cfg.Google, redisClient)

	// Presentation layer - Handlers
	userHandler := user.NewHandler(userService, cfg)
	authHandler := auth.NewHandler(authService, cfg)

	// Router setup
	router := gin.Default()
	router.Use(utils.CORS(&cfg.Server.Cors))

	// Register domain routes
	user.RegisterRoutes(router, userHandler)
	auth.RegisterRoutes(router, authHandler)

	return router
}
```

### Benefits of This Approach

1. ✅ **No Race Conditions** - Each call creates a new client (or use sync.Once if you really need singleton)
2. ✅ **Testable** - Easy to inject mock Redis client
3. ✅ **Explicit Dependencies** - Clear what depends on Redis
4. ✅ **Lifecycle Management** - Can close/reconnect as needed
5. ✅ **Flexible Configuration** - Different configs for different environments
6. ✅ **Compile-Time Safety** - Compiler ensures dependencies are provided

---

## If You Really Need a Singleton (Not Recommended)

If you absolutely must have a singleton (rarely needed), use `sync.Once`:

```go
var (
	redisClient *redis.Client
	redisOnce   sync.Once  // ← Thread-safe initialization
)

func ConnectRedisClient(cfg *config.RedisConfig) *redis.Client {
	redisOnce.Do(func() {
		client, err := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Password: cfg.Password,
			DB:       cfg.DB,
		})
		if err != nil {
			log.Fatalln("failed to create Redis client:", err)
		}
		
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		if err := client.Ping(ctx).Err(); err != nil {
			log.Fatalln("failed to connect to Redis:", err)
		}
		
		redisClient = client
	})
	return redisClient
}
```

**But this still has testing and lifecycle issues!** Dependency injection is better.

---

## Testing Example with Dependency Injection

### Before (With Singleton - Hard to Test)

```go
func TestAuthService_Logout(t *testing.T) {
    // Problem: Can't inject mock Redis
    // Must use real Redis or hack global variable
    database.redisClient = &MockRedisClient{}  // ← HACK!
    
    service := auth.NewService(...)
    // Test...
    
    // Problem: How to cleanup?
    database.redisClient = nil  // ← Breaks other tests!
}
```

### After (With Dependency Injection - Easy to Test)

```go
func TestAuthService_Logout(t *testing.T) {
    // Easy: Create mock and inject
    mockRedis := &MockRedisClient{}
    mockRedis.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
        Return(nil)
    
    service := auth.NewService(userService, cfg.Google, mockRedis)
    
    // Test runs...
    // No cleanup needed - mock is isolated
}
```

---

## Migration Steps

1. **Remove global variable** from `internal/database/redis.go`
2. **Rename `ConnectRedisClient`** to `NewRedisClient` and return `(*redis.Client, error)`
3. **Update router setup** to handle the error and pass client as dependency
4. **Remove `GetRedisClient()`** function (no longer needed)
5. **Update all usages** to receive Redis client as parameter
6. **Add graceful shutdown** to close Redis connection properly

---

## Summary

| Problem | Impact | Severity |
|---------|--------|----------|
| Race Condition | Multiple clients created, resource leaks | 🔴 CRITICAL |
| Testing Issues | Can't test properly, tests interfere | 🟡 HIGH |
| No Lifecycle Management | Can't shutdown gracefully | 🟡 HIGH |
| Hidden Dependencies | Hard to maintain, refactor | 🟡 MEDIUM |
| Configuration Locked | Can't change config after init | 🟡 MEDIUM |
| No Error Recovery | Broken connections persist | 🟡 MEDIUM |

**Solution:** Use dependency injection instead of singleton pattern. It's the idiomatic Go way and solves all these problems.

