// +build windows

package collector

import (
	"errors"
	"github.com/StackExchange/wmi"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"strconv"
)

// proctestCollector collects process IO data rates
type ProctestCollector struct {
	IOReadBytesPersec  *prometheus.Desc
	IOWriteBytesPersec *prometheus.Desc
	PrivateBytes       *prometheus.Desc
	PageFaultsPersec   *prometheus.Desc
	ThreadCount        *prometheus.Desc
}

func init() {
	Factories["proctest"] = NewProctestCollector
}

// NewProctestCollector ...
func NewProctestCollector() (Collector, error) {
	return &ProctestCollector{
		IOReadBytesPersec: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "proctest", "io_read_bytes_persec"),
			"IO read bytes/s",
			[]string{"process", "process_id"},
			nil,
		),
		IOWriteBytesPersec: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "proctest", "io_write_bytes_persec"),
			"IO write bytes/s",
			[]string{"process", "process_id"},
			nil,
		),

		PrivateBytes: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "proctest", "private_bytes"),
			"IO read bytes/s",
			[]string{"process", "process_id"},
			nil,
		),

		PageFaultsPersec: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "proctest", "page_faults_persec"),
			"Page Faults/s",
			[]string{"process", "process_id"},
			nil,
		),
		ThreadCount: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "proctest", "thread_count"),
			"Thread Count",
			[]string{"process", "process_id"},
			nil,
		),
	}, nil
}

// Collect collects metrics
func (c *ProctestCollector) Collect(ctx *ScrapeContext, ch chan<- prometheus.Metric) error {
	{
		var procs []Win32_PerfRawData_PerfProc_Process
		q := queryAll(procs)
		log.Infof("Query: %s", q)
		if err := wmi.Query(q, &procs); err != nil {
			log.Errorf("WMI query error while collecting %s-metrics: %s", "proctest", err)
			return err
		}

		if len(procs) == 0 {
			log.Errorf("Query returned 0 rows")
			return errors.New("No rows returned")
		}

		for _, proc := range procs {

			if proc.Name == "_Total" {
				continue
			}

			name := proc.Name
			pid := strconv.FormatUint(uint64(proc.IDProcess), 10)

			ch <- prometheus.MustNewConstMetric(
				c.IOReadBytesPersec,
				prometheus.GaugeValue,
				float64(proc.IOReadBytesPersec),
				name,
				pid,
			)
			ch <- prometheus.MustNewConstMetric(
				c.IOWriteBytesPersec,
				prometheus.GaugeValue,
				float64(proc.IOWriteBytesPersec),
				name,
				pid,
			)
			ch <- prometheus.MustNewConstMetric(
				c.PrivateBytes,
				prometheus.GaugeValue,
				float64(proc.PrivateBytes),
				name,
				pid,
			)
			ch <- prometheus.MustNewConstMetric(
				c.PageFaultsPersec,
				prometheus.GaugeValue,
				float64(proc.PageFaultsPersec),
				name,
				pid,
			)
			ch <- prometheus.MustNewConstMetric(
				c.ThreadCount,
				prometheus.GaugeValue,
				float64(proc.ThreadCount),
				name,
				pid,
			)
		}
	}

	return nil
}
