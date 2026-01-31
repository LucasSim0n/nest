# Cafe ☕

A minimalist HTTP micro-framework for Go, inspired by *Nest / Express*-style composition with routers and middlewares, but built **entirely on top of `net/http`**.

Cafe is designed to:

* Keep the core small and easy to understand
* Avoid external dependencies
* Compose HTTP applications using **mounted routers**, **middlewares**, and **clear HTTP methods**

---

## 🎯 Motivation

The Go standard library already provides a solid HTTP foundation, but building structured applications on top of `net/http` often leads to repetitive boilerplate or early adoption of heavy frameworks.

Cafe exists to fill that gap.

The goal is to offer:

* A **thin abstraction** over `net/http`, not a replacement
* Familiar composition patterns (routers + middlewares)
* A codebase small enough to read and understand in one sitting
* A learning-friendly framework that shows *how* things work

If you enjoy knowing exactly what happens when a request hits your server, Cafe is for you.

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

## ⚡ Quick Start

```go
package main

import (
    "net/http"
    "github.com/LucasSim0n/cafe"
)

func main() {
    app := cafe.NewServer()

    app.Get("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello Cafe ☕"))
    })

    app.Listen(":3000")
}
```

---

## 🚀 Usage

### 🛣️ HTTP Routes

Both the server and routers support the basic HTTP methods.

Cafe inherits native path parameter support from **`net/http` (Go 1.22+)**.

```go
app.Get("/users/{id}", handler)
app.Post("/users", postHandler)
app.Put("/users/{id}", putHandler)
app.Delete("/users/{id}", deleteHandler)
```

Access parameters:

```go
id := r.PathValue("id")
```

---

### 🧩 Routers

```go
api := cafe.NewRouter()
api.Get("/users", usersHandler)
app.UseRouter("/api", api)
```

This produces:

```
GET  /api/users
POST /api/users
```

---

### 🌳 Nested routers

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

### 🧠 Middlewares

A middleware is a function that wraps an `http.HandlerFunc`:

```go
type middleware func(next http.HandlerFunc) http.HandlerFunc
```

#### Global middleware

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

#### Router-level middleware

```go
api.Use(authMiddleware)
```

Applied only to that router’s routes (and its child routers).

---

#### Execution order

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

Cafe **does not aim to replace** larger frameworks.

---

## 🤝 Contributing

Contributions are welcome. Keep changes small, focused, and dependency-free.

---

Built with ❤️ and `net/http`.

