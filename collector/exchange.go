// +build windows

package collector

import (
	"github.com/StackExchange/wmi"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

const subsystem string = "exchange"

// placeholder values should be replaced
type placeholder uint64

type exchangeCollector struct {
	LDAPReadTime   *prometheus.Desc
	LDAPSearchTime *prometheus.Desc

	LDAPSearchesTimedOutperMinute          *prometheus.Desc
	LongRunningLDAPOperationsPermin        *prometheus.Desc
	LDAPSearchesTimeLimitExceededperMinute *prometheus.Desc

	ExternalActiveRemoteDeliveryQueueLength *prometheus.Desc
	InternalActiveRemoteDeliveryQueueLength *prometheus.Desc
	ActiveMailboxDeliveryQueueLength        *prometheus.Desc
	RetryMailboxDeliveryQueueLength         *prometheus.Desc
	UnreachableQueueLength                  *prometheus.Desc
	ExternalLargestDeliveryQueueLength      *prometheus.Desc
	InternalLargestDeliveryQueueLength      *prometheus.Desc
	PoisonQueueLength                       *prometheus.Desc

	IODatabaseReadsAverageLatency          *prometheus.Desc
	IODatabaseWritesAverageLatency         *prometheus.Desc
	IOLogWritesAverageLatency              *prometheus.Desc
	IODatabaseReadsRecoveryAverageLatency  *prometheus.Desc
	IODatabaseWritesRecoveryAverageLatency *prometheus.Desc
}

type win32_PerfRawData_MSExchangeADAccess_MSExchangeADAccessProcesses struct {
	Name string

	LDAPReadTime                           uint64
	LDAPSearchTime                         uint64
	LDAPSearchesTimedOutperMinute          uint64
	LongRunningLDAPOperationsPermin        uint64
	LDAPSearchesTimeLimitExceededperMinute uint64
}

type win32_PerfFormattedData_MSExchangeADAccess_MSExchangeADAccessDomainControllers struct {
	Name string

	LDAPSearchesTimedOutperMinute          uint64
	LongRunningLDAPOperationsPermin        uint64
	LDAPSearchesTimeLimitExceededperMinute uint64
}

type win32_PerfRawData_MSExchangeTransportQueues_MSExchangeTransportQueues struct {
	Name string

	ExternalActiveRemoteDeliveryQueueLength placeholder
	InternalActiveRemoteDeliveryQueueLength placeholder
	ActiveMailboxDeliveryQueueLength        placeholder
	RetryMailboxDeliveryQueueLength         placeholder
	UnreachableQueueLength                  placeholder
	ExternalLargestDeliveryQueueLength      placeholder
	InternalLargestDeliveryQueueLength      placeholder
	PoisonQueueLength                       placeholder
}

type win32_PerfRawData_ESE_MSExchangeDatabaseInstances struct {
	Name string

	IODatabaseReadsAverageLatency          placeholder
	IODatabaseWritesAverageLatency         placeholder
	IOLogWritesAverageLatency              placeholder
	IODatabaseReadsRecoveryAverageLatency  placeholder
	IODatabaseWritesRecoveryAverageLatency placeholder
}

func init() {
	Factories[subsystem] = newExchangeCollector
}

// newExchangeCollector returns a new Collector
func newExchangeCollector() (Collector, error) {
	return &exchangeCollector{
		LDAPReadTime: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "ldap_read_time"),
			"LDAP Read Time", []string{}, nil,
		),
		LDAPSearchTime: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "ldap_search_time"),
			"LDAP Search Time", []string{}, nil,
		),
		LDAPSearchesTimedOutperMinute: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "ldap_searches_timed_out_per_min"),
			"LDAP Searches timeout pr minute", []string{}, nil,
		),
		LongRunningLDAPOperationsPermin: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "ldap_long_running_ops_per_min"),
			"Long Running LDAP operations pr minute", []string{}, nil,
		),
		LDAPSearchesTimeLimitExceededperMinute: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "ldap_searches_timed_out_per_min"),
			"LDAP Searches Time Limit Exceeded pr minute", []string{}, nil,
		),
		ExternalActiveRemoteDeliveryQueueLength: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "ext_active_remote_delivery_queue"),
			"External Active Remote Delivery Queue Length", []string{}, nil,
		),
		InternalActiveRemoteDeliveryQueueLength: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "internal_active_remote_delivery_queue"),
			"Internal Active Remote Delivery Queue Length", []string{}, nil,
		),
		ActiveMailboxDeliveryQueueLength: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "active_mailbox_delivery_queue"),
			"Active Mailbox Delivery Queue Length", []string{}, nil,
		),
		RetryMailboxDeliveryQueueLength: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "retry_mailbox_delivery_queue"),
			"Retry Mailbox Delivery Queue Length", []string{}, nil,
		),
		UnreachableQueueLength: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "unreachable_queue"),
			"Unreachable Queue Length", []string{}, nil,
		),
		ExternalLargestDeliveryQueueLength: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "external_largest_delivery_queue"),
			"External Largest Delivery Queue Length", []string{}, nil,
		),
		InternalLargestDeliveryQueueLength: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "inernal_largest_delivery_queue"),
			"Internal Largest Delivery Queue Length", []string{}, nil,
		),
		PoisonQueueLength: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "poison_queue"),
			"Poison Queue Length", []string{}, nil,
		),
	}, nil
}

