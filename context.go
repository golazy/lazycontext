// lazycontext provides acumulative context that can store values like [context.Context] or by type
package lazycontext

import (
	"context"
	"reflect"
	"time"
)

// AppContext is a small wrapper around the context.AppContext type that provides a more convenient API
type AppContext interface {
	context.Context
	AddValue(key any, value any) AppContext
}

// New creates a new Context
func New() AppContext {
	ctx := context.Background()
	ctx = context.WithValue(ctx, reflect.TypeFor[AppContext](), ctx)
	return &lctx{
		ctx: ctx,
	}
}

var ctxKey = struct{}{}

// FromContext retrieves the Context from a context.Context
func FromContext(ctx context.Context) (AppContext, bool) {
	if ctx == nil {
		return nil, false
	}
	c, ok := ctx.Value(ctxKey).(AppContext)
	if ok {
		return nil, false
	}
	return c, true
}

// NewWithContext creates a new Context from a context.Context
func NewWithContext(ctx context.Context) AppContext {
	return &lctx{ctx: context.WithValue(ctx, ctxKey, ctx)}
}

type lctx struct {
	ctx context.Context
}

func (c *lctx) base() context.Context {
	if c.ctx == nil {
		return context.Background()
	}
	return c.ctx
}

// Deadline returns the time when work done on behalf of this context
// should be canceled. Deadline returns ok==false when no deadline is
// set. Successive calls to Deadline return the same results.
func (ctx *lctx) Deadline() (deadline time.Time, ok bool) {
	return ctx.base().Deadline()
}

// Done returns a channel that's closed when work done on behalf of this
// context should be canceled. Done may return nil if this context can
// never be canceled. Successive calls to Done return the same value.
// The close of the Done channel may happen asynchronously,
// after the cancel function returns.
//
// WithCancel arranges for Done to be closed when cancel is called;
// WithDeadline arranges for Done to be closed when the deadline
// expires; WithTimeout arranges for Done to be closed when the timeout
// elapses.
//
// Done is provided for use in select statements:
//
//	// Stream generates values with DoSomething and sends them to out
//	// until DoSomething returns an error or ctx.Done is closed.
//	func Stream(ctx context.Context, out chan<- Value) error {
//		for {
//			v, err := DoSomething(ctx)
//			if err != nil {
//				return err
//			}
//			select {
//			case <-ctx.Done():
//				return ctx.Err()
//			case out <- v:
//			}
//		}
//	}
//
// See https://blog.golang.org/pipelines for more examples of how to use
// a Done channel for cancellation.
func (c *lctx) Done() <-chan struct{} {
	return c.base().Done()
}

// Err returns a non-nil error value after Done is closed.
// If Done is closed, Err returns a non-nil error explaining why:
// Canceled if the context was canceled
// or DeadlineExceeded if the context's deadline passed.
// After Err returns a non-nil error, successive calls to Err return the same error.
func (c *lctx) Err() error {
	return c.base().Err()
}

// Value returns the value associated with this context for key, or nil
// if no value is associated with key. Successive calls to Value with
// the same key returns the same result.
//
// Use context values only for request-scoped data that transits
// processes and API boundaries, not for passing optional parameters to
// functions.
//
// A key identifies a specific value in a Context. Functions that wish
// to store values in Context typically allocate a key in a global
// variable then use that key as the argument to context.WithValue and
// Context.Value. A key can be any type that supports equality;
// packages should define keys as an unexported type to avoid
// collisions.
//
// Packages that define a Context key should provide type-safe accessors
// for the values stored using that key:
//
//	// Package user defines a User type that's stored in Contexts.
//	package user
//
//	import "context"
//
//	// User is the type of value stored in the Contexts.
//	type User struct {...}
//
//	// key is an unexported type for keys defined in this package.
//	// This prevents collisions with keys defined in other packages.
//	type key int
//
//	// userKey is the key for user.User values in Contexts. It is
//	// unexported; clients use user.NewContext and user.FromContext
//	// instead of using this key directly.
//	var userKey key
//
//	// NewContext returns a new Context that carries value u.
//	func NewContext(ctx context.Context, u *User) context.Context {
//		return context.WithValue(ctx, userKey, u)
//	}
//
//	// FromContext returns the User value stored in ctx, if any.
//	func FromContext(ctx context.Context) (*User, bool) {
//		u, ok := ctx.Value(userKey).(*User)
//		return u, ok
//	}
func (c *lctx) Value(key any) any {
	return c.base().Value(key)
}

// AddValue adds a value to the context
func (c *lctx) AddValue(key any, value any) AppContext {
	c.ctx = context.WithValue(c.base(), key, value)
	return c
}

// Get retrieves the last value added to the context with the given type
func Get[T any](app context.Context) T {
	t, ok := app.Value(reflect.TypeOf(*new(T))).(T)
	if ok {
		return t
	}
	return *new(T)
}

// Set stores a value in the context using the type of the value as the key
func Set[T any](ctx AppContext, value T) {
	ctx.AddValue(reflect.TypeOf(*new(T)), value)
}
