# REST APIs in Go — Complete Notes & Practice

---

## 1. Introduction to REST APIs

**REST** (Representational State Transfer) is an architectural style for designing networked APIs, built on top of HTTP.

### Core REST Principles
- **Resources** — everything is a resource, identified by a URL (e.g., `/users/42`).
- **HTTP methods map to actions**:
  | Method | Purpose | Idempotent? |
  |---|---|---|
  | `GET` | Read a resource | Yes |
  | `POST` | Create a resource | No |
  | `PUT` | Replace a resource entirely | Yes |
  | `PATCH` | Partially update a resource | No (usually) |
  | `DELETE` | Remove a resource | Yes |
- **Stateless** — each request contains all information needed; the server holds no client session state between requests.
- **Representations** — resources are typically exchanged as JSON (also XML, but JSON dominates modern APIs).
- **Status codes communicate outcome** — see section 8.

### Example REST resource design
```
GET    /users          -> list users
GET    /users/42       -> get user 42
POST   /users          -> create a user
PUT    /users/42       -> replace user 42
PATCH  /users/42       -> partially update user 42
DELETE /users/42       -> delete user 42
```

### Why Go for REST APIs
- Excellent standard library (`net/http`) — you can build production APIs with **zero external dependencies**.
- Fast compilation, static binaries, great concurrency for handling many simultaneous requests.
- Strong typing pairs naturally with JSON (de)serialization via struct tags.

---

## 2. Building APIs using net/http

Go's standard library `net/http` package is enough to build a full REST API without any framework.

### 2.1 Minimal server
```go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, REST API!")
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### 2.2 The Handler interface
Anything satisfying this interface can handle HTTP requests:
```go
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
```
`http.HandleFunc` is a convenience wrapper — it adapts a plain function `func(w http.ResponseWriter, r *http.Request)` into a `Handler`.

### 2.3 ServeMux — Go's built-in router
```go
mux := http.NewServeMux()
mux.HandleFunc("/users", usersHandler)
mux.HandleFunc("/users/", userByIDHandler) // trailing slash = prefix match

server := &http.Server{
	Addr:         ":8080",
	Handler:      mux,
	ReadTimeout:  5 * time.Second,
	WriteTimeout: 10 * time.Second,
	IdleTimeout:  120 * time.Second,
}
log.Fatal(server.ListenAndServe())
```
> **Note:** since **Go 1.22**, `http.ServeMux` supports method matching and path parameters natively:
```go
mux := http.NewServeMux()
mux.HandleFunc("GET /users/{id}", getUserHandler)
mux.HandleFunc("POST /users", createUserHandler)
mux.HandleFunc("DELETE /users/{id}", deleteUserHandler)

// inside a handler:
id := r.PathValue("id")
```
This removed the need for a third-party router in many simple projects. For complex routing needs, libraries like `gorilla/mux`, `chi`, or `gin` are still popular.

### 2.4 Always set a request timeout / use `http.Server`
Don't use the bare `http.ListenAndServe(":8080", nil)` in production — it uses `http.DefaultServeMux` and has no timeouts (risk of resource exhaustion from slow clients). Always configure an explicit `http.Server`.

---

## 3. Routing and Middleware

### 3.1 Routing with `chi` (popular lightweight router)
```go
import "github.com/go-chi/chi/v5"

r := chi.NewRouter()
r.Get("/users", listUsers)
r.Get("/users/{id}", getUser)
r.Post("/users", createUser)
r.Put("/users/{id}", updateUser)
r.Delete("/users/{id}", deleteUser)

http.ListenAndServe(":8080", r)

// inside handler:
id := chi.URLParam(r, "id")
```

### 3.2 What is Middleware?
Middleware wraps a handler to add cross-cutting behavior (logging, auth, recovery, CORS) without modifying the handler itself. It's just a function that takes a `Handler` and returns a `Handler`.

```go
type Middleware func(http.Handler) http.Handler
```

### 3.3 Writing your own middleware
```go
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r) // call the next handler in the chain
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic recovered: %v", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
```

### 3.4 Chaining middleware
```go
func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

// usage
finalHandler := Chain(mux, LoggingMiddleware, RecoveryMiddleware, AuthMiddleware)
http.ListenAndServe(":8080", finalHandler)
```
With `chi`, this is built in:
```go
r.Use(middleware.Logger)
r.Use(middleware.Recoverer)
r.Use(AuthMiddleware)
```

### 3.5 Common middleware types
- **Logging** — record method, path, status, latency.
- **Recovery** — catch panics so one bad request doesn't crash the server.
- **CORS** — allow cross-origin browser requests.
- **Authentication** — verify tokens/sessions before letting requests through.
- **Rate limiting** — throttle abusive clients.
- **Request ID / tracing** — attach a unique ID to each request for log correlation.

---

## 4. CRUD Operations

Full example: an in-memory "users" API (`Create, Read, Update, Delete`).

```go
package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserStore struct {
	mu     sync.Mutex
	users  map[int]User
	nextID int
}

