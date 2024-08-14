package lazycontext

import (
	"fmt"
	"testing"
)

func TestContext(t *testing.T) {

	ctx := New()

	Set(ctx, "test")

	value := Get[string](ctx)

	if value != "test" {
		t.Errorf("Expected 'test', got %v", value)
	}

}

func ExampleContext() {
	ctx := New()

	// You can store values as with normal context
	var userKey string
	ctx.AddValue(userKey, "user_33")
	fmt.Println(ctx.Value(userKey))

	// Or if you need to reference only by specific type you can omit the key
	type myConfig struct{ Name string }
	Set(ctx, myConfig{Name: "test"})
	cfg := Get[myConfig](ctx)
	fmt.Println(cfg.Name)

	// Output:
	// user_33
	// test

}
