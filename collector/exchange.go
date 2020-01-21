// +build windows

package collector

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/tabwriter"

	"github.com/StackExchange/wmi"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"gopkg.in/alecthomas/kingpin.v2"
)

const subsystem string = "exchange"

var (
	tabw            *tabwriter.Writer
	exchangeVerbose = kingpin.Flag("collector.exchange.verbose", "Be more verbose").Bool()
)

type exchangeCollector struct {
	LDAPReadTime                            *prometheus.Desc
	LDAPSearchTime                          *prometheus.Desc
	LDAPTimeoutErrorsPersec                 *prometheus.Desc
	LongRunningLDAPOperationsPermin         *prometheus.Desc
	LDAPSearchesTimeLimitExceededperMinute  *prometheus.Desc
	ExternalActiveRemoteDeliveryQueueLength *prometheus.Desc
	InternalActiveRemoteDeliveryQueueLength *prometheus.Desc
	ActiveMailboxDeliveryQueueLength        *prometheus.Desc
	RetryMailboxDeliveryQueueLength         *prometheus.Desc
	UnreachableQueueLength                  *prometheus.Desc
	ExternalLargestDeliveryQueueLength      *prometheus.Desc
	InternalLargestDeliveryQueueLength      *prometheus.Desc
	PoisonQueueLength                       *prometheus.Desc
	IODatabaseReadsAverageLatency           *prometheus.Desc
	IODatabaseWritesAverageLatency          *prometheus.Desc
	IOLogWritesAverageLatency               *prometheus.Desc
	IODatabaseReadsRecoveryAverageLatency   *prometheus.Desc
	IODatabaseWritesRecoveryAverageLatency  *prometheus.Desc

	invalidProcName *regexp.Regexp
}

type win32_PerfRawData_MSExchangeADAccess_MSExchangeADAccessProcesses struct {
	Name string

	LDAPReadTime                           uint64
	LDAPSearchTime                         uint64
	LDAPTimeoutErrorsPersec                uint64
	LongRunningLDAPOperationsPermin        uint64
	LDAPSearchesTimeLimitExceededperMinute uint64
}

type win32_PerfRawData_MSExchangeTransportQueues_MSExchangeTransportQueues struct {
	Name string

	ExternalActiveRemoteDeliveryQueueLength uint64
	InternalActiveRemoteDeliveryQueueLength uint64
	ActiveMailboxDeliveryQueueLength        uint64
	RetryMailboxDeliveryQueueLength         uint64
	UnreachableQueueLength                  uint64
	ExternalLargestDeliveryQueueLength      uint64
	InternalLargestDeliveryQueueLength      uint64
	PoisonQueueLength                       uint64
}

type win32_PerfRawData_ESE_MSExchangeDatabaseInstances struct {
	Name string

	IODatabaseReadsAverageLatency          uint64
	IODatabaseWritesAverageLatency         uint64
	IOLogWritesAverageLatency              uint64
	IODatabaseReadsRecoveryAverageLatency  uint64
	IODatabaseWritesRecoveryAverageLatency uint64
}

func init() {
	Factories[subsystem] = newExchangeCollector
	tabw = tabwriter.NewWriter(os.Stdout, 50, 50, 10, '.', tabwriter.AlignRight|tabwriter.Debug)

}

// desc creates a new prometheus description
func desc(metricName string, labels []string, desc string) *prometheus.Desc {
	if *exchangeVerbose {
		fmt.Fprintf(tabw, "wmi_exchange_%-50s %-50s {%s}\n", metricName, desc, strings.Join(labels, ","))
		tabw.Flush()
	}
	return prometheus.NewDesc(prometheus.BuildFQName(Namespace, subsystem, metricName), desc, labels, nil)
}

// newExchangeCollector returns a new Collector
func newExchangeCollector() (Collector, error) {
	return &exchangeCollector{
		LDAPReadTime:                            desc("ldap_read_time", []string{"name"}, "LDAP Read Time"),
		LDAPSearchTime:                          desc("ldap_search_time", []string{"name"}, "LDAP Search Time"),
		LDAPTimeoutErrorsPersec:                 desc("ldap_timeout_errors_per_sec", []string{"name"}, "LDAP timeout errors per second"),
		LongRunningLDAPOperationsPermin:         desc("ldap_long_running_ops_per_min", []string{"name"}, "Long Running LDAP operations pr minute"),
		LDAPSearchesTimeLimitExceededperMinute:  desc("ldap_searches_timed_out_per_min", []string{"name"}, "LDAP Searches Time Limit Exceeded pr minute"),
		ExternalActiveRemoteDeliveryQueueLength: desc("ext_active_remote_delivery_queue", []string{"name"}, "External Active Remote Delivery Queue Length"),
		InternalActiveRemoteDeliveryQueueLength: desc("internal_active_remote_delivery_queue", []string{"name"}, "Internal Active Remote Delivery Queue Length"),
		ActiveMailboxDeliveryQueueLength:        desc("active_mailbox_delivery_queue", []string{"name"}, "Active Mailbox Delivery Queue Length"),
		RetryMailboxDeliveryQueueLength:         desc("retry_mailbox_delivery_queue", []string{"name"}, "Retry Mailbox Delivery Queue Length"),
		UnreachableQueueLength:                  desc("unreachable_queue", []string{"name"}, "Unreachable Queue Length"),
		ExternalLargestDeliveryQueueLength:      desc("external_largest_delivery_queue", []string{"name"}, "External Largest Delivery Queue Length"),
		InternalLargestDeliveryQueueLength:      desc("inernal_largest_delivery_queue", []string{"name"}, "Internal Largest Delivery Queue Length"),
		PoisonQueueLength:                       desc("poison_queue", []string{"name"}, "Poison Queue Length"),
		IODatabaseReadsAverageLatency:           desc("io_db_avg_read_latency", []string{"name"}, "Average database read latency"),
		IODatabaseWritesAverageLatency:          desc("io_db_avg_write_latency", []string{"name"}, "Average database write latency"),
		IOLogWritesAverageLatency:               desc("io_log_writes_avg_latency", []string{"name"}, "Average Log Writes Latency"),
		IODatabaseReadsRecoveryAverageLatency:   desc("io_db_reads_recovery_avg_latency", []string{"name"}, "Database reads recovery avrage latency"),
		IODatabaseWritesRecoveryAverageLatency:  desc("io_db_writes_recovery_avg_latency", []string{"name"}, "Database writes recovery latency"),
		invalidProcName:                         regexp.MustCompile(`#[0-9]{0,2}`),
	}, nil
}

