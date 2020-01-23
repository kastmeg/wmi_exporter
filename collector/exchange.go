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
	LDAPReadTime                               *prometheus.Desc
	LDAPSearchTime                             *prometheus.Desc
	LDAPTimeoutErrorsPersec                    *prometheus.Desc
	LongRunningLDAPOperationsPermin            *prometheus.Desc
	LDAPSearchesTimeLimitExceededperMinute     *prometheus.Desc
	ExternalActiveRemoteDeliveryQueueLength    *prometheus.Desc
	InternalActiveRemoteDeliveryQueueLength    *prometheus.Desc
	ActiveMailboxDeliveryQueueLength           *prometheus.Desc
	RetryMailboxDeliveryQueueLength            *prometheus.Desc
	UnreachableQueueLength                     *prometheus.Desc
	ExternalLargestDeliveryQueueLength         *prometheus.Desc
	InternalLargestDeliveryQueueLength         *prometheus.Desc
	PoisonQueueLength                          *prometheus.Desc
	IODatabaseReadsAverageLatency              *prometheus.Desc
	IODatabaseWritesAverageLatency             *prometheus.Desc
	IOLogWritesAverageLatency                  *prometheus.Desc
	IODatabaseReadsRecoveryAverageLatency      *prometheus.Desc
	IODatabaseWritesRecoveryAverageLatency     *prometheus.Desc
	MailboxServerLocatorAverageLatency         *prometheus.Desc
	AverageAuthenticationLatency               *prometheus.Desc
	AverageClientAccessServerProcessingLatency *prometheus.Desc
	MailboxServerProxyFailureRate              *prometheus.Desc
	OutstandingProxyRequests                   *prometheus.Desc
	ProxyRequestsPerSec                        *prometheus.Desc
	ActiveSyncRequestsPersec                   *prometheus.Desc // name collision
	PingCommandsPending                        *prometheus.Desc
	SyncCommandsPersec                         *prometheus.Desc
	AvailabilityRequestssec                    *prometheus.Desc
	CurrentUniqueUsers                         *prometheus.Desc
	OWARequestsPersec                          *prometheus.Desc // name collision
	AutodiscoverRequestsPersec                 *prometheus.Desc // name collision
	ActiveTasks                                *prometheus.Desc
	CompletedTasks                             *prometheus.Desc
	QueuedTasks                                *prometheus.Desc

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

type win32_PerfRawData_MSExchangeHttpProxy_MSExchangeHttpProxy struct {
	Name string

	MailboxServerLocatorAverageLatency         uint64
	AverageAuthenticationLatency               uint64
	AverageClientAccessServerProcessingLatency uint64
	MailboxServerProxyFailureRate              uint64
	OutstandingProxyRequests                   uint64
	ProxyRequestsPerSec                        uint64
}

type win32_PerfRawData_MSExchangeActiveSync_MSExchangeActiveSync struct {
	Name string

	RequestsPersec      uint64 // ActiveSyncRequestsPersec
	PingCommandsPending uint64
	SyncCommandsPersec  uint64
}

type win32_PerfRawData_MSExchangeAvailabilityService_MSExchangeAvailabilityService struct {
	Name string

	Requestssec uint64 // AvailabilityRequestssec TODO: Is this really correct?
}

type win32_PerfRawData_MSExchangeOWA_MSExchangeOWA struct {
	Name string

	CurrentUniqueUsers uint64
	RequestsPersec     uint64 // OWARequestsPersec
}

type win32_PerfRawData_MSExchangeAutodiscover_MSExchangeAutodiscover struct {
	Name string

	RequestsPersec uint64 // AutodiscoverRequestsPersec
}

type win32_PerfRawData_MSExchangeWorkloadManagementWorkloads_MSExchangeWorkloadManagementWorkloads struct {
	Name string

	ActiveTasks    uint64
	CompletedTasks uint64
	QueuedTasks    uint64
}

func init() {
	Factories[subsystem] = newExchangeCollector
}

// desc creates a new prometheus description
func desc(metricName string, labels []string, desc string) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName(Namespace, subsystem, metricName), desc, labels, nil)
}

