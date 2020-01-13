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

type ExchangeCollector struct {
	// Class: Win32_PerfRawData_MSExchangeADAccess_MSExchangeADAccessProcesses
	LDAPReadTime   *prometheus.Desc
	LDAPSearchTime *prometheus.Desc

	// Class: Win32_PerfFormattedData_MSExchangeADAccess_MSExchangeADAccessDomainControllers
	LDAPSearchesTimedOutperMinute          *prometheus.Desc
	LongRunningLDAPOperationsPermin        *prometheus.Desc
	LDAPSearchesTimeLimitExceededperMinute *prometheus.Desc

	// Class: Win32_PerfRawData_MSExchangeTransportQueues_MSExchangeTransportQueues
	ExternalActiveRemoteDeliveryQueueLength *prometheus.Desc
	InternalActiveRemoteDeliveryQueueLength *prometheus.Desc
	ActiveMailboxDeliveryQueueLength        *prometheus.Desc
	RetryMailboxDeliveryQueueLength         *prometheus.Desc
	UnreachableQueueLength                  *prometheus.Desc
	ExternalLargestDeliveryQueueLength      *prometheus.Desc
	InternalLargestDeliveryQueueLength      *prometheus.Desc
	PoisonQueueLength                       *prometheus.Desc

	// Class: Win32_PerfRawData_ESE_MSExchangeDatabaseInstances
	IODatabaseReadsAverageLatency          *prometheus.Desc
	IODatabaseWritesAverageLatency         *prometheus.Desc
	IOLogWritesAverageLatency              *prometheus.Desc
	IODatabaseReadsRecoveryAverageLatency  *prometheus.Desc
	IODatabaseWritesRecoveryAverageLatency *prometheus.Desc
}

// Win32_PerfRawData_MSExchangeADAccess_MSExchangeADAccessProcesses uses reflection
type Win32_PerfRawData_MSExchangeADAccess_MSExchangeADAccessProcesses struct {
	Name string

	LDAPReadTime                           uint64
	LDAPSearchTime                         uint64
	LDAPSearchesTimedOutperMinute          uint64
	LongRunningLDAPOperationsPermin        uint64
	LDAPSearchesTimeLimitExceededperMinute uint64
}

// Win32_PerfFormattedData_MSExchangeADAccess_MSExchangeADAccessDomainControllers uses reflection
type Win32_PerfFormattedData_MSExchangeADAccess_MSExchangeADAccessDomainControllers struct {
	Name string

	LDAPSearchesTimedOutperMinute          uint64
	LongRunningLDAPOperationsPermin        uint64
	LDAPSearchesTimeLimitExceededperMinute uint64
}

// Win32_PerfRawData_MSExchangeTransportQueues_MSExchangeTransportQueues  uses reflection
type Win32_PerfRawData_MSExchangeTransportQueues_MSExchangeTransportQueues struct {
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

// Win32_PerfRawData_ESE_MSExchangeDatabaseInstances  uses reflection
type Win32_PerfRawData_ESE_MSExchangeDatabaseInstances struct {
	Name string

	IODatabaseReadsAverageLatency          placeholder
	IODatabaseWritesAverageLatency         placeholder
	IOLogWritesAverageLatency              placeholder
	IODatabaseReadsRecoveryAverageLatency  placeholder
	IODatabaseWritesRecoveryAverageLatency placeholder
}

func init() {
	Factories[subsystem] = NewExchangeCollector
}

// NewExchangeCollector returns a new Collector
func NewExchangeCollector() (Collector, error) {
	return &ExchangeCollector{
		// Class: Win32_PerfRawData_MSExchangeADAccess_MSExchangeADAccessProcesses
		LDAPReadTime: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "ldap_read_time"),
			"LDAP Read Time", []string{}, nil,
		),
		LDAPSearchTime: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "ldap_search_time"),
			"LDAP Search Time", []string{}, nil,
		),

		// Class: Win32_PerfFormattedData_MSExchangeADAccess_MSExchangeADAccessDomainControllers
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

		// Class: Win32_PerfRawData_MSExchangeTransportQueues_MSExchangeTransportQueues
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

// Collect sends
func (c *ExchangeCollector) Collect(ctx *ScrapeContext, ch chan<- prometheus.Metric) error {
	if desc, err := c.collect(ch); err != nil {
		log.Error("failed collecting Exchange metrics:", desc, err)
		return err
	}
	return nil
}

func (c *ExchangeCollector) collect(ch chan<- prometheus.Metric) (*prometheus.Desc, error) {
	{
		var dst []Win32_PerfRawData_MSExchangeADAccess_MSExchangeADAccessProcesses
		if err := wmi.Query(queryAll(dst), &dst); err != nil {
			return nil, err
		}

		for _, app := range dst {
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
		var dst []Win32_PerfFormattedData_MSExchangeADAccess_MSExchangeADAccessDomainControllers
		if err := wmi.Query(queryAll(dst), &dst); err != nil {
			return nil, err
		}

		for _, app := range dst {
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
		var dst []Win32_PerfRawData_MSExchangeTransportQueues_MSExchangeTransportQueues
		if err := wmi.Query(queryAll(dst), &dst); err != nil {
			return nil, err
		}

		for _, app := range dst {
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
	return nil, nil
}
