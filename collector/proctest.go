// +build windows

package collector

import (
	"errors"
	"github.com/StackExchange/wmi"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

// proctestCollector collects process IO data rates
type proctestCollector struct {
	IOReadBytesPersec  *prometheus.Desc
	IOWriteBytesPersec *prometheus.Desc
	PrivateBytes       *prometheus.Desc
	PageFaultsPersec   *prometheus.Desc
	ThreadCount        *prometheus.Desc
}

type win32_PerfRawData_PerfProc_Process struct {
	Name string

	IOReadBytesPersec  uint64
	IOWriteBytesPersec uint64
	PrivateBytes       uint64
	PageFaultsPersec   uint64
	ThreadCount        uint64
}

func init() {
	Factories["proctest"] = newProctestCollector
}

func newProctestCollector() (Collector, error) {
	return &proctestCollector{
		IOReadBytesPersec: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "proctest", "proctest_io_read_bytes_persec"),
			"IO read bytes/s", []string{}, nil,
		),
		IOWriteBytesPersec: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "proctest", "proctest_io_write_bytes_persec"),
			"IO write bytes/s", []string{}, nil,
		),

		PrivateBytes: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "proctest", "proctest_private_bytes"),
			"IO read bytes/s", []string{}, nil,
		),

		PageFaultsPersec: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "proctest", "proctest_page_faults_persec"),
			"Page Faults/s", []string{}, nil,
		),
		ThreadCount: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "proctest", "proctest_thread_count"),
			"Thread Count", []string{}, nil,
		),
	}, nil
}

// Collect collects Exchange-metrics and provides them to prometheus through the ch channel
func (c *proctestCollector) Collect(ctx *ScrapeContext, ch chan<- prometheus.Metric) error {
	{
		var data []win32_PerfRawData_PerfProc_Process
		if err := wmi.Query(queryAll(data), &data); err != nil {
			log.Errorf("WMI query error while collecting %s-metrics: %s", "proctest", err)
			return err
		}

		if len(data) == 0 {
			log.Errorf("Query returned 0 rows")
			return errors.New("No rows returned")
		}

		for _, app := range data {
			ch <- prometheus.MustNewConstMetric(
				c.IOReadBytesPersec,
				prometheus.GaugeValue,
				float64(app.IOReadBytesPersec),
			)
			ch <- prometheus.MustNewConstMetric(
				c.IOWriteBytesPersec,
				prometheus.GaugeValue,
				float64(app.IOWriteBytesPersec),
			)
			ch <- prometheus.MustNewConstMetric(
				c.PrivateBytes,
				prometheus.GaugeValue,
				float64(app.PrivateBytes),
			)
			ch <- prometheus.MustNewConstMetric(
				c.PageFaultsPersec,
				prometheus.GaugeValue,
				float64(app.PageFaultsPersec),
			)
			ch <- prometheus.MustNewConstMetric(
				c.ThreadCount,
				prometheus.GaugeValue,
				float64(app.ThreadCount),
			)
		}
	}

	return nil
}