func NewUserStore() *UserStore {
	return &UserStore{users: make(map[int]User), nextID: 1}
}

// CREATE
func (s *UserStore) Create(u User) User {
	s.mu.Lock()
	defer s.mu.Unlock()
	u.ID = s.nextID
	s.users[u.ID] = u
	s.nextID++
	return u
}

// READ (single)
func (s *UserStore) Get(id int) (User, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, ok := s.users[id]
	return u, ok
}

// READ (all)
func (s *UserStore) List() []User {
	s.mu.Lock()
	defer s.mu.Unlock()
	list := make([]User, 0, len(s.users))
	for _, u := range s.users {
		list = append(list, u)
	}
	return list
}

// UPDATE
func (s *UserStore) Update(id int, u User) (User, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.users[id]; !ok {
		return User{}, false
	}
	u.ID = id
	s.users[id] = u
	return u, true
}

// DELETE
func (s *UserStore) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.users[id]; !ok {
		return false
	}
	delete(s.users, id)
	return true
}

var store = NewUserStore()

func usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		json.NewEncoder(w).Encode(store.List())
	case http.MethodPost:
		var u User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "invalid JSON body", http.StatusBadRequest)
			return
		}
		created := store.Create(u)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(created)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func userByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/users/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		u, ok := store.Get(id)
		if !ok {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(u)

	case http.MethodPut:
		var u User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "invalid JSON body", http.StatusBadRequest)
			return
		}
		updated, ok := store.Update(id, u)
		if !ok {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(updated)

	case http.MethodDelete:
		if !store.Delete(id) {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/users", usersHandler)
	http.HandleFunc("/users/", userByIDHandler)
	http.ListenAndServe(":8080", nil)
}
```

### Testing with curl
```bash
curl -X POST localhost:8080/users -d '{"name":"Alice","email":"alice@test.com"}'
curl localhost:8080/users
curl localhost:8080/users/1
curl -X PUT localhost:8080/users/1 -d '{"name":"Alice Updated","email":"a2@test.com"}'
curl -X DELETE localhost:8080/users/1
```

---

## 5. Request and Response Handling

### 5.1 Reading request data
```go
// Query params
name := r.URL.Query().Get("name")

// Path params (Go 1.22+ ServeMux)
id := r.PathValue("id")

// Headers
token := r.Header.Get("Authorization")

// Body (JSON)
var input SomeStruct
if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
	http.Error(w, "invalid body", http.StatusBadRequest)
	return
}
defer r.Body.Close()

// Form data (application/x-www-form-urlencoded)
r.ParseForm()
value := r.FormValue("key")
```

### 5.2 Limiting body size (protect against huge payloads)
```go
r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1MB limit
```

### 5.3 Writing responses
```go
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOK) // must be called BEFORE writing body, and only once
json.NewEncoder(w).Encode(response)
```
> **Gotcha:** Once you call `w.Write()` (or `json.NewEncoder(w).Encode()`), the status code defaults to `200` if you haven't explicitly set one. `WriteHeader()` must be called before any write, or it's ignored.

### 5.4 Helper functions (common pattern in real projects)
```go
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// usage
writeJSON(w, http.StatusOK, user)
writeError(w, http.StatusNotFound, "user not found")
```

### 5.5 Context in requests (timeouts, cancellation, request-scoped values)
```go
func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context() // cancelled automatically if client disconnects
	result, err := doSlowDBQuery(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "query failed")
		return
	}
	writeJSON(w, http.StatusOK, result)
}
```

---

## 6. JSON APIs in Go

### 6.1 Struct tags control JSON field names
```go
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email,omitempty"`     // omit if empty
	Password  string    `json:"-"`                    // never serialize
	CreatedAt time.Time `json:"created_at"`
}
```

### 6.2 Encoding (struct → JSON)
```go
data, err := json.Marshal(user)               // []byte
json.NewEncoder(w).Encode(user)               // write directly to a writer (streaming, preferred for HTTP responses)
prettyData, _ := json.MarshalIndent(user, "", "  ") // pretty-printed
```

### 6.3 Decoding (JSON → struct)
```go
var user User
err := json.Unmarshal(data, &user)             // from []byte
err := json.NewDecoder(r.Body).Decode(&user)   // from a reader (preferred for HTTP requests)
```

### 6.4 Reject unknown fields (stricter validation)
```go
dec := json.NewDecoder(r.Body)
dec.DisallowUnknownFields()
if err := dec.Decode(&user); err != nil {
	http.Error(w, "unexpected field in request", http.StatusBadRequest)
	return
}
```

### 6.5 Handling nested / dynamic JSON
```go
type Address struct {
	City    string `json:"city"`
	Country string `json:"country"`
}

