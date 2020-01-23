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
	LDAPTimeoutErrorsPerSec                    *prometheus.Desc
	LongRunningLDAPOperationsPerMin            *prometheus.Desc
	LDAPSearchesTimeLimitExceededPerMinute     *prometheus.Desc
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
	ActiveSyncRequestsPerSec                   *prometheus.Desc
	PingCommandsPending                        *prometheus.Desc
	SyncCommandsPerSec                         *prometheus.Desc
	AvailabilityRequestsSec                    *prometheus.Desc
	CurrentUniqueUsers                         *prometheus.Desc
	OWARequestsPerSec                          *prometheus.Desc
	AutodiscoverRequestsPerSec                 *prometheus.Desc
	ActiveTasks                                *prometheus.Desc
	CompletedTasks                             *prometheus.Desc
	QueuedTasks                                *prometheus.Desc
	RPCAveragedLatency                         *prometheus.Desc
	RPCRequests                                *prometheus.Desc
	ActiveUserCount                            *prometheus.Desc
	ConnectionCount                            *prometheus.Desc
	RPCOperationsPerSec                        *prometheus.Desc
	UserCount                                  *prometheus.Desc

	invalidProcName *regexp.Regexp
}

type win32_PerfRawData_MSExchangeRpcClientAccess_MSExchangeRpcClientAccess struct {
	Name string

	RPCAveragedLatency  uint64
	RPCRequests         uint64
	ActiveUserCount     uint64
	ConnectionCount     uint64
	RPCOperationsPerSec uint64
	UserCount           uint64
}

type win32_PerfRawData_MSExchangeADAccess_MSExchangeADAccessProcesses struct {
	Name string

	LDAPReadTime                           uint64
	LDAPSearchTime                         uint64
	LDAPTimeoutErrorsPerSec                uint64
	LongRunningLDAPOperationsPerMin        uint64
	LDAPSearchesTimeLimitExceededPerMinute uint64
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

	RequestsPerSec      uint64 // ActiveSyncRequestsPerSec
	PingCommandsPending uint64
	SyncCommandsPerSec  uint64
}

type win32_PerfRawData_MSExchangeAvailabilityService_MSExchangeAvailabilityService struct {
	Name string

	RequestsSec uint64 // AvailabilityRequestsSec TODO: Is this really correct?
}

type win32_PerfRawData_MSExchangeOWA_MSExchangeOWA struct {
	Name string

	CurrentUniqueUsers uint64
	RequestsPerSec     uint64 // OWARequestsPerSec
}

