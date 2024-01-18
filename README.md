# 👞 clog

[![Go Reference](https://pkg.go.dev/badge/github.com/chainguard-dev/clog.svg)](https://pkg.go.dev/github.com/chainguard-dev/clog)

Context-aware [`slog`](https://pkg.go.dev/log/slog)

## Usage

### Context Logger

The context Logger can be used to use Loggers from the context. This is
sometimes preferred over the [Context Handler](#context-handler), since this can
make it easier to use different loggers in different contexts (e.g. testing).

This approach is heavily inspired by
[`knative.dev/pkg/logging`](https://pkg.go.dev/knative.dev/pkg/logging), but with [zero dependencies outside the standard library](https://github.com/chainguard-dev/clog/blob/main/go.mod) (compare with [`pkg/logging`'s deps](https://pkg.go.dev/knative.dev/pkg/logging?tab=imports)).

```go
package main

import (
	"context"
	"log/slog"

	"github.com/chainguard-dev/clog"
)

func main() {
	// One-time setup
	log := clog.New(slog.Default().Handler()).With("a", "b")
	ctx := clog.WithLogger(context.Background(), log)

	f(ctx)
}

func f(ctx context.Context) {
	// Grab logger from context and use.
	log := clog.FromContext(ctx)
	log.Info("in f")

	// Add logging context and pass on.
	ctx = clog.WithLogger(ctx, log.With("f", "hello"))
	g(ctx)
}

func g(ctx context.Context) {
	// Grab logger from context and use.
	log := clog.FromContext(ctx)
	log.Info("in g")

	// Package level context loggers are also aware
	clog.ErrorContext(ctx, "asdf")
}

```

```sh
$ go run .
2009/11/10 23:00:00 INFO in f a=b
2009/11/10 23:00:00 INFO in g a=b f=hello
2009/11/10 23:00:00 ERROR asdf a=b f=hello
```

#### Testing

The `slogtest` package provides utilities to make it easy to create loggers that
will use the native testing logging.

```go
func TestFoo(t *testing.T) {
	ctx := slogtest.TestContextWithLogger(t)

	for _, tc := range []string{"a", "b"} {
		t.Run(tc, func(t *testing.T) {
			clog.FromContext(ctx).Infof("hello world")
		})
	}
}
```

```sh
$ go test -v ./examples/logger
=== RUN   TestLog
=== RUN   TestLog/a
=== NAME  TestLog
    slogtest.go:20: time=2023-12-12T18:42:53.020-05:00 level=INFO msg="hello world"

=== RUN   TestLog/b
=== NAME  TestLog
    slogtest.go:20: time=2023-12-12T18:42:53.020-05:00 level=INFO msg="hello world"

--- PASS: TestLog (0.00s)
    --- PASS: TestLog/a (0.00s)
    --- PASS: TestLog/b (0.00s)
PASS
ok      github.com/chainguard-dev/clog/examples/logger
```

### Context Handler

The context Handler can be used to insert values from the context.

```go
func init() {
	slog.SetDefault(slog.New(clog.NewHandler(slog.NewTextHandler(os.Stdout, nil))))
}

func main() {
	ctx := context.Background()
	ctx = clog.WithValues(ctx, "foo", "bar")

	// Use slog package directly
	slog.InfoContext(ctx, "hello world", slog.Bool("baz", true))

	// glog / zap style (note: can't pass additional attributes)
	clog.ErrorContextf(ctx, "hello %s", "world")
}
```

```sh
$ go run .
time=2009-11-10T23:00:00.000Z level=INFO msg="hello world" baz=true foo=bar
time=2009-11-10T23:00:00.000Z level=ERROR msg="hello world" foo=bar
```
