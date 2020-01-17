// +build windows

package collector

import (
	"regexp"

	"github.com/StackExchange/wmi"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

const subsystem string = "exchange"

type exchangeCollector struct {
	LDAPReadTime   *prometheus.Desc
	LDAPSearchTime *prometheus.Desc

	LDAPTimeoutErrorsPersec                *prometheus.Desc
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
}

// newExchangeCollector returns a new Collector
func newExchangeCollector() (Collector, error) {
	return &exchangeCollector{
		LDAPReadTime: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "ldap_read_time"),
			"LDAP Read Time", []string{"name"}, nil,
		),
		LDAPSearchTime: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "ldap_search_time"),
			"LDAP Search Time", []string{"name"}, nil,
		),
		LDAPTimeoutErrorsPersec: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "ldap_timeout_errors_per_sec"),
			"LDAP Searches timeout pr minute", []string{"name"}, nil,
		),
		LongRunningLDAPOperationsPermin: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "ldap_long_running_ops_per_min"),
			"Long Running LDAP operations pr minute", []string{"name"}, nil,
		),
		LDAPSearchesTimeLimitExceededperMinute: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "ldap_searches_timed_out_per_min"),
			"LDAP Searches Time Limit Exceeded pr minute", []string{"name"}, nil,
		), /*
			ExternalActiveRemoteDeliveryQueueLength: prometheus.NewDesc(
				prometheus.BuildFQName(Namespace, subsystem, "ext_active_remote_delivery_queue"),
				"External Active Remote Delivery Queue Length", nil, nil,
			),
			InternalActiveRemoteDeliveryQueueLength: prometheus.NewDesc(
				prometheus.BuildFQName(Namespace, subsystem, "internal_active_remote_delivery_queue"),
				"Internal Active Remote Delivery Queue Length", nil, nil,
			),
			ActiveMailboxDeliveryQueueLength: prometheus.NewDesc(
				prometheus.BuildFQName(Namespace, subsystem, "active_mailbox_delivery_queue"),
				"Active Mailbox Delivery Queue Length", nil, nil,
			),
			RetryMailboxDeliveryQueueLength: prometheus.NewDesc(
				prometheus.BuildFQName(Namespace, subsystem, "retry_mailbox_delivery_queue"),
				"Retry Mailbox Delivery Queue Length", nil, nil,
			),
			UnreachableQueueLength: prometheus.NewDesc(
				prometheus.BuildFQName(Namespace, subsystem, "unreachable_queue"),
				"Unreachable Queue Length", nil, nil,
			),
			ExternalLargestDeliveryQueueLength: prometheus.NewDesc(
				prometheus.BuildFQName(Namespace, subsystem, "external_largest_delivery_queue"),
				"External Largest Delivery Queue Length", nil, nil,
			),
			InternalLargestDeliveryQueueLength: prometheus.NewDesc(
				prometheus.BuildFQName(Namespace, subsystem, "inernal_largest_delivery_queue"),
				"Internal Largest Delivery Queue Length", nil, nil,
			),
			// func NewDesc(fqName, help string, variableLabels []string, constLabels Labels) *Desc {
			PoisonQueueLength: prometheus.NewDesc(
				prometheus.BuildFQName(Namespace, subsystem, "poison_queue"),
				"Poison Queue Length", nil, nil,
			),*/

		invalidProcName: regexp.MustCompile(`#[0-9]{0,2}`),
	}, nil
}

// Collect collects Exchange-metrics and provides them to prometheus through the ch channel
func (c *exchangeCollector) Collect(ctx *ScrapeContext, ch chan<- prometheus.Metric) error {

	var procData []win32_PerfRawData_MSExchangeADAccess_MSExchangeADAccessProcesses
	if err := wmi.Query(queryAll(procData), &procData); err != nil {
		log.Errorf("WMI query error while collecting %s-metrics: %s", subsystem, err)
		return err
	}

	if len(procData) > 0 {

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
	} else {
		log.Warnln("Length of []procData is zero")
	}

	/*
		var data []win32_PerfRawData_MSExchangeTransportQueues_MSExchangeTransportQueues
		if err := wmi.Query(queryAll(data), &data); err != nil {
			log.Errorf("WMI query error while collecting %s-metrics: %s", subsystem, err)
			return err
		}

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
	*/
	return nil
}
