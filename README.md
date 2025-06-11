# anygo

**anygo** is an idiomatic Go library inspired by Rust's `Result<T, E>` type.
It enables safe and expressive error handling without relying on traditional `error`-only returns.

## Installation

```bash
go get github.com/daxartio/anygo
```

## Overview

Core type:

```go
type Result[T any]
```

It represents either a success value (Ok) or an error (Err).

## Basic Usage

```go
import (
    "errors"
    "fmt"

    "github.com/daxartio/anygo"
)

func compute() anygo.Result[int] {
    if success {
        return anygo.Ok(42)
    }
    return anygo.Err[int](errors.New("something went wrong"))
}

res := compute()
if res.IsOk() {
    val := res.MustUnwrap()
    fmt.Println("Value:", val)
} else {
    fmt.Println("Error occurred")
}
```

## API

### Creation

- `Ok(value T) Result[T]` — creates a successful result.
- `Err[T](err error) Result[T]` — creates a failed result.

### Inspection

- `IsOk() bool` — true if result is Ok.
- `IsErr() bool` — true if result is Err.

### Unwrapping

- `Unwrap() (T, error)` — returns value and error.
- `UnwrapOr(default T) T` — value or default if Err.
- `UnwrapOrElse(func() T) T` — value or result of fallback function.
- `MustUnwrap() T` — panics if Err.
- `Expect(msg string) T` — panics with message if Err.

### Combinators

- `Map(Result[T], func(T) U) Result[U]` — transforms value.
- `MapErr(func(error) error) Result[T]` — transforms error.
- `AndThen(Result[T], func(T) Result[U]) Result[U]` — chains computations.
- `Inspect(func(T)) Result[T]` — performs side effect if Ok.
- `Or(Result[T]) Result[T]` — fallback result if Err.
- `OrElse(func() Result[T]) Result[T]` — fallback result from function.
- `ToPtr() *T` — pointer to value or nil.

## License

MIT
