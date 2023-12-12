# slogctx

Context aware slog

## Usage

```go
func init() {
	slog.SetDefault(slog.New(slogctx.NewHandler(slog.NewTextHandler(os.Stdout, nil))))
}

func main() {
	ctx := context.Background()
	ctx = slogctx.With(ctx, "foo", "bar")

	// Use slog package directly
	slog.InfoContext(ctx, "hello world", slog.Bool("baz", true))

	// glog / zap style (note: can't pass additional attributes)
	slogctx.Errorf(ctx, "hello %s", "world")
}
```

```sh
$ go run .
time=2023-12-06T16:29:33.440-07:00 level=INFO msg="hello world" foo=bar
```