// Collect collects Exchange-metrics and provides them to prometheus through the ch channel
func (c *exchangeCollector) Collect(ctx *ScrapeContext, ch chan<- prometheus.Metric) error {
	var procData []win32_PerfRawData_MSExchangeADAccess_MSExchangeADAccessProcesses
	if err := wmi.Query(queryAll(procData), &procData); err != nil {
		log.Errorf("WMI query error while collecting %s-metrics: %s", subsystem, err)
		return err
	}

	for _, proc := range procData {
		if proc.Name == "_Total" {
			continue
		}
		// Skip processes with # or #n-suffix
		if c.invalidProcName.Match([]byte(proc.Name)) {
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			c.LDAPReadTime,
			prometheus.CounterValue,
			float64(proc.LDAPReadTime),
			proc.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.LDAPSearchTime,
			prometheus.CounterValue,
			float64(proc.LDAPSearchTime),
			proc.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.LDAPTimeoutErrorsPersec,
			prometheus.CounterValue,
			float64(proc.LDAPTimeoutErrorsPersec),
			proc.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.LongRunningLDAPOperationsPermin,
			prometheus.CounterValue,
			float64(proc.LongRunningLDAPOperationsPermin),
			proc.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.LDAPSearchesTimeLimitExceededperMinute,
			prometheus.CounterValue,
			float64(proc.LDAPSearchesTimeLimitExceededperMinute),
			proc.Name,
		)
	}

	var transportQueues []win32_PerfRawData_MSExchangeTransportQueues_MSExchangeTransportQueues
	if err := wmi.Query(queryAll(transportQueues), &transportQueues); err != nil {
		log.Errorf("WMI query error while collecting %s-metrics: %s", subsystem, err)
		return err
	}
	for _, queue := range transportQueues {
		ch <- prometheus.MustNewConstMetric(
			c.ExternalActiveRemoteDeliveryQueueLength,
			prometheus.CounterValue,
			float64(queue.ExternalActiveRemoteDeliveryQueueLength),
			queue.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.InternalActiveRemoteDeliveryQueueLength,
			prometheus.CounterValue,
			float64(queue.InternalActiveRemoteDeliveryQueueLength),
			queue.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.ActiveMailboxDeliveryQueueLength,
			prometheus.CounterValue,
			float64(queue.ActiveMailboxDeliveryQueueLength),
			queue.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.RetryMailboxDeliveryQueueLength,
			prometheus.CounterValue,
			float64(queue.RetryMailboxDeliveryQueueLength),
			queue.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.UnreachableQueueLength,
			prometheus.CounterValue,
			float64(queue.UnreachableQueueLength),
			queue.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.ExternalLargestDeliveryQueueLength,
			prometheus.CounterValue,
			float64(queue.ExternalLargestDeliveryQueueLength),
			queue.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.InternalLargestDeliveryQueueLength,
			prometheus.CounterValue,
			float64(queue.InternalLargestDeliveryQueueLength),
			queue.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.PoisonQueueLength,
			prometheus.CounterValue,
			float64(queue.PoisonQueueLength),
			queue.Name,
		)
	}

	var databaseInstances []win32_PerfRawData_ESE_MSExchangeDatabaseInstances
	if err := wmi.Query(queryAll(databaseInstances), &databaseInstances); err != nil {
		log.Errorf("WMI query error while collecting %s-metrics: %s", subsystem, err)
		return err
	}
	for _, instance := range databaseInstances {
		ch <- prometheus.MustNewConstMetric(
			c.IODatabaseReadsAverageLatency,
			prometheus.CounterValue,
			float64(instance.IODatabaseReadsAverageLatency),
			instance.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.IODatabaseWritesAverageLatency,
			prometheus.CounterValue,
			float64(instance.IODatabaseWritesAverageLatency),
			instance.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.IOLogWritesAverageLatency,
			prometheus.CounterValue,
			float64(instance.IOLogWritesAverageLatency),
			instance.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.IODatabaseReadsRecoveryAverageLatency,
			prometheus.CounterValue,
			float64(instance.IODatabaseReadsRecoveryAverageLatency),
			instance.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.IODatabaseWritesRecoveryAverageLatency,
			prometheus.CounterValue,
			float64(instance.IODatabaseWritesRecoveryAverageLatency),
			instance.Name,
		)
	}

	return nil
}
