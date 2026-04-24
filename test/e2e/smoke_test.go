package e2e

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var endpoints = map[string]string{
	"gateway":    "http://localhost:8080",
	"identity":   "http://localhost:8081",
	"submission": "http://localhost:8082",
	"grading":    "http://localhost:8083",
	"vcs":        "http://localhost:8084",
}

func TestHealthz(t *testing.T) {
	if os.Getenv("LMS_E2E_SKIP") != "" {
		t.Skip("LMS_E2E_ set")
	}

	client := &http.Client{Timeout: 3 * time.Second}
	for name, base := range endpoints {
		t.Run(name+"_healthz", func(t *testing.T) {
			resp, err := client.Get(base + "/healthz")
			require.NoError(t, err)
			defer resp.Body.Close()
			require.Equal(t, http.StatusOK, resp.StatusCode)
		})
		t.Run(name+"_readyz", func(t *testing.T) {
			resp, err := client.Get(base + "/readyz")
			require.NoError(t, err)
			defer resp.Body.Close()
			require.Equal(t, http.StatusOK, resp.StatusCode)
		})
	}
}