// newExchangeCollector returns a new Collector
func newExchangeCollector() (Collector, error) {
	return &exchangeCollector{
		LDAPReadTime:                               desc("ldap_read_time", []string{"name"}, "LDAP Read Time"),
		LDAPSearchTime:                             desc("ldap_search_time", []string{"name"}, "LDAP Search Time"),
		LDAPTimeoutErrorsPersec:                    desc("ldap_timeout_errors_per_sec", []string{"name"}, "LDAP timeout errors per second"),
		LongRunningLDAPOperationsPermin:            desc("long_running_ldap_operations_permin", []string{"name"}, "Long Running LDAP operations pr minute"),
		LDAPSearchesTimeLimitExceededperMinute:     desc("ldap_searches_time_limit_exceeded_per_minute", []string{"name"}, "LDAP searches time limit exceeded per minute"),
		ExternalActiveRemoteDeliveryQueueLength:    desc("external_active_remote_delivery_queue_length", []string{"name"}, "External Active Remote Delivery Queue Length"),
		InternalActiveRemoteDeliveryQueueLength:    desc("internal_active_remote_delivery_queue_length", []string{"name"}, "Internal Active Remote Delivery Queue Length"),
		ActiveMailboxDeliveryQueueLength:           desc("active_mailbox_delivery_queue_length", []string{"name"}, "Active Mailbox Delivery Queue Length"),
		RetryMailboxDeliveryQueueLength:            desc("retry_mailbox_delivery_queue_length", []string{"name"}, "Retry Mailbox Delivery Queue Length"),
		UnreachableQueueLength:                     desc("unreachable_queue_length", []string{"name"}, "Unreachable Queue Length"),
		ExternalLargestDeliveryQueueLength:         desc("external_largest_delivery_queue_length", []string{"name"}, "External Largest Delivery Queue Length"),
		InternalLargestDeliveryQueueLength:         desc("inernal_largest_delivery_queue_length", []string{"name"}, "Internal Largest Delivery Queue Length"),
		PoisonQueueLength:                          desc("poison_queue_length", []string{"name"}, "Poison Queue Length"),
		IODatabaseReadsAverageLatency:              desc("io_database_reads_average_latency", []string{"name"}, "Average database read latency"),
		IODatabaseWritesAverageLatency:             desc("io_database_writes_average_latency", []string{"name"}, "Average database write latency"),
		IOLogWritesAverageLatency:                  desc("io_log_writes_average_latency", []string{"name"}, "Average Log Writes Latency"),
		IODatabaseReadsRecoveryAverageLatency:      desc("io_database_reads_recovery_average_latency", []string{"name"}, "Database reads recovery avrage latency"),
		IODatabaseWritesRecoveryAverageLatency:     desc("io_database_writes_recovery_average_latency", []string{"name"}, "Database writes recovery average latency"),
		MailboxServerLocatorAverageLatency:         desc("mailbox_server_locator_average_latency", []string{"name"}, "Exchange HTTP Proxy Mailbox Server Locator latency (avg)"),
		AverageAuthenticationLatency:               desc("average_authentication_latency", []string{"name"}, "Exchange HTTP Proxy Authentication Latency (avg)"),
		AverageClientAccessServerProcessingLatency: desc("average_client_access_server_processing_latency", []string{"name"}, "Exchange HTTP Proxy Client Access Server Processing Latency (avg)"),
		MailboxServerProxyFailureRate:              desc("mailbox_server_proxy_failure_rate", []string{"name"}, "Exchange HTTP Proxy Mailbox Server Proxy Failure Rate"),
		OutstandingProxyRequests:                   desc("outstanding_proxy_requests", []string{"name"}, "Exchange HTTP Proxy Outstanding Proxy Requests"),
		ProxyRequestsPerSec:                        desc("proxy_requests_per_sec", []string{"name"}, "Exchange HTTP Proxy Requests/s"),
		ActiveSyncRequestsPersec:                   desc("active_sync_requests_per_sec", []string{"name"}, "Active Sync requests/s "),
		PingCommandsPending:                        desc("ping_commands_pending", []string{"name"}, "Pending Active Sync ping-commands"),
		SyncCommandsPersec:                         desc("sync_commands_per_sec", []string{"name"}, "Active Sync sync-commands/s"),
		AvailabilityRequestssec:                    desc("availability_requests_per_sec", []string{"name"}, "Availability Service / Availability Requests/s (wat?)"),
		CurrentUniqueUsers:                         desc("current_unique_users", []string{"name"}, "Outlook Web Access current unique users"),
		OWARequestsPersec:                          desc("owa_requests_per_sec", []string{"name"}, "Outlook Web Access requests/s"),
		AutodiscoverRequestsPersec:                 desc("autodiscover_requests_per_sec", []string{"name"}, "Autodiscovery requests/s"),
		ActiveTasks:                                desc("active_tasks", []string{"name"}, "Active Workload Management Tasks"),
		CompletedTasks:                             desc("completed_tasks", []string{"name"}, "Completed Workload Management Tasks"),
		QueuedTasks:                                desc("queued_tasks", []string{"name"}, "Queued Workload Management Tasks"),

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

	var httpproxy []win32_PerfRawData_MSExchangeHttpProxy_MSExchangeHttpProxy
	if err := wmi.Query(queryAll(&httpproxy), &httpproxy); err != nil {
		return err
	}
	ch <- prometheus.MustNewConstMetric(
		c.MailboxServerLocatorAverageLatency,
		prometheus.CounterValue,
		float64(httpproxy[0].MailboxServerLocatorAverageLatency),
	)
	ch <- prometheus.MustNewConstMetric(
		c.AverageAuthenticationLatency,
		prometheus.CounterValue,
		float64(httpproxy[0].AverageAuthenticationLatency),
	)
	ch <- prometheus.MustNewConstMetric(
		c.AverageClientAccessServerProcessingLatency,
		prometheus.CounterValue,
		float64(httpproxy[0].AverageClientAccessServerProcessingLatency),
	)
	ch <- prometheus.MustNewConstMetric(
		c.MailboxServerProxyFailureRate,
		prometheus.CounterValue,
		float64(httpproxy[0].MailboxServerProxyFailureRate),
	)
	ch <- prometheus.MustNewConstMetric(
		c.OutstandingProxyRequests,
		prometheus.CounterValue,
		float64(httpproxy[0].OutstandingProxyRequests),
	)
	ch <- prometheus.MustNewConstMetric(
		c.ProxyRequestsPerSec,
		prometheus.CounterValue,
		float64(httpproxy[0].ProxyRequestsPerSec),
	)

	var activesync []win32_PerfRawData_MSExchangeActiveSync_MSExchangeActiveSync
	if err := wmi.Query(queryAll(&activesync), &activesync); err != nil {
		return err
	}
	ch <- prometheus.MustNewConstMetric(
		c.ActiveSyncRequestsPersec,
		prometheus.CounterValue,
		float64(activesync[0].RequestsPersec), // ActiveSyncRequestsPersec
	)
	ch <- prometheus.MustNewConstMetric(
		c.PingCommandsPending,
		prometheus.CounterValue,
		float64(activesync[0].PingCommandsPending),
	)
	ch <- prometheus.MustNewConstMetric(
		c.SyncCommandsPersec,
		prometheus.CounterValue,
		float64(activesync[0].SyncCommandsPersec),
	)

	var availservice []win32_PerfRawData_MSExchangeAvailabilityService_MSExchangeAvailabilityService
	if err := wmi.Query(queryAll(&availservice), &availservice); err != nil {
		return err
	}
	ch <- prometheus.MustNewConstMetric(
		c.AvailabilityRequestssec,
		prometheus.CounterValue,
		float64(availservice[0].Requestssec), // AvailabilityRequestssec TODO: Correct?
	)

	var owa []win32_PerfRawData_MSExchangeOWA_MSExchangeOWA
	if err := wmi.Query(queryAll(&owa), &owa); err != nil {
		return err
	}
	ch <- prometheus.MustNewConstMetric(
		c.CurrentUniqueUsers,
		prometheus.CounterValue,
		float64(owa[0].CurrentUniqueUsers),
	)
	ch <- prometheus.MustNewConstMetric(
		c.OWARequestsPersec,
		prometheus.CounterValue,
		float64(owa[0].RequestsPersec), // OWARequestsPerSec
	)

	var autodisc []win32_PerfRawData_MSExchangeAutodiscover_MSExchangeAutodiscover
	if err := wmi.Query(queryAll(&autodiscq), &autodisc); err != nil {
		return err
	}
	ch <- prometheus.MustNewConstMetric(
		c.AutodiscoverRequestsPersec,
		prometheus.CounterValue,
		float64(autodisc[0].RequestsPersec), // AutodiscoveRequestsPersec
	)

	var mgmtworkload []win32_PerfRawData_MSExchangeWorkloadManagementWorkloads_MSExchangeWorkloadManagementWorkloads
	if err := wmi.Query(queryAll(&mgmtworkload), &mgmtworkload); err != nil {
		return err
	}
	ch <- prometheus.MustNewConstMetric(
		c.ActiveTasks,
		prometheus.CounterValue,
		float64(mgmtworkload[0].ActiveTasks),
	)
	ch <- prometheus.MustNewConstMetric(
		c.CompletedTasks,
		prometheus.CounterValue,
		float64(mgmtworkload[0].CompletedTasks),
	)
	ch <- prometheus.MustNewConstMetric(
		c.QueuedTasks,
		prometheus.CounterValue,
		float64(mgmtworkload[0].QueuedTasks),
	)

	return nil
}