type win32_PerfRawData_MSExchangeAutodiscover_MSExchangeAutodiscover struct {
	Name string

	RequestsPerSec uint64 // AutodiscoverRequestsPerSec
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
		LDAPTimeoutErrorsPerSec:                    desc("ldap_timeout_errors_per_sec", []string{"name"}, "LDAP timeout errors per second"),
		LongRunningLDAPOperationsPerMin:            desc("long_running_ldap_operations_per_min", []string{"name"}, "Long Running LDAP operations pr minute"),
		LDAPSearchesTimeLimitExceededPerMinute:     desc("ldap_searches_time_limit_exceeded_per_min", []string{"name"}, "LDAP searches time limit exceeded per minute"),
		ExternalActiveRemoteDeliveryQueueLength:    desc("external_active_remote_delivery_queue_len", []string{"name"}, "External Active Remote Delivery Queue length"),
		InternalActiveRemoteDeliveryQueueLength:    desc("internal_active_remote_delivery_queue_len", []string{"name"}, "Internal Active Remote Delivery Queue length"),
		ActiveMailboxDeliveryQueueLength:           desc("active_mailbox_delivery_queue_len", []string{"name"}, "Active Mailbox Delivery Queue length"),
		RetryMailboxDeliveryQueueLength:            desc("retry_mailbox_delivery_queue_len", []string{"name"}, "Retry Mailbox Delivery Queue length"),
		UnreachableQueueLength:                     desc("unreachable_queue_len", []string{"name"}, "Unreachable Queue length"),
		ExternalLargestDeliveryQueueLength:         desc("external_largest_delivery_queue_len", []string{"name"}, "External Largest Delivery Queue length"),
		InternalLargestDeliveryQueueLength:         desc("inernal_largest_delivery_queue_len", []string{"name"}, "Internal Largest Delivery Queue length"),
		PoisonQueueLength:                          desc("poison_queue_len", []string{"name"}, "Poison Queue length"),
		IODatabaseReadsAverageLatency:              desc("io_database_reads_average_latency", []string{"name"}, "Average database read latency"),
		IODatabaseWritesAverageLatency:             desc("io_database_writes_average_latency", []string{"name"}, "Average database write latency"),
		IOLogWritesAverageLatency:                  desc("io_log_writes_average_latency", []string{"name"}, "Average Log Writes Latency"),
		IODatabaseReadsRecoveryAverageLatency:      desc("io_database_reads_recovery_average_latency", []string{"name"}, "Database reads recovery avrage latency"),
		IODatabaseWritesRecoveryAverageLatency:     desc("io_database_writes_recovery_average_latency", []string{"name"}, "Database writes recovery average latency"),
		MailboxServerLocatorAverageLatency:         desc("mailbox_server_locator_average_latency", []string{"name"}, "Exchange HTTP Proxy Mailbox Server Locator latency (avg)"),
		AverageAuthenticationLatency:               desc("average_authentication_latency", []string{"name"}, "Exchange HTTP Proxy Authentication Latency (avg)"),
		AverageClientAccessServerProcessingLatency: desc("average_client_access_server_processing_latency", []string{"name"}, "Exchange HTTP Proxy Client Access Server Processing Latency (avg)"),
		MailboxServerProxyFailureRate:              desc("mailbox_server_proxy_failure_rate", []string{"name"}, "Exchange HTTP Proxy Mailbox Server Proxy Failure Rate"),
		OutstandingProxyRequests:                   desc("outstanding_proxy_requests", []string{"name"}, "Exchange HTTP Proxy outstanding proxy requests"),
		ProxyRequestsPerSec:                        desc("proxy_requests_per_sec", []string{"name"}, "Exchange HTTP Proxy requests/s"),
		ActiveSyncRequestsPerSec:                   desc("active_sync_requests_per_sec", []string{"name"}, "Active Sync requests/s "),
		PingCommandsPending:                        desc("ping_commands_pending", []string{"name"}, "Pending Active Sync ping-commands"),
		SyncCommandsPerSec:                         desc("sync_commands_per_sec", []string{"name"}, "Active Sync sync-commands/s"),
		AvailabilityRequestsSec:                    desc("availability_requests_per_sec", []string{"name"}, "Availability Service / Availability requests/s (wat?)"),
		CurrentUniqueUsers:                         desc("current_unique_users", []string{"name"}, "Outlook Web Access current unique users"),
		OWARequestsPerSec:                          desc("owa_requests_per_sec", []string{"name"}, "Outlook Web Access requests/s"),
		AutodiscoverRequestsPerSec:                 desc("autodiscover_requests_per_sec", []string{"name"}, "Autodiscovery requests/s"),
		ActiveTasks:                                desc("active_tasks", []string{"name"}, "Active Workload Management Tasks"),
		CompletedTasks:                             desc("completed_tasks", []string{"name"}, "Completed Workload Management Tasks"),
		QueuedTasks:                                desc("queued_tasks", []string{"name"}, "Queued Workload Management Tasks"),
		RPCAveragedLatency:                         desc("rpc_averaged_latency", []string{"name"}, "RPC Client Access averaged latency"),
		RPCRequests:                                desc("rpc_requests", []string{"name"}, "RPC Client Access requests"),
		ActiveUserCount:                            desc("active_user_count", []string{"name"}, "RPC Client Access active user count"),
		ConnectionCount:                            desc("connection_count", []string{"name"}, "RPC Client Access connection count"),
		RPCOperationsPerSec:                        desc("rpc_operations_per_sec", []string{"name"}, "RPC Client Access operations per sec"),
		UserCount:                                  desc("user_count", []string{"name"}, "RPC Client Access user count"),

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
			c.LDAPTimeoutErrorsPerSec,
			prometheus.CounterValue,
			float64(proc.LDAPTimeoutErrorsPerSec),
			proc.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.LongRunningLDAPOperationsPerMin,
			prometheus.CounterValue,
			float64(proc.LongRunningLDAPOperationsPerMin),
			proc.Name,
		)
		ch <- prometheus.MustNewConstMetric(
			c.LDAPSearchesTimeLimitExceededPerMinute,
			prometheus.CounterValue,
			float64(proc.LDAPSearchesTimeLimitExceededPerMinute),
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
		c.ActiveSyncRequestsPerSec,
		prometheus.CounterValue,
		float64(activesync[0].RequestsPerSec), // ActiveSyncRequestsPerSec
	)
	ch <- prometheus.MustNewConstMetric(
		c.PingCommandsPending,
		prometheus.CounterValue,
		float64(activesync[0].PingCommandsPending),
	)
	ch <- prometheus.MustNewConstMetric(
		c.SyncCommandsPerSec,
		prometheus.CounterValue,
		float64(activesync[0].SyncCommandsPerSec),
	)

	var availservice []win32_PerfRawData_MSExchangeAvailabilityService_MSExchangeAvailabilityService
	if err := wmi.Query(queryAll(&availservice), &availservice); err != nil {
		return err
	}
	ch <- prometheus.MustNewConstMetric(
		c.AvailabilityRequestsSec,
		prometheus.CounterValue,
		float64(availservice[0].RequestsSec), // AvailabilityRequestsSec TODO: Correct?
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
		c.OWARequestsPerSec,
		prometheus.CounterValue,
		float64(owa[0].RequestsPerSec), // OWARequestsPerSec
	)

	var autodisc []win32_PerfRawData_MSExchangeAutodiscover_MSExchangeAutodiscover
	if err := wmi.Query(queryAll(&autodisc), &autodisc); err != nil {
		return err
	}
	ch <- prometheus.MustNewConstMetric(
		c.AutodiscoverRequestsPerSec,
		prometheus.CounterValue,
		float64(autodisc[0].RequestsPerSec), // AutodiscoveRequestsPerSec
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

	var rpcCliAccess []win32_PerfRawData_MSExchangeRpcClientAccess_MSExchangeRpcClientAccess
	if err := wmi.Query(queryAll(&rpcCliAccess), &rpcCliAccess); err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(
		c.RPCAveragedLatency,
		prometheus.CounterValue,
		float64(rpcCliAccess[0].RPCAveragedLatency),
	)

	ch <- prometheus.MustNewConstMetric(
		c.RPCRequests,
		prometheus.CounterValue,
		float64(rpcCliAccess[0].RPCRequests),
	)

	ch <- prometheus.MustNewConstMetric(
		c.ActiveUserCount,
		prometheus.CounterValue,
		float64(rpcCliAccess[0].ActiveUserCount),
	)

	ch <- prometheus.MustNewConstMetric(
		c.ConnectionCount,
		prometheus.CounterValue,
		float64(rpcCliAccess[0].ConnectionCount),
	)

	ch <- prometheus.MustNewConstMetric(
		c.RPCOperationsPerSec,
		prometheus.CounterValue,
		float64(rpcCliAccess[0].RPCOperationsPerSec),
	)

	ch <- prometheus.MustNewConstMetric(
		c.UserCount,
		prometheus.CounterValue,
		float64(rpcCliAccess[0].UserCount),
	)

	return nil
}
