package anygo

import "fmt"

// Result represents a value of type T or an error.
type Result[T any] struct {
	value T
	err   error
}

// Ok returns a successful Result containing value.
//
// Example:
//
//	r := anygo.Ok("hello")
//	fmt.Println(r.IsOk()) // true
func Ok[T any](val T) Result[T] {
	return Result[T]{value: val}
}

// Err returns a failed Result containing an error.
//
// Example:
//
//	r := anygo.Err[string](errors.New("oops"))
//	fmt.Println(r.IsErr()) // true
func Err[T any](err error) Result[T] {
	return Result[T]{err: err}
}

// IsOk returns true if the Result has no error.
//
// Example:
//
//	r := anygo.Ok(123)
//	fmt.Println(r.IsOk()) // true
func (r Result[T]) IsOk() bool {
	return r.err == nil
}

// IsErr returns true if the Result has an error.
//
// Example:
//
//	r := anygo.Err[int](errors.New("fail"))
//	fmt.Println(r.IsErr()) // true
func (r Result[T]) IsErr() bool {
	return r.err != nil
}

// Unwrap returns the value and error.
//
// Example:
//
//	r := anygo.Ok("hi")
//	v, err := r.Unwrap()
//	fmt.Println(v, err) // "hi", nil
func (r Result[T]) Unwrap() (T, error) {
	return r.value, r.err
}

// UnwrapError returns the error if present, or nil if ok.
func (r Result[T]) UnwrapError() error {
	if r.IsOk() {
		return nil
	}
	return r.err
}

// UnwrapOr returns the value if ok, or the default otherwise.
//
// Example:
//
//	r := anygo.Err[string](errors.New("fail"))
//	fmt.Println(r.UnwrapOr("default")) // "default"
func (r Result[T]) UnwrapOr(def T) T {
	if r.IsOk() {
		return r.value
	}
	return def
}

// UnwrapOrElse returns the value if ok, or calls the fallback function otherwise.
//
// Example:
//
//	r := anygo.Err[string](errors.New("fail"))
//	fmt.Println(r.UnwrapOrElse(func() string { return "fallback" })) // "fallback"
func (r Result[T]) UnwrapOrElse(f func() T) T {
	if r.IsOk() {
		return r.value
	}
	return f()
}

// MustUnwrap returns the value or panics if there's an error.
//
// Example:
//
//	r := anygo.Ok("safe")
//	fmt.Println(r.MustUnwrap()) // "safe"
//
//	// anygo.Err[string](errors.New("boom")).MustUnwrap() // panics
func (r Result[T]) MustUnwrap() T {
	if r.IsErr() {
		panic(r.err)
	}
	return r.value
}

// Map applies a function to the value if ok, propagates error otherwise.
//
// Example:
//
//	r := anygo.Ok(2)
//	squared := anygo.Map(r, func(x int) int { return x * x })
//	fmt.Println(squared.MustUnwrap()) // 4
func Map[T any, U any](r Result[T], f func(T) U) Result[U] {
	if r.IsErr() {
		return Err[U](r.err)
	}
	return Ok(f(r.value))
}

// MapErr transforms the error if present.
func (r Result[T]) MapErr(f func(error) error) Result[T] {
	if r.IsOk() {
		return r
	}
	return Err[T](f(r.err))
}

// Map applies a function to the value if ok, propagates error otherwise.
//
// Example:
//
//	r := anygo.Ok(2)
//	squared := r.Map(func(x int) int { return x * x })
func (r Result[T]) Map(f func(T) T) Result[T] {
	if r.IsErr() {
		return r
	}
	return Ok(f(r.value))
}

// Inspect calls a function on the value if Result is Ok.
func (r Result[T]) Inspect(f func(T)) Result[T] {
	if r.IsOk() {
		f(r.value)
	}
	return r
}

// Expect panics with the provided message if Result is Err.
func (r Result[T]) Expect(msg string) T {
	if r.IsErr() {
		panic(fmt.Sprintf("%s: %v", msg, r.err))
	}
	return r.value
}

// ToPtr returns a pointer to the value if Ok, or nil if Err.
func (r Result[T]) ToPtr() *T {
	if r.IsOk() {
		return &r.value
	}
	return nil
}

// Or returns self if Ok, otherwise returns the alternative.
func (r Result[T]) Or(other Result[T]) Result[T] {
	if r.IsOk() {
		return r
	}
	return other
}

// OrElse calls the fallback function if Err.
func (r Result[T]) OrElse(f func() Result[T]) Result[T] {
	if r.IsOk() {
		return r
	}
	return f()
}

// Errorf adds context to the error if Result is Err.
func (r Result[T]) Errorf(format string, a ...any) Result[T] {
	if r.IsOk() {
		return r
	}
	return Err[T](fmt.Errorf("%s: %w", fmt.Sprintf(format, a...), r.err))
}

// AndThen chains another Result-producing function on success.
type andThenFunc[T any, U any] func(T) Result[U]

func AndThen[T any, U any](r Result[T], f andThenFunc[T, U]) Result[U] {
	if r.IsErr() {
		return Err[U](r.err)
	}
	return f(r.value)
}
