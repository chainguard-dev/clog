package gcp

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/chainguard-dev/clog"
)

func TestTrace(t *testing.T) {
	// This ensures the metadata server is not called at all during tests.
	md := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("metadata server called")
	}))
	defer md.Close()
	t.Setenv("GCE_METADATA_HOST", md.URL)

	slog.SetDefault(slog.New(NewHandler(slog.LevelDebug)))
	for _, c := range []struct {
		name      string
		env       string
		wantTrace bool
	}{
		{"no env set", "", false},
		{"env set", "my-project", true},
	} {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("GOOGLE_CLOUD_PROJECT", c.env)

			// Set up a server that logs a message with trace context added.
			slog.SetDefault(slog.New(NewHandler(slog.LevelDebug)))
			h := WithCloudTraceContext(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				clog.InfoContext(ctx, "hello world")

				// TODO: This doesn't propagate the trace context to the logger.
				//clog.FromContext(ctx).Info("hello world")

				if r.Header.Get("X-Cloud-Trace-Context") == "" {
					t.Error("got empty trace context header, want non-empty")
				}

				if found := ctx.Value("trace") != nil; found != c.wantTrace {
					t.Fatalf("got trace context %t, want %t", found, c.wantTrace)
					if c.wantTrace {
						if trace := ctx.Value("trace"); !strings.Contains(trace.(string), "/"+c.env+"/") {
							t.Errorf("got trace context %q, want %q", trace, c.env)
						}
					}
				}
			}))
			srv := httptest.NewServer(h)
			defer srv.Close()

			// Send a request to the server with a trace context header.
			req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("X-Cloud-Trace-Context", "trace/id/yay")
			if _, err := http.DefaultClient.Do(req); err != nil {
				t.Fatal(err)
			}
		})
	}
}
