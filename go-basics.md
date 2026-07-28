# Go (Golang) Basics

## 1. Hello World

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

- Every Go file belongs to a `package`.
- `package main` + `func main()` = executable program.
- `import` brings in standard library or external packages.

---

## 2. Variables

```go
var x int = 10
var y = 20        // type inferred
z := 30           // short declaration (inside functions only)

var a, b, c = 1, 2, 3
const Pi = 3.14
```

- `:=` can only be used inside functions.
- Unused variables cause a compile error.

---

## 3. Basic Types

| Type      | Example              |
|-----------|----------------------|
| int       | `var i int = 5`      |
| float64   | `var f float64 = 3.2`|
| string    | `var s string = "hi"`|
| bool      | `var b bool = true`  |
| byte      | alias for `uint8`    |
| rune      | alias for `int32`    |

Zero values: `0`, `0.0`, `""`, `false`.

---

## 4. Functions

```go
func add(a int, b int) int {
    return a + b
}

// shorthand for same-type params
func add2(a, b int) int {
    return a + b
}

// multiple return values
func divmod(a, b int) (int, int) {
    return a / b, a % b
}

// named return values
func split(sum int) (x, y int) {
    x = sum * 4 / 9
    y = sum - x
    return
}
```

---

## 5. Control Flow

### If / Else
```go
if x > 0 {
    fmt.Println("positive")
} else if x == 0 {
    fmt.Println("zero")
} else {
    fmt.Println("negative")
}
```

### For (Go's only loop)
```go
for i := 0; i < 10; i++ {
    fmt.Println(i)
}

// while-style
for x < 100 {
    x *= 2
}

// infinite loop
for {
    break
}
```

### Switch
```go
switch day {
case "Mon", "Tue":
    fmt.Println("early week")
case "Fri":
    fmt.Println("almost weekend")
default:
    fmt.Println("other day")
}
```

---

## 6. Arrays, Slices, Maps

```go
// Array (fixed size)
var arr [5]int

// Slice (dynamic array)
nums := []int{1, 2, 3}
nums = append(nums, 4)

// Slicing
sub := nums[1:3]

// Map
m := map[string]int{"a": 1, "b": 2}
m["c"] = 3
val, ok := m["a"]   // ok = true if key exists
delete(m, "a")
```

---

## 7. Structs

```go
type Person struct {
    Name string
    Age  int
}

p := Person{Name: "Alice", Age: 30}
p.Age = 31

// Pointer to struct
func birthday(p *Person) {
    p.Age++
}
birthday(&p)
```

---

## 8. Methods & Interfaces

```go
func (p Person) Greet() string {
    return "Hi, I'm " + p.Name
}

type Greeter interface {
    Greet() string
}
```

- Methods are functions with a receiver (`(p Person)`).
- Use pointer receiver `(p *Person)` to modify the original value.
- Interfaces are satisfied implicitly (no `implements` keyword).

---

## 9. Error Handling

```go
func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, fmt.Errorf("cannot divide by zero")
    }
    return a / b, nil
}

result, err := divide(10, 0)
if err != nil {
    fmt.Println("Error:", err)
}
```

- No exceptions; errors are returned as normal values.
- `panic`/`recover` exist for exceptional/unrecoverable situations.

---

## 10. Goroutines & Channels

```go
func say(s string) {
    fmt.Println(s)
}

go say("hello")   // runs concurrently

ch := make(chan int)

go func() {
    ch <- 42       // send value
}()

val := <-ch        // receive value
```

- `go` keyword starts a goroutine (lightweight thread).
- Channels are used to communicate/synchronize between goroutines.
- `select` lets you wait on multiple channel operations.

```go
select {
case msg1 := <-ch1:
    fmt.Println(msg1)
case msg2 := <-ch2:
    fmt.Println(msg2)
default:
    fmt.Println("no message")
}
```

---

## 11. Packages & Modules

```bash
go mod init myproject   # start a new module
go build                # compile
go run main.go          # compile + run
go get <package>        # add dependency
go test                 # run tests
```

---

## 12. Quick Reference

- `fmt.Println`, `fmt.Printf`, `fmt.Sprintf` — output/formatting
- `len()` — length of string, slice, array, map, channel
- `make()` — create slices, maps, channels
- `new()` — allocate zeroed memory, return pointer
- `defer` — schedule a call to run after function returns
- `iota` — auto-incrementing constant generator

```go
const (
    A = iota // 0
    B        // 1
    C        // 2
)
```
