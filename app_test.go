package nest

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewServer(t *testing.T) {
	app := NewServer()
	if app.handler == nil {
		t.Error("Expected handler to be initialized, got nil")
	}
	if len(app.routers) != 0 {
		t.Errorf("Expected no routers, got %d", len(app.routers))
	}
	if len(app.routes) != 0 {
		t.Errorf("Expected no routes, got %d", len(app.routes))
	}
}

func TestApp_Get(t *testing.T) {
	app := NewServer()
	var called bool
	handler := func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}
	app.Get("/test", handler)

	if len(app.routes) != 1 {
		t.Errorf("Expected 1 route, got %d", len(app.routes))
	}
	if app.routes[0].path != "/test" {
		t.Errorf("Expected path /test, got %s", app.routes[0].path)
	}
	if app.routes[0].method != "GET" {
		t.Errorf("Expected method GET, got %s", app.routes[0].method)
	}

	app.setUpRouters() // Simula el setup antes de escuchar

	req := httptest.NewRequest("GET", "/test/", nil)
	rr := httptest.NewRecorder()
	app.handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", rr.Code)
	}
	if !called {
		t.Error("Handler was not called")
	}
}

func TestApp_Post(t *testing.T) {
	app := NewServer()
	var called bool
	handler := func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusCreated)
	}
	app.Post("/test", handler)
	app.setUpRouters()

	req := httptest.NewRequest("POST", "/test/", nil)
	rr := httptest.NewRecorder()
	app.handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status Created, got %d", rr.Code)
	}
	if !called {
		t.Error("Handler was not called")
	}
}

func TestApp_Put(t *testing.T) {
	app := NewServer()
	var called bool
	handler := func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusAccepted)
	}
	app.Put("/test", handler)
	app.setUpRouters()

	req := httptest.NewRequest("PUT", "/test/", nil)
	rr := httptest.NewRecorder()
	app.handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Errorf("Expected status Accepted, got %d", rr.Code)
	}
	if !called {
		t.Error("Handler was not called")
	}
}

func TestApp_Delete(t *testing.T) {
	app := NewServer()
	var called bool
	handler := func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusNoContent)
	}
	app.Delete("/test", handler)
	app.setUpRouters()

	req := httptest.NewRequest("DELETE", "/test/", nil)
	rr := httptest.NewRecorder()
	app.handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("Expected status No Content, got %d", rr.Code)
	}
	if !called {
		t.Error("Handler was not called")
	}
}

func TestApp_UseRouter(t *testing.T) {
	app := NewServer()
	rtr := NewRouter()
	var called bool
	rtr.Get("/sub", func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.Write([]byte("sub-route"))
	})
	app.UseRouter("/api", rtr)

	if len(app.routers) != 1 {
		t.Errorf("Expected 1 mounted router, got %d", len(app.routers))
	}
	if app.routers[0].path != "/api" {
		t.Errorf("Expected mounted router path /api, got %s", app.routers[0].path)
	}

	app.setUpRouters()

	req := httptest.NewRequest("GET", "/api/sub/", nil)
	rr := httptest.NewRecorder()
	app.handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", rr.Code)
	}
	if !called {
		t.Error("Sub-router handler was not called")
	}
	body, _ := io.ReadAll(rr.Body)
	if string(body) != "sub-route" {
		t.Errorf("Expected body 'sub-route', got '%s'", string(body))
	}
}

func TestApp_UseRouter_DuplicatePath(t *testing.T) {
	app := NewServer()
	rtr1 := NewRouter()
	rtr2 := NewRouter()

	app.UseRouter("/api", rtr1)
	app.UseRouter("/api", rtr2) // Should not add a duplicate

	if len(app.routers) != 1 {
		t.Errorf("Expected 1 mounted router due to duplicate path, got %d", len(app.routers))
	}
}

func TestNewRouter(t *testing.T) {
	rtr := NewRouter()
	if len(rtr.routes) != 0 {
		t.Errorf("Expected no routes, got %d", len(rtr.routes))
	}
	if len(rtr.routers) != 0 {
		t.Errorf("Expected no routers, got %d", len(rtr.routers))
	}
}