type User struct {
	Name    string  `json:"name"`
	Address Address `json:"address"`
}
```
For truly dynamic/unknown structure: `map[string]interface{}` or `json.RawMessage` (defer parsing part of the payload).

### 6.6 Custom marshal/unmarshal (e.g. custom date format)
```go
type CustomDate struct {
	time.Time
}

func (d CustomDate) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.Format("2006-01-02") + `"`), nil
}

func (d *CustomDate) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}
```

---

## 7. Authentication Basics

### 7.1 API Keys (simplest)
```go
func APIKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-API-Key")
		if key != "expected-secret-key" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
```

### 7.2 Basic Auth
```go
func BasicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != "admin" || pass != "secret" {
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted"`)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
```

### 7.3 JWT (JSON Web Tokens) — most common for real APIs
```go
import "github.com/golang-jwt/jwt/v5"

var jwtSecret = []byte("your-secret-key")

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userID int) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader { // no "Bearer " prefix found
			http.Error(w, "missing bearer token", http.StatusUnauthorized)
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		// attach user ID to context for downstream handlers
		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
```
Retrieving the user ID in a handler:
```go
userID := r.Context().Value("userID").(int)
```
> Best practice: use a custom typed key instead of a raw string for `context.WithValue` to avoid collisions:
```go
type ctxKey string
const userIDKey ctxKey = "userID"
```

### 7.4 Password hashing (never store plaintext passwords)
```go
import "golang.org/x/crypto/bcrypt"

func HashPassword(pw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(hash), err
}

func CheckPassword(hash, pw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))
	return err == nil
}
```

### 7.5 Auth flow summary
```
1. POST /login {email, password}
2. Server verifies password against bcrypt hash
3. Server issues JWT, returns it to client
4. Client sends "Authorization: Bearer <token>" on subsequent requests
5. JWTMiddleware validates token on protected routes
```

---

## 8. API Error Handling and Validation

### 8.1 HTTP status codes to use
| Code | Meaning | When to use |
|---|---|---|
| `200 OK` | Success | Successful GET/PUT/PATCH |
| `201 Created` | Resource created | Successful POST |
| `204 No Content` | Success, no body | Successful DELETE |
| `400 Bad Request` | Client sent invalid data | Malformed JSON, failed validation |
| `401 Unauthorized` | Missing/invalid credentials | No token, bad token |
| `403 Forbidden` | Authenticated but not allowed | Valid user, insufficient permission |
| `404 Not Found` | Resource doesn't exist | Bad ID |
| `405 Method Not Allowed` | Wrong HTTP verb for the route | |
| `409 Conflict` | State conflict | Duplicate resource, version conflict |
| `422 Unprocessable Entity` | Semantically invalid data | Failed business-rule validation |
| `429 Too Many Requests` | Rate limited | |
| `500 Internal Server Error` | Unexpected server-side failure | |

### 8.2 Consistent error response format
```go
type ErrorResponse struct {
	Error   string            `json:"error"`
	Details map[string]string `json:"details,omitempty"`
}

func writeErrorJSON(w http.ResponseWriter, status int, msg string, details map[string]string) {
	writeJSON(w, status, ErrorResponse{Error: msg, Details: details})
}
```

### 8.3 Input validation
```go
type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func (req CreateUserRequest) Validate() map[string]string {
	errs := map[string]string{}
	if strings.TrimSpace(req.Name) == "" {
		errs["name"] = "name is required"
	}
	if !strings.Contains(req.Email, "@") {
		errs["email"] = "invalid email format"
	}
	if req.Age < 0 || req.Age > 150 {
		errs["age"] = "age must be between 0 and 150"
	}
	return errs
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorJSON(w, http.StatusBadRequest, "invalid JSON", nil)
		return
	}
	if errs := req.Validate(); len(errs) > 0 {
		writeErrorJSON(w, http.StatusUnprocessableEntity, "validation failed", errs)
		return
	}
	// proceed to create user...
}
```

### 8.4 Using a validation library (`go-playground/validator`)
```go
import "github.com/go-playground/validator/v10"

type CreateUserRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Age   int    `json:"age" validate:"gte=0,lte=150"`
}

var validate = validator.New()

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	json.NewDecoder(r.Body).Decode(&req)

	if err := validate.Struct(req); err != nil {
		writeErrorJSON(w, http.StatusUnprocessableEntity, "validation failed",
			map[string]string{"details": err.Error()})
		return
	}
	// ...
}
```

### 8.5 Centralized error handling pattern (custom handler type)
```go
type APIError struct {
	Status int
	Msg    string
}

func (e *APIError) Error() string { return e.Msg }

type APIHandler func(w http.ResponseWriter, r *http.Request) error

func Wrap(h APIHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			if apiErr, ok := err.(*APIError); ok {
				writeErrorJSON(w, apiErr.Status, apiErr.Msg, nil)
				return
			}
			log.Println("unexpected error:", err)
			writeErrorJSON(w, http.StatusInternalServerError, "internal server error", nil)
		}
	}
}

// usage — handlers just return errors, no manual response writing needed
func getUserHandler(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return &APIError{Status: http.StatusBadRequest, Msg: "invalid id"}
	}
	user, ok := store.Get(id)
	if !ok {
		return &APIError{Status: http.StatusNotFound, Msg: "user not found"}
	}
	return writeJSONErr(w, http.StatusOK, user)
}

mux.HandleFunc("GET /users/{id}", Wrap(getUserHandler))
```

### 8.6 Recovering from panics (never let one bad request crash the server)
```go
func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic: %v\n%s", rec, debug.Stack())
				writeErrorJSON(w, http.StatusInternalServerError, "internal server error", nil)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
```

---

## 9. Putting It All Together — Minimal Production-Style API Skeleton

```go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("GET /users", listUsers)
	mux.HandleFunc("POST /users", createUser)
	mux.HandleFunc("GET /users/{id}", getUser)
	mux.HandleFunc("PUT /users/{id}", updateUser)
	mux.HandleFunc("DELETE /users/{id}", deleteUser)

	var handler http.Handler = mux
	handler = LoggingMiddleware(handler)
	handler = RecoverMiddleware(handler)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("listening on :8080")
	log.Fatal(server.ListenAndServe())
}
```

---

## 10. Practice Questions & Answers

**Q1. Why should you avoid `http.ListenAndServe(":8080", nil)` in production?**
**A:** It uses `http.DefaultServeMux` (a global, shared mux any imported package could register routes onto — a security/maintainability risk) and has no read/write timeouts configured, leaving the server vulnerable to slow-client resource exhaustion. Always construct an explicit `http.Server{}` with timeouts and your own mux.

**Q2. What's the difference between `PUT` and `PATCH`?**
**A:** `PUT` replaces the entire resource (client sends the full representation); `PATCH` applies a partial update (only the fields provided are changed).

**Q3. Why use `json.NewDecoder(r.Body).Decode(&v)` instead of `io.ReadAll` + `json.Unmarshal`?**
**A:** The decoder streams directly from the request body without buffering the whole payload into memory first, which is more efficient for large payloads and integrates cleanly with `io.Reader`.

**Q4. What does middleware actually do, structurally, in Go?**
**A:** A middleware is a function `func(http.Handler) http.Handler` — it wraps an existing handler and returns a new one that runs some logic before/after calling the wrapped handler's `ServeHTTP`. Chaining multiple middlewares nests these wrappers.

**Q5. Why is `context.WithValue` risky for storing auth data, and how do you mitigate it?**
**A:** Using plain string keys risks collisions between packages using the same key name. Mitigate by defining an unexported custom type for context keys (e.g., `type ctxKey string; const userIDKey ctxKey = "userID"`), which guarantees uniqueness via the type system.

**Q6. What status code should a successful DELETE return, and why no body?**
**A:** `204 No Content` — the operation succeeded but there's nothing meaningful to return, so the body is empty by convention.

**Q7. How do you prevent a single panicking request handler from crashing the whole server?**
**A:** Wrap handlers with a recovery middleware that uses `defer` + `recover()` to catch panics, log them, and respond with a `500` instead of letting the panic propagate and (in default `net/http` behavior) just close that one request's goroutine — note Go's `net/http` actually already recovers panics per-request by default and logs them, but writing your own gives you control over the JSON error response format.

---

*End of notes — build the CRUD example above, add JWT auth to protect the `POST/PUT/DELETE` routes, and add validation with `go-playground/validator` to solidify all these pieces together.*
