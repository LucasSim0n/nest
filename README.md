# Cafe ☕

A minimalist HTTP micro-framework for Go, inspired by *Nest / Express*-style composition with routers and middlewares, but built **entirely on top of `net/http`**.

Cafe is designed to:

* Keep the core small and easy to understand
* Avoid external dependencies
* Compose HTTP applications using **mounted routers**, **middlewares**, and **clear HTTP methods**

---

## ✨ Features

* ✅ Simple API (`Get`, `Post`, `Put`, `Delete`)
* 🧩 Nested routers (routers inside routers)
* 🧠 Chainable middlewares
* 🔒 Prevents duplicate routes
* ⚙️ Based on the standard `http.ServeMux`
* 📦 Zero external dependencies

---

## 📦 Installation

```bash
go get github.com/LucasSim0n/cafe
```

---

## 🚀 Basic usage

```go
package main

import (
    "net/http"
    "github.com/LucasSim0n/cafe"
)

func main() {
    app := cafe.NewServer()

    app.Get("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello world"))
    })

    app.Listen(":3000")
}
```

---

## 🛣️ HTTP Routes

> **Path parameters support**
>
> Cafe inherits native path parameter support from **`net/http` (Go 1.22+)**.
> This means you can use parameterized paths directly, following the standard library syntax:
>
> ```go
> app.Get("/users/{id}", handler)
> app.Get("/posts/{slug}", handler)
> ```
>
> Parameters can be accessed from the request using `r.PathValue("param")`:
>
> ```go
> id := r.PathValue("id")
> ```
>
> No custom router or parameter parser is implemented in Cafe — it intentionally relies on the behavior and guarantees of `net/http`.

Both the server and routers support the basic HTTP methods:

```go
app.Get(path, handler)
app.Post(path, handler)
app.Put(path, handler)
app.Delete(path, handler)
```

Duplicate routes (same method + path) are automatically ignored.

---

## 🧩 Routers

You can group routes using routers and mount them under a base path.

```go
api := cafe.NewRouter()

api.Get("/users", usersHandler)
api.Post("/users", createUserHandler)

app.UseRouter("/api", api)
```

This produces:

```
GET  /api/users
POST /api/users
```

---

## 🌳 Nested routers

Routers can also contain other routers:

```go
admin := cafe.NewRouter()
admin.Get("/dashboard", dashboardHandler)

api.UseRouter("/admin", admin)
```

Final result:

```
GET /api/admin/dashboard
```

---

## 🧠 Middlewares

A middleware is a function that wraps an `http.HandlerFunc`:

```go
type middleware func(next http.HandlerFunc) http.HandlerFunc
```

### Global middleware

```go
app.Use(func(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        log.Println(r.Method, r.URL.Path)
        next(w, r)
    }
})
```

Applied to **all** server routes.

---

### Router-level middleware

```go
api.Use(authMiddleware)
```

Applied only to that router’s routes (and its child routers).

---

### Execution order

Middlewares run in declaration order:

```go
app.Use(mw1)
app.Use(mw2)
```

Execution flow:

```
mw1 → mw2 → handler
```

---

## 🔧 Internals (brief)

* Uses patterns like:

  ```
  METHOD /path/{$}
  ```

  to simulate method-based routing using `ServeMux`
* Middlewares are applied:

  * Globally at the `App` level
  * Locally at each `router`
* Routers are resolved recursively

---

## 📌 Philosophy

Cafe **does not aim to replace** larger frameworks like Gin, Echo, or Fiber.

It’s ideal if you want to:

* Learn how an HTTP framework works internally
* Have full control over request flow
* Keep your stack small and explicit

---

## 🛠️ Roadmap (ideas)

* [ ] Route-level middleware
* [ ] Context helpers
* [ ] Error handling
* [ ] Parametrized route groups

---

Built with ❤️ and `net/http`.
