package gcpslog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

// WithCloudTraceContext returns an http.handler that adds the GCP Cloud Trace
// ID to the context. This is used to correlate the structured logs with the
// request log.
func WithCloudTraceContext(h http.Handler) http.Handler {
	// Get the project ID from the environment if specified
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		// By default use the metadata IP; otherwise use the environment variable
		// for consistency with https://pkg.go.dev/cloud.google.com/go/compute/metadata#Client.Get
		host := "169.254.169.254"
		if h := os.Getenv("GCE_METADATA_HOST"); h != "" {
			host = h
		}
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/computeMetadata/v1/project/project-id", host), nil)
		if err != nil {
			slog.Debug("WithCloudTraceContext: could not get GCP project ID from metadata server", "err", err)
			return h
		}
		req.Header.Set("Metadata-Flavor", "Google")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			slog.Debug("WithCloudTraceContext: could not get GCP project ID from metadata server", "err", err)
			return h
		}
		if resp.StatusCode != http.StatusOK {
			slog.Debug("WithCloudTraceContext: could not get GCP project ID from metadata server", "code", resp.StatusCode, "status", resp.Status)
			return h
		}
		defer resp.Body.Close()
		all, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Debug("WithCloudTraceContext: could not get GCP project ID from metadata server", "err", err)
			return h
		}
		projectID = string(all)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var trace string
		traceHeader := r.Header.Get("X-Cloud-Trace-Context")
		traceParts := strings.Split(traceHeader, "/")
		if len(traceParts) > 0 && len(traceParts[0]) > 0 {
			trace = fmt.Sprintf("projects/%s/traces/%s", projectID, traceParts[0])
		}
		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "trace", trace)))
	})
}

func traceFromContext(ctx context.Context) string {
	trace := ctx.Value("trace")
	if trace == nil {
		return ""
	}
	return trace.(string)
}