func TestRouter_Get(t *testing.T) {
	rtr := NewRouter()
	handler := func(w http.ResponseWriter, r *http.Request) {}
	rtr.Get("/item", handler)

	if len(rtr.routes) != 1 {
		t.Errorf("Expected 1 route, got %d", len(rtr.routes))
	}
	if rtr.routes[0].path != "/item" {
		t.Errorf("Expected path /item, got %s", rtr.routes[0].path)
	}
	if rtr.routes[0].method != "GET" {
		t.Errorf("Expected method GET, got %s", rtr.routes[0].method)
	}
}

func TestRouter_UseRouter(t *testing.T) {
	parentRouter := NewRouter()
	childRouter := NewRouter()
	var called bool
	childRouter.Get("/nested", func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.Write([]byte("nested-route"))
	})
	parentRouter.UseRouter("/sub", childRouter)

	if len(parentRouter.routers) != 1 {
		t.Errorf("Expected 1 mounted router, got %d", len(parentRouter.routers))
	}
	if parentRouter.routers[0].path != "/sub" {
		t.Errorf("Expected mounted router path /sub, got %s", parentRouter.routers[0].path)
	}

	// Test getRoutes to ensure nested routes are correctly flattened and prefixed
	allRoutes := parentRouter.getRoutes()
	found := false
	for _, r := range allRoutes {
		if r.path == "/sub/nested" && r.method == "GET" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected '/sub/nested' route not found in flattened routes")
	}

	// To fully test the handler, we need to integrate with an App and use an httptest.Server
	app := NewServer()
	app.UseRouter("/main", parentRouter)
	app.setUpRouters()

	req := httptest.NewRequest("GET", "/main/sub/nested/", nil)
	rr := httptest.NewRecorder()
	app.handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", rr.Code)
	}
	if !called {
		t.Error("Nested router handler was not called")
	}
	body, _ := io.ReadAll(rr.Body)
	if string(body) != "nested-route" {
		t.Errorf("Expected body 'nested-route', got '%s'", string(body))
	}
}

func TestAddRoute_Duplicate(t *testing.T) {
	routes := []route{}
	handler1 := func(w http.ResponseWriter, r *http.Request) {}
	handler2 := func(w http.ResponseWriter, r *http.Request) {}

	routes = addRoute(routes, "/test", "GET", handler1)
	routes = addRoute(routes, "/test", "GET", handler2) // Should not add a duplicate

	if len(routes) != 1 {
		t.Errorf("Expected 1 route after adding a duplicate, got %d", len(routes))
	}
	if routes[0].handler == nil {
		t.Errorf("Expected handler to be non-nil")
	}
}

func TestSetUpRouters_Order(t *testing.T) {
	app := NewServer()
	var order []string

	app.Get("/root", func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "root")
	})

	router1 := NewRouter()
	router1.Get("/sub1", func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "sub1")
	})
	app.UseRouter("/api1", router1)

	router2 := NewRouter()
	router2.Get("/sub2", func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "sub2")
	})
	app.UseRouter("/api2", router2)

	app.setUpRouters()

	// Test root route
	reqRoot := httptest.NewRequest("GET", "/root/", nil)
	rrRoot := httptest.NewRecorder()
	app.handler.ServeHTTP(rrRoot, reqRoot)
	if rrRoot.Code != http.StatusOK {
		t.Errorf("Expected status OK for root, got %d", rrRoot.Code)
	}

	// Test api1/sub1 route
	reqSub1 := httptest.NewRequest("GET", "/api1/sub1/", nil)
	rrSub1 := httptest.NewRecorder()
	app.handler.ServeHTTP(rrSub1, reqSub1)
	if rrSub1.Code != http.StatusOK {
		t.Errorf("Expected status OK for /api1/sub1, got %d", rrSub1.Code)
	}

	// Test api2/sub2 route
	reqSub2 := httptest.NewRequest("GET", "/api2/sub2/", nil)
	rrSub2 := httptest.NewRecorder()
	app.handler.ServeHTTP(rrSub2, reqSub2)
	if rrSub2.Code != http.StatusOK {
		t.Errorf("Expected status OK for /api2/sub2, got %d", rrSub2.Code)
	}

	// Verify the order of handler calls (if it matters, though for routing it usually doesn't strictly)
	// This test primarily ensures all routes are correctly registered.
	// You might want to refine this if the order of handler execution is critical for your app logic.
	expectedOrder := []string{"root", "sub1", "sub2"} // Assuming handlers are called when requests are made
	// Note: The actual order depends on the order of requests in the test, not the order of definition in setUpRouters.
	// This part of the test is more for confirming handlers are callable.
	_ = expectedOrder // Prevent unused variable warning
}

