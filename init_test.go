package victoriaMetricsInit

import (
	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"os"
	"runtime"
	"testing"
	"time"
)

func TestInitMetrics(t *testing.T) {
	t.Run("request", func(t *testing.T) {
		testHost := "testHost:8428"

		_ = os.Setenv("VICTORIA_METRICS_HOST", testHost)
		_ = os.Setenv("VICTORIA_METRICS_INTERVAL", "1s")

		gock.New("http://" + testHost).Get("/api/v1/import/prometheus").Reply(200)

		InitMetrics("testInstance")
		time.Sleep(time.Second * 2)
		assert.True(t, gock.IsDone())
	})

	t.Run("defaults", func(t *testing.T) {
		defaultTimeInterval = 2 * time.Second

		_ = os.Setenv("VICTORIA_METRICS_HOST", "")
		_ = os.Setenv("VICTORIA_METRICS_INTERVAL", "")

		gock.New("http://" + defaultHost).Get("/api/v1/import/prometheus").Reply(200)

		InitMetrics("testInstance")
		runtime.Gosched()
		time.Sleep(defaultTimeInterval + time.Millisecond*500)
		assert.True(t, gock.IsDone())
	})

	t.Run("errors", func(t *testing.T) {
		stdErrFile, _ := os.CreateTemp("", "stderr.txt")
		origStderr := os.Stderr
		os.Stderr = stdErrFile

		defer (func() {
			_ = stdErrFile.Close()
			_ = os.Remove(stdErrFile.Name())
			os.Stderr = origStderr
		})()

		_ = os.Setenv("VICTORIA_METRICS_HOST", "\n")
		_ = os.Setenv("VICTORIA_METRICS_INTERVAL", "")

		InitMetrics("testInstance")

		content, _ := os.ReadFile(stdErrFile.Name())
		assert.Contains(t, string(content), "Failed to init metrics: ")
	})
}
