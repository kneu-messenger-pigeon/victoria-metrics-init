package victoria_metrics_init

import (
	"fmt"
	"github.com/VictoriaMetrics/metrics"
	"os"
	"time"
)

const defaultHost = "victoria-metrics:8428"

var defaultTimeInterval = 10 * time.Second

func InitMetrics(instance string) {
	host := os.Getenv("VICTORIA_METRICS_HOST")
	if host == "" {
		host = defaultHost
	}

	timeInterval, err := time.ParseDuration(os.Getenv("VICTORIA_METRICS_INTERVAL"))
	if err != nil || timeInterval <= 0 {
		timeInterval = defaultTimeInterval
	}

	err = metrics.InitPush(
		fmt.Sprintf("http://%s/api/v1/import/prometheus", host),
		timeInterval,
		fmt.Sprintf(`instance="%s"`, instance),
		true,
	)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to init metrics: %s\n", err)
	}
}