// Collect collects Exchange-metrics and provides them to prometheus through the ch channel
func (c *exchangeCollector) Collect(ctx *ScrapeContext, ch chan<- prometheus.Metric) error {
	{
		var data []win32_PerfRawData_MSExchangeADAccess_MSExchangeADAccessProcesses
		if err := wmi.Query(queryAll(data), &data); err != nil {
			log.Errorf("WMI query error while collecting %s-metrics: %s", subsystem, err)
			return err
		}

		for _, app := range data {
			ch <- prometheus.MustNewConstMetric(
				c.LDAPReadTime,
				prometheus.CounterValue,
				float64(app.LDAPReadTime),
			)

			ch <- prometheus.MustNewConstMetric(
				c.LDAPSearchTime,
				prometheus.CounterValue,
				float64(app.LDAPSearchTime),
			)
		}
	}

	{
		var data []win32_PerfFormattedData_MSExchangeADAccess_MSExchangeADAccessDomainControllers
		if err := wmi.Query(queryAll(data), &data); err != nil {
			log.Errorf("WMI query error while collecting %s-metrics: %s", subsystem, err)
			return err
		}

		for _, app := range data {
			ch <- prometheus.MustNewConstMetric(
				c.LDAPSearchesTimedOutperMinute,
				prometheus.CounterValue,
				float64(app.LDAPSearchesTimedOutperMinute),
			)
			ch <- prometheus.MustNewConstMetric(
				c.LongRunningLDAPOperationsPermin,
				prometheus.CounterValue,
				float64(app.LongRunningLDAPOperationsPermin),
			)
			ch <- prometheus.MustNewConstMetric(
				c.LDAPSearchesTimeLimitExceededperMinute,
				prometheus.CounterValue,
				float64(app.LDAPSearchesTimeLimitExceededperMinute),
			)
		}
	}

	{
		var data []win32_PerfRawData_MSExchangeTransportQueues_MSExchangeTransportQueues
		if err := wmi.Query(queryAll(data), &data); err != nil {
			log.Errorf("WMI query error while collecting %s-metrics: %s", subsystem, err)
			return err
		}

		for _, app := range data {
			ch <- prometheus.MustNewConstMetric(
				c.ExternalActiveRemoteDeliveryQueueLength,
				prometheus.CounterValue,
				float64(app.ExternalActiveRemoteDeliveryQueueLength),
			)
			ch <- prometheus.MustNewConstMetric(
				c.InternalActiveRemoteDeliveryQueueLength,
				prometheus.CounterValue,
				float64(app.InternalActiveRemoteDeliveryQueueLength),
			)
			ch <- prometheus.MustNewConstMetric(
				c.ActiveMailboxDeliveryQueueLength,
				prometheus.CounterValue,
				float64(app.ActiveMailboxDeliveryQueueLength),
			)
			ch <- prometheus.MustNewConstMetric(
				c.RetryMailboxDeliveryQueueLength,
				prometheus.CounterValue,
				float64(app.RetryMailboxDeliveryQueueLength),
			)
			ch <- prometheus.MustNewConstMetric(
				c.UnreachableQueueLength,
				prometheus.CounterValue,
				float64(app.UnreachableQueueLength),
			)
			ch <- prometheus.MustNewConstMetric(
				c.ExternalLargestDeliveryQueueLength,
				prometheus.CounterValue,
				float64(app.ExternalLargestDeliveryQueueLength),
			)
			ch <- prometheus.MustNewConstMetric(
				c.InternalLargestDeliveryQueueLength,
				prometheus.CounterValue,
				float64(app.InternalLargestDeliveryQueueLength),
			)
			ch <- prometheus.MustNewConstMetric(
				c.PoisonQueueLength,
				prometheus.CounterValue,
				float64(app.PoisonQueueLength),
			)
		}
	}

	return nil
}
