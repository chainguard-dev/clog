package gcp

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const CloudTraceHeader = "X-Cloud-Trace-Context"

func insideTest() bool {
	// Ask runtime.Callers for up to 10 PCs, including runtime.Callers itself.
	pc := make([]uintptr, 10)
	n := runtime.Callers(0, pc)
	if n == 0 {
		slog.Debug("WithCloudTraceContext: no PCs available")
		return true
	}
	frames := runtime.CallersFrames(pc[:n])
	for {
		frame, more := frames.Next()
		if !more {
			break
		}
		if strings.HasPrefix(frame.Function, "testing.") &&
			strings.HasSuffix(frame.File, "src/testing/testing.go") {
			slog.Debug("WithCloudTraceContext: inside test", "function", frame.Function, "file", frame.File, "line", frame.Line)
			return true
		}
	}
	return false
}

var (
	projectID  string
	lookupOnce sync.Once
)

func extractProjectID() string {
	// Get the project ID from the environment if specified
	fromEnv := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if fromEnv != "" {
		projectID = fromEnv
		return projectID
	}
	lookupOnce.Do(func() {
		if insideTest() {
			slog.Debug("WithCloudTraceContext: inside test, not looking up project ID")
			return
		}

		// By default use the metadata IP; otherwise use the environment variable
		// for consistency with https://pkg.go.dev/cloud.google.com/go/compute/metadata#Client.Get
		host := "169.254.169.254"
		if h := os.Getenv("GCE_METADATA_HOST"); h != "" {
			host = h
		}
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/computeMetadata/v1/project/project-id", host), nil)
		if err != nil {
			slog.Debug("WithCloudTraceContext: could not get GCP project ID from metadata server", "err", err)
			return
		}
		req.Header.Set("Metadata-Flavor", "Google")
		resp, err := (&http.Client{ // Timeouts copied from https://pkg.go.dev/cloud.google.com/go/compute/metadata#Get
			Transport: &http.Transport{
				Dial: (&net.Dialer{Timeout: 2 * time.Second}).Dial,
			},
			Timeout: 5 * time.Second,
		}).Do(req)
		if err != nil {
			slog.Debug("WithCloudTraceContext: could not get GCP project ID from metadata server", "err", err)
			return
		}
		if resp.StatusCode != http.StatusOK {
			slog.Debug("WithCloudTraceContext: could not get GCP project ID from metadata server", "code", resp.StatusCode, "status", resp.Status)
			return
		}
		defer resp.Body.Close()
		all, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Debug("WithCloudTraceContext: could not get GCP project ID from metadata server", "err", err)
			return
		}
		projectID = string(all)
	})
	return projectID
}

func gcpTraceUri(projectID, traceHeader string) string {
	traceParts := strings.Split(traceHeader, "/")
	if len(traceParts) > 0 && len(traceParts[0]) > 0 {
		return fmt.Sprintf("projects/%s/traces/%s", projectID, traceParts[0])
	}
	return ""
}

// WithCloudTraceContext returns an http.handler that adds the GCP Cloud Trace
// ID to the context. This is used to correlate the structured logs with the
// request log.
func WithCloudTraceContext(h http.Handler) http.Handler {
	projectID := extractProjectID()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if projectID != "" {
			if trace := gcpTraceUri(projectID, r.Header.Get(CloudTraceHeader)); trace != "" {
				r = r.WithContext(context.WithValue(r.Context(), "trace", trace))
			}
		}
		h.ServeHTTP(w, r)
	})
}

// CloudTraceContextUnaryInterceptor is a gRPC unary interceptor that adds the
// Cloud Trace ID to the context. This is used to correlate the structured
// logs with the request log.
func CloudTraceContextUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	projectID := extractProjectID()
	if projectID != "" {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			traceHeader := ""
			if x := md.Get(CloudTraceHeader); len(x) > 0 {
				traceHeader = x[0]
			}
			if trace := gcpTraceUri(projectID, traceHeader); trace != "" {
				ctx = context.WithValue(ctx, "trace", trace)
			}
		}
	}
	return handler(ctx, req)
}

func traceFromContext(ctx context.Context) string {
	trace := ctx.Value("trace")
	if trace == nil {
		return ""
	}
	return trace.(string)
}
