# slogctx

[![Go Reference](https://pkg.go.dev/badge/github.com/chainguard-dev/slogctx.svg)](https://pkg.go.dev/github.com/chainguard-dev/slogctx)

Context aware slog

## Usage

### Context Logger

The context Logger can be used to use Loggers from the context. This is
sometimes preferred over the [Context Handler](#context-handler), since this can
make it easier to use different loggers in different contexts (e.g. testing).

This approach is heavily inspired by
[knative.dev/pkg/logging](https://pkg.go.dev/knative.dev/pkg/logging)

```go
func main() {
	log := slogctx.New(slog.Default).With("a", "b")
	ctx := slogctx.WithLogger(log)

	// Grab logger from context and use
	slogctx.FromContext(ctx).With("foo", "bar").Infof("hello world")

	// Package level context loggers are also aware
	slogctx.ErrorContext(ctx, "asdf")
}
```

```sh
2023/12/12 18:27:27 INFO hello world a=b foo=bar
2023/12/12 18:27:27 ERROR asdf a=b
```

#### Testing

The `slogtest` package provides utilities to make it easy to create loggers that
will use the native testing logging.

```go
func TestFoo(t *testing.T) {
	ctx := slogtest.TestContextWithLogger(t)

	for _, tc := range []string{"a", "b"} {
		t.Run(tc, func(t *testing.T) {
			slogctx.FromContext(ctx).Infof("hello world")
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
ok      github.com/chainguard-dev/slogctx/examples/logger
```

### Context Handler

The context Handler can be used to insert values from the context.

```go
func init() {
	slog.SetDefault(slog.New(slogctx.NewHandler(slog.NewTextHandler(os.Stdout, nil))))
}

func main() {
	ctx := context.Background()
	ctx = slogctx.WithValue(ctx, "foo", "bar")

	// Use slog package directly
	slog.InfoContext(ctx, "hello world", slog.Bool("baz", true))

	// glog / zap style (note: can't pass additional attributes)
	slogctx.Errorf(ctx, "hello %s", "world")
}
```

```sh
$ go run .
time=2023-12-12T14:29:02.336-05:00 level=INFO msg="hello world" baz=true foo=bar
time=2023-12-12T14:29:02.337-05:00 level=ERROR msg="hello world" foo=bar
```