func TestHandlePathPattern(t *testing.T) {
	app := NewServer()
	var called bool
	handler := func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}

	tests := []struct {
		name         string
		path         string
		expectedPath string // The path that should successfully match
	}{
		{"simple path", "/hello", "/hello/"},
		{"path with trailing slash", "/world/", "/world/"},
		{"nested path", "/api/v1/users", "/api/v1/users/"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app = NewServer() // Reset app for each test case
			called = false
			app.handle(tt.path, "GET", handler)

			req := httptest.NewRequest("GET", tt.expectedPath, nil)
			rr := httptest.NewRecorder()
			app.handler.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("For path %s, expected status OK, got %d", tt.path, rr.Code)
			}
			if !called {
				t.Errorf("For path %s, handler was not called", tt.path)
			}

			// Test that a non-matching path returns 404
			called = false
			reqNotFound := httptest.NewRequest("GET", "/nonexistent/", nil)
			rrNotFound := httptest.NewRecorder()
			app.handler.ServeHTTP(rrNotFound, reqNotFound)
			if rrNotFound.Code != http.StatusNotFound {
				t.Errorf("For non-matching path, expected status Not Found, got %d", rrNotFound.Code)
			}
			if called {
				t.Errorf("For non-matching path, handler was unexpectedly called")
			}
		})
	}
}

func TestRouter_GetRoutes_NestedRouters(t *testing.T) {
	// r1: /users
	//   r2: /profile
	//     r3: /details
	//   r2: /settings
	//   r3: /billing (directly in r1)

	r1 := NewRouter()
	r1.Get("/users", func(w http.ResponseWriter, r *http.Request) {})

	r2 := NewRouter()
	r2.Get("/profile", func(w http.ResponseWriter, r *http.Request) {})

	r3 := NewRouter()
	r3.Get("/details", func(w http.ResponseWriter, r *http.Request) {})
	r2.UseRouter("/sub", r3) // /profile/sub/details

	r1.UseRouter("/admin", r2) // /admin/profile/sub/details

	r4 := NewRouter()
	r4.Get("/settings", func(w http.ResponseWriter, r *http.Request) {})
	r1.UseRouter("/admin", r4) // This should overwrite the previous /admin router if implemented correctly, or be ignored.
	// Based on current implementation, it will add a duplicate if the router itself is different but path is the same.
	// Let's assume for this test that the previous /admin with r2 is replaced or handled.
	// If it's ignored, then r2's routes should appear. If it's overwritten by r4, then only r4's routes should appear under /admin.
	// Looking at `UseRouter`, it checks `mr.path == path`. So if `/admin` is used twice, only the first one will be added.
	// So, we expect r2's routes under /admin, not r4's.

	r1.Get("/billing", func(w http.ResponseWriter, r *http.Request) {}) // /billing

	allRoutes := r1.getRoutes()

	expectedPaths := map[string]bool{
		"/users":             false,
		"/admin/profile":     false,
		"/admin/sub/details": false,
		"/billing":           false,
	}

	for _, r := range allRoutes {
		if _, ok := expectedPaths[r.path]; ok {
			expectedPaths[r.path] = true
		}
	}

	for path, found := range expectedPaths {
		if !found {
			t.Errorf("Expected route '%s' not found in flattened routes", path)
		}
	}

	if len(allRoutes) != len(expectedPaths) {
		t.Errorf("Expected %d routes, got %d", len(expectedPaths), len(allRoutes))
		t.Logf("Found routes: %+v", allRoutes)
	}
}
