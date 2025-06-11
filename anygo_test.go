package anygo_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/daxartio/anygo"
)

func TestOk(t *testing.T) {
	r := anygo.Ok(42)
	if !r.IsOk() || r.IsErr() {
		t.Fatal("expected Ok result")
	}
}

func TestErr(t *testing.T) {
	err := errors.New("fail")
	r := anygo.Err[int](err)
	if !r.IsErr() || r.IsOk() {
		t.Fatal("expected Err result")
	}
}

func TestUnwrapOr(t *testing.T) {
	r := anygo.Err[int](errors.New("fail"))
	if v := r.UnwrapOr(100); v != 100 {
		t.Fatalf("expected fallback value, got %d", v)
	}
}

func TestUnwrapOrElse(t *testing.T) {
	r := anygo.Err[int](errors.New("fail"))
	if v := r.UnwrapOrElse(func() int { return 99 }); v != 99 {
		t.Fatalf("expected fallback function value, got %d", v)
	}
}

func TestMustUnwrapOk(t *testing.T) {
	r := anygo.Ok("hello")
	v := r.MustUnwrap()
	if v != "hello" {
		t.Fatalf("expected 'hello', got %v", v)
	}
}

func TestMap(t *testing.T) {
	r := anygo.Ok(3)
	mapped := anygo.Map(r, func(i int) string { return fmt.Sprintf("%d!", i) })
	if v := mapped.MustUnwrap(); v != "3!" {
		t.Fatalf("expected '3!', got %v", v)
	}
}

func TestResultMap(t *testing.T) {
	r := anygo.Ok(3)
	mapped := r.Map(func(i int) int { return i + 1 })
	if v := mapped.MustUnwrap(); v != 4 {
		t.Fatalf("expected '4', got %v", v)
	}
}

func TestMapErr(t *testing.T) {
	r := anygo.Err[int](errors.New("fail"))
	wrapped := r.MapErr(func(e error) error {
		return fmt.Errorf("wrapped: %w", e)
	})
	if !wrapped.IsErr() || wrapped.UnwrapOr(0) != 0 {
		t.Fatal("expected mapped error")
	}
}

func TestAndThen(t *testing.T) {
	r := anygo.Ok(5)
	res := anygo.AndThen(r, func(i int) anygo.Result[string] {
		return anygo.Ok(fmt.Sprintf("%d ok", i))
	})
	if v := res.MustUnwrap(); v != "5 ok" {
		t.Fatalf("unexpected result: %v", v)
	}
}

func TestInspect(t *testing.T) {
	r := anygo.Ok("value")
	called := false
	r.Inspect(func(s string) {
		if s != "value" {
			t.Fatal("unexpected inspect value")
		}
		called = true
	})
	if !called {
		t.Fatal("inspect function was not called")
	}
}

func TestExpect(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	anygo.Err[int](errors.New("fail")).Expect("should not fail")
}

func TestToPtr(t *testing.T) {
	r := anygo.Ok(123)
	ptr := r.ToPtr()
	if ptr == nil || *ptr != 123 {
		t.Fatal("unexpected pointer value")
	}
}

func TestOr(t *testing.T) {
	err := anygo.Err[string](errors.New("bad"))
	fallback := anygo.Ok("ok")
	res := err.Or(fallback)
	if v := res.MustUnwrap(); v != "ok" {
		t.Fatal("expected fallback value")
	}
}

func TestOrElse(t *testing.T) {
	err := anygo.Err[int](errors.New("bad"))
	res := err.OrElse(func() anygo.Result[int] {
		return anygo.Ok(77)
	})
	if v := res.MustUnwrap(); v != 77 {
		t.Fatal("expected fallback value")
	}
}
