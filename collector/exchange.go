// +build windows

package collector

import (
	"fmt"
	"strings"

	"github.com/StackExchange/wmi"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"gopkg.in/alecthomas/kingpin.v2"
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

	ActiveCollFuncs []collectorFunc
}

type win32_PerfRawData_MSExchangeRpcClientAccess_MSExchangeRpcClientAccess struct {
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
	RequestsPerSec      uint64
	RequestsTotal       uint64
	PingCommandsPending uint64
	SyncCommandsPerSec  uint64
}

type win32_PerfRawData_MSExchangeAvailabilityService_MSExchangeAvailabilityService struct {
	RequestsSec uint64
}

type win32_PerfRawData_MSExchangeOWA_MSExchangeOWA struct {
	CurrentUniqueUsers uint64
	RequestsPerSec     uint64
}

type win32_PerfRawData_MSExchangeAutodiscover_MSExchangeAutodiscover struct {
	RequestsPerSec uint64
}

type win32_PerfRawData_MSExchangeWorkloadManagementWorkloads_MSExchangeWorkloadManagementWorkloads struct {
	Name string

	ActiveTasks    uint64
	CompletedTasks uint64
	QueuedTasks    uint64
}

// collectorFunc is a function that collects metrics
type collectorFunc func(ch chan<- prometheus.Metric) error

var (
	// All available collector functions
	exchangeAllCollectorFuncs = []string{
		"ad_access_procs",
		"transport_queues",
		"database_instances",
		"http_proxy",
		"active_sync",
		"availability_service",
		"owa",
		"auto_descover",
		"management_workloads",
		"rpc_client_access",
	}

	exchangeCollectorFuncDesc map[string]string = map[string]string{
		"ad_access_procs":      "(WMI Class: win32_PerfRawData_MSExchangeADAccess_MSExchangeADAccessProcesses)",
		"transport_queues":     "(WMI Class: win32_PerfRawData_MSExchangeTransportQueues_MSExchangeTransportQueues)",
		"database_instances":   "(WMI Class: win32_PerfRawData_ESE_MSExchangeDatabaseInstances)",
		"http_proxy":           "(WMI Class: win32_PerfRawData_MSExchangeHttpProxy_MSExchangeHttpProxy)",
		"active_sync":          "(WMI Class: win32_PerfRawData_MSExchangeActiveSync_MSExchangeActiveSync)",
		"availability_service": "(WMI Class: win32_PerfRawData_MSExchangeAvailabilityService_MSExchangeAvailabilityService)",
		"owa":                  "(WMI Class: win32_PerfRawData_MSExchangeOWA_MSExchangeOWA)",
		"auto_descover":        "(WMI Class: win32_PerfRawData_MSExchangeAutodiscover_MSExchangeAutodiscover)",
		"management_workloads": "(WMI Class: win32_PerfRawData_MSExchangeWorkloadManagementWorkloads_MSExchangeWorkloadManagementWorkloads)",
		"rpc_client_access":    "(WMI Class: win32_PerfRawData_MSExchangeRpcClientAccess_MSExchangeRpcClientAccess)",
	}

	argExchangeListAllCollectors = kingpin.Flag(
		"collectors.exchange.list",
		"Lists all available exchange collectors and their description",
	).Bool()

	argExchangeEnabledCollectors = kingpin.Flag(
		"collectors.exchange.enable",
		"comma-separated list of exchange collectors to use",
	).Default(strings.Join(exchangeAllCollectorFuncs, ",")).String()

	argExchangeDisabledCollectors = kingpin.Flag(
		"collectors.exchange.disable",
		"comma-separated list of exchange collectors NOT to use",
	).Default().String()
)

func init() {
	registerCollector("exchange", newExchangeCollector)
}

// desc creates a new prometheus description
func desc(metricName string, description string, labels ...string) *prometheus.Desc {
	return prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, metricName),
		description,
		labels,
		nil,
	)
}

// contains checks if element e exists in slice s
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// toLabelName converts strings to lowercase and replaces all whitespace and dots with underscores
func toLabelName(name string) string {
	return strings.ReplaceAll(strings.Join(strings.Fields(strings.ToLower(name)), "_"), ".", "_")
}

// newExchangeCollector returns a new Collector
func newExchangeCollector() (Collector, error) {
	c := exchangeCollector{
		LDAPReadTime:                               desc("ldap_read_time", "LDAP Read Time", "name"),
		LDAPSearchTime:                             desc("ldap_search_time", "LDAP Search Time", "name"),
		LDAPTimeoutErrorsPerSec:                    desc("ldap_timeout_errors_per_sec", "LDAP timeout errors per second", "name"),
		LongRunningLDAPOperationsPerMin:            desc("long_running_ldap_operations_per_min", "Long Running LDAP operations pr minute", "name"),
		LDAPSearchesTimeLimitExceededPerMinute:     desc("ldap_searches_time_limit_exceeded_per_min", "LDAP searches time limit exceeded per minute", "name"),
		ExternalActiveRemoteDeliveryQueueLength:    desc("external_active_remote_delivery_queue_len", "External Active Remote Delivery Queue length", "name"),
		InternalActiveRemoteDeliveryQueueLength:    desc("internal_active_remote_delivery_queue_len", "Internal Active Remote Delivery Queue length", "name"),
		ActiveMailboxDeliveryQueueLength:           desc("active_mailbox_delivery_queue_len", "Active Mailbox Delivery Queue length", "name"),
		RetryMailboxDeliveryQueueLength:            desc("retry_mailbox_delivery_queue_len", "Retry Mailbox Delivery Queue length", "name"),
		UnreachableQueueLength:                     desc("unreachable_queue_len", "Unreachable Queue length", "name"),
		ExternalLargestDeliveryQueueLength:         desc("external_largest_delivery_queue_len", "External Largest Delivery Queue length", "name"),
		InternalLargestDeliveryQueueLength:         desc("inernal_largest_delivery_queue_len", "Internal Largest Delivery Queue length", "name"),
		PoisonQueueLength:                          desc("poison_queue_len", "Poison Queue length", "name"),
		IODatabaseReadsAverageLatency:              desc("io_database_reads_average_latency", "Average database read latency", "name"),
		IODatabaseWritesAverageLatency:             desc("io_database_writes_average_latency", "Average database write latency", "name"),
		IOLogWritesAverageLatency:                  desc("io_log_writes_average_latency", "Average Log Writes Latency", "name"),
		IODatabaseReadsRecoveryAverageLatency:      desc("io_database_reads_recovery_average_latency", "Database reads recovery avrage latency", "name"),
		IODatabaseWritesRecoveryAverageLatency:     desc("io_database_writes_recovery_average_latency", "Database writes recovery average latency", "name"),
		MailboxServerLocatorAverageLatency:         desc("mailbox_server_locator_average_latency", "HTTP Proxy Mailbox Server Locator latency (avg)", "name"),
		AverageAuthenticationLatency:               desc("average_authentication_latency", "HTTP Proxy Authentication Latency (avg)", "name"),
		AverageClientAccessServerProcessingLatency: desc("average_client_access_server_processing_latency", "HTTP Proxy Client Access Server Processing Latency (avg)", "name"),
		MailboxServerProxyFailureRate:              desc("mailbox_server_proxy_failure_rate", "HTTP Proxy Mailbox Server Proxy Failure Rate", "name"),
		OutstandingProxyRequests:                   desc("outstanding_proxy_requests", "HTTP Proxy outstanding proxy requests", "name"),
		ProxyRequestsPerSec:                        desc("proxy_requests_per_sec", "HTTP Proxy requests/s", "name"),
		ActiveSyncRequestsPerSec:                   desc("active_sync_requests_per_sec", "Active Sync requests/s "),
		PingCommandsPending:                        desc("ping_commands_pending", "Pending Active Sync ping-commands"),
		SyncCommandsPerSec:                         desc("sync_commands_per_sec", "Active Sync sync-commands/s"),
		AvailabilityRequestsSec:                    desc("availability_requests_per_sec", "Availability Service / Availability requests/s"),
		CurrentUniqueUsers:                         desc("current_unique_users", "Outlook Web Access current unique users"),
		OWARequestsPerSec:                          desc("owa_requests_per_sec", "Outlook Web Access requests/s"),
		AutodiscoverRequestsPerSec:                 desc("autodiscover_requests_per_sec", "Autodiscovery requests/s"),
		ActiveTasks:                                desc("active_tasks", "Active Workload Management Tasks"),
		CompletedTasks:                             desc("completed_tasks", "Completed Workload Management Tasks"),
		QueuedTasks:                                desc("queued_tasks", "Queued Workload Management Tasks"),
		RPCAveragedLatency:                         desc("rpc_averaged_latency", "RPC Client Access averaged latency"),
		RPCRequests:                                desc("rpc_requests", "RPC Client Access requests"),
		ActiveUserCount:                            desc("active_user_count", "RPC Client Access active user count"),
		ConnectionCount:                            desc("connection_count", "RPC Client Access connection count"),
		RPCOperationsPerSec:                        desc("rpc_operations_per_sec", "RPC Client Access operations per sec"),
		UserCount:                                  desc("user_count", "RPC Client Access user count"),
	}

	collectorFuncLookup := map[string]collectorFunc{
		"ad_access_procs":      c.collectADAccessProcs,
		"transport_queues":     c.collectTransportQueues,
		"database_instances":   c.collectDatabaseInstances,
		"http_proxy":           c.collectHTTPProxy,
		"active_sync":          c.collectActiveSync,
		"availability_service": c.collectAvailService,
		"owa":                  c.collectOWA,
		"auto_descover":        c.collectAutoDiscover,
		"management_workloads": c.collectMgmtWorkloads,
		"rpc_client_access":    c.collectRPCClientAccess,
	}

	// get the disabled and enabled collectors into slices
	disabledCollectors := strings.Split(*argExchangeDisabledCollectors, ",")
	enabledCollectors := strings.Split(*argExchangeEnabledCollectors, ",")

	// collFuncNames that are not also disabledCollectorFuncs gets added to the ActiveCollFuncs slice.
	for _, collFuncName := range enabledCollectors {
		collFunc, isValidName := collectorFuncLookup[collFuncName]

		if !isValidName {
			return nil, fmt.Errorf("No such collector function %s", collFuncName)
		}

		// skip collector func names that are explicitly disabled
		if contains(disabledCollectors, collFuncName) {
			continue
		}

		c.ActiveCollFuncs = append(c.ActiveCollFuncs, collFunc)
	}

	if *argExchangeListAllCollectors {
		state := ""
		for _, name := range exchangeAllCollectorFuncs {
			if contains(disabledCollectors, name) {
				state = "\x1b[31m => "
			}

			if contains(enabledCollectors, name) {
				state = "\x1b[32m => "
			}

			fmt.Printf("%-15s %-32s %-32s\n", state, name, exchangeCollectorFuncDesc[name])
		}
	}

	return &c, nil
}

// Collect collects exchange metrics and sends them to prometheus
func (c *exchangeCollector) Collect(ctx *ScrapeContext, ch chan<- prometheus.Metric) error {
	for _, collFunc := range c.ActiveCollFuncs {
		if err := collFunc(ch); err != nil {
			log.Errorf("Error in %s: %s", className(collFunc), err)
		}
	}
	return nil
}

func (c *exchangeCollector) collectADAccessProcs(ch chan<- prometheus.Metric) error {
	data := []win32_PerfRawData_MSExchangeADAccess_MSExchangeADAccessProcesses{}
	if err := wmi.Query("SELECT * FROM "+className(data), &data); err != nil {
		return err
	}
	for _, proc := range data {
		labelName := toLabelName(proc.Name)
		if strings.HasSuffix(labelName, "_total") {
			continue
		}

		ch <- prometheus.MustNewConstMetric(
			c.LDAPReadTime,
			prometheus.GaugeValue,
			float64(proc.LDAPReadTime),
			labelName,
		)
		ch <- prometheus.MustNewConstMetric(
			c.LDAPSearchTime,
			prometheus.GaugeValue,
			float64(proc.LDAPSearchTime),
			labelName,
		)
		ch <- prometheus.MustNewConstMetric(
			c.LDAPTimeoutErrorsPerSec,
			prometheus.GaugeValue,
			float64(proc.LDAPTimeoutErrorsPerSec),
			labelName,
		)
		ch <- prometheus.MustNewConstMetric(
			c.LongRunningLDAPOperationsPerMin,
			prometheus.GaugeValue,
			float64(proc.LongRunningLDAPOperationsPerMin),
			labelName,
		)
		ch <- prometheus.MustNewConstMetric(
			c.LDAPSearchesTimeLimitExceededPerMinute,
			prometheus.GaugeValue,
			float64(proc.LDAPSearchesTimeLimitExceededPerMinute),
			labelName,
		)
	}
	return nil
}

func (c *exchangeCollector) collectTransportQueues(ch chan<- prometheus.Metric) error {
	data := []win32_PerfRawData_MSExchangeTransportQueues_MSExchangeTransportQueues{}
	if err := wmi.Query("SELECT * FROM "+className(data), &data); err != nil {
		return err
	}
	for _, queue := range data {
		labelName := toLabelName(queue.Name)
		if strings.HasSuffix(labelName, "_total") {
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			c.ExternalActiveRemoteDeliveryQueueLength,
			prometheus.GaugeValue,
			float64(queue.ExternalActiveRemoteDeliveryQueueLength),
			labelName,
		)
		ch <- prometheus.MustNewConstMetric(
			c.InternalActiveRemoteDeliveryQueueLength,
			prometheus.GaugeValue,
			float64(queue.InternalActiveRemoteDeliveryQueueLength),
			labelName,
		)
		ch <- prometheus.MustNewConstMetric(
			c.ActiveMailboxDeliveryQueueLength,
			prometheus.GaugeValue,
			float64(queue.ActiveMailboxDeliveryQueueLength),
			labelName,
		)
		ch <- prometheus.MustNewConstMetric(
			c.RetryMailboxDeliveryQueueLength,
			prometheus.GaugeValue,
			float64(queue.RetryMailboxDeliveryQueueLength),
			labelName,
		)
		ch <- prometheus.MustNewConstMetric(
			c.UnreachableQueueLength,
			prometheus.GaugeValue,
			float64(queue.UnreachableQueueLength),
			labelName,
		)
		ch <- prometheus.MustNewConstMetric(
			c.ExternalLargestDeliveryQueueLength,
			prometheus.GaugeValue,
			float64(queue.ExternalLargestDeliveryQueueLength),
			labelName,
		)
		ch <- prometheus.MustNewConstMetric(
			c.InternalLargestDeliveryQueueLength,
			prometheus.GaugeValue,
			float64(queue.InternalLargestDeliveryQueueLength),
			labelName,
		)
		ch <- prometheus.MustNewConstMetric(
			c.PoisonQueueLength,
			prometheus.GaugeValue,
			float64(queue.PoisonQueueLength),
			labelName,
		)
	}
	return nil
}

func (c *exchangeCollector) collectDatabaseInstances(ch chan<- prometheus.Metric) error {
	data := []win32_PerfRawData_ESE_MSExchangeDatabaseInstances{}
	if err := wmi.Query("SELECT * FROM "+className(data), &data); err != nil {
		return err
	}
	for _, db := range data {
		labelName := toLabelName(db.Name)
		if strings.HasSuffix(labelName, "_total") {
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			c.IODatabaseReadsAverageLatency,
			prometheus.GaugeValue,
			float64(db.IODatabaseReadsAverageLatency),
			labelName,
		)
		ch <- prometheus.MustNewConstMetric(
			c.IODatabaseWritesAverageLatency,
			prometheus.GaugeValue,
			float64(db.IODatabaseWritesAverageLatency),
			labelName,
		)
		ch <- prometheus.MustNewConstMetric(
			c.IOLogWritesAverageLatency,
			prometheus.GaugeValue,
			float64(db.IOLogWritesAverageLatency),
			labelName,
		)
		ch <- prometheus.MustNewConstMetric(
			c.IODatabaseReadsRecoveryAverageLatency,
			prometheus.GaugeValue,
			float64(db.IODatabaseReadsRecoveryAverageLatency),
			labelName,
		)
		ch <- prometheus.MustNewConstMetric(
			c.IODatabaseWritesRecoveryAverageLatency,
			prometheus.GaugeValue,
			float64(db.IODatabaseWritesRecoveryAverageLatency),
			labelName,
		)
	}
	return nil
}

func (c *exchangeCollector) collectHTTPProxy(ch chan<- prometheus.Metric) error {
	data := []win32_PerfRawData_MSExchangeHttpProxy_MSExchangeHttpProxy{}
	if err := wmi.Query("SELECT * FROM "+className(data), &data); err != nil {
		return err
	}
	for _, proxy := range data {
		labelName := toLabelName(proxy.Name)
		ch <- prometheus.MustNewConstMetric(
			c.MailboxServerLocatorAverageLatency,
			prometheus.GaugeValue,
			float64(proxy.MailboxServerLocatorAverageLatency),
			labelName,
		)
		ch <- prometheus.MustNewConstMetric(
			c.AverageAuthenticationLatency,
			prometheus.GaugeValue,
			float64(proxy.AverageAuthenticationLatency),
			labelName,
		)
		ch <- prometheus.MustNewConstMetric(
			c.AverageClientAccessServerProcessingLatency,
			prometheus.GaugeValue,
			float64(proxy.AverageClientAccessServerProcessingLatency),
			labelName,
		)
		ch <- prometheus.MustNewConstMetric(
			c.MailboxServerProxyFailureRate,
			prometheus.GaugeValue,
			float64(proxy.MailboxServerProxyFailureRate),
			labelName,
		)
		ch <- prometheus.MustNewConstMetric(
			c.OutstandingProxyRequests,
			prometheus.GaugeValue,
			float64(proxy.OutstandingProxyRequests),
			labelName,
		)
		ch <- prometheus.MustNewConstMetric(
			c.ProxyRequestsPerSec,
			prometheus.GaugeValue,
			float64(proxy.ProxyRequestsPerSec),
			labelName,
		)
	}
	return nil
}

func (c *exchangeCollector) collectActiveSync(ch chan<- prometheus.Metric) error {
	data := []win32_PerfRawData_MSExchangeActiveSync_MSExchangeActiveSync{}
	if err := wmi.Query("SELECT * FROM "+className(data), &data); err != nil {
		return err
	}
	for _, acsync := range data {
		ch <- prometheus.MustNewConstMetric(
			c.ActiveSyncRequestsPerSec,
			prometheus.GaugeValue,
			float64(acsync.RequestsPerSec),
		)
		ch <- prometheus.MustNewConstMetric(
			c.PingCommandsPending,
			prometheus.GaugeValue,
			float64(acsync.PingCommandsPending),
		)
		ch <- prometheus.MustNewConstMetric(
			c.SyncCommandsPerSec,
			prometheus.GaugeValue,
			float64(acsync.SyncCommandsPerSec),
		)
	}
	return nil
}

func (c *exchangeCollector) collectAvailService(ch chan<- prometheus.Metric) error {
	data := []win32_PerfRawData_MSExchangeAvailabilityService_MSExchangeAvailabilityService{}
	if err := wmi.Query("SELECT * FROM "+className(data), &data); err != nil {
		return err
	}
	for _, availservice := range data {
		ch <- prometheus.MustNewConstMetric(
			c.AvailabilityRequestsSec,
			prometheus.GaugeValue,
			float64(availservice.RequestsSec),
		)
	}
	return nil
}

func (c *exchangeCollector) collectOWA(ch chan<- prometheus.Metric) error {
	data := []win32_PerfRawData_MSExchangeOWA_MSExchangeOWA{}
	if err := wmi.Query("SELECT * FROM "+className(data), &data); err != nil {
		return err
	}
	for _, owa := range data {
		ch <- prometheus.MustNewConstMetric(
			c.CurrentUniqueUsers,
			prometheus.GaugeValue,
			float64(owa.CurrentUniqueUsers),
		)
		ch <- prometheus.MustNewConstMetric(
			c.OWARequestsPerSec,
			prometheus.GaugeValue,
			float64(owa.RequestsPerSec),
		)
	}
	return nil
}

func (c *exchangeCollector) collectAutoDiscover(ch chan<- prometheus.Metric) error {
	data := []win32_PerfRawData_MSExchangeAutodiscover_MSExchangeAutodiscover{}
	if err := wmi.Query("SELECT * FROM "+className(data), &data); err != nil {
		return err
	}
	for _, autodisc := range data {
		ch <- prometheus.MustNewConstMetric(
			c.AutodiscoverRequestsPerSec,
			prometheus.GaugeValue,
			float64(autodisc.RequestsPerSec),
		)
	}
	return nil
}

func (c *exchangeCollector) collectMgmtWorkloads(ch chan<- prometheus.Metric) error {
	data := []win32_PerfRawData_MSExchangeWorkloadManagementWorkloads_MSExchangeWorkloadManagementWorkloads{}
	if err := wmi.Query("SELECT * FROM "+className(data), &data); err != nil {
		return err
	}
	ch <- prometheus.MustNewConstMetric(
		c.ActiveTasks,
		prometheus.GaugeValue,
		float64(data[0].ActiveTasks),
	)
	ch <- prometheus.MustNewConstMetric(
		c.CompletedTasks,
		prometheus.CounterValue,
		float64(data[0].CompletedTasks),
	)
	ch <- prometheus.MustNewConstMetric(
		c.QueuedTasks,
		prometheus.CounterValue,
		float64(data[0].QueuedTasks),
	)
	return nil
}

func (c *exchangeCollector) collectRPCClientAccess(ch chan<- prometheus.Metric) error {
	data := []win32_PerfRawData_MSExchangeRpcClientAccess_MSExchangeRpcClientAccess{}
	if err := wmi.Query("SELECT * FROM "+className(data), &data); err != nil {
		return err
	}
	for _, rpc := range data {
		ch <- prometheus.MustNewConstMetric(
			c.RPCAveragedLatency,
			prometheus.GaugeValue,
			float64(rpc.RPCAveragedLatency),
		)
		ch <- prometheus.MustNewConstMetric(
			c.RPCRequests,
			prometheus.CounterValue,
			float64(rpc.RPCRequests),
		)
		ch <- prometheus.MustNewConstMetric(
			c.ActiveUserCount,
			prometheus.GaugeValue,
			float64(rpc.ActiveUserCount),
		)
		ch <- prometheus.MustNewConstMetric(
			c.ConnectionCount,
			prometheus.CounterValue,
			float64(rpc.ConnectionCount),
		)
		ch <- prometheus.MustNewConstMetric(
			c.RPCOperationsPerSec,
			prometheus.GaugeValue,
			float64(rpc.RPCOperationsPerSec),
		)
		ch <- prometheus.MustNewConstMetric(
			c.UserCount,
			prometheus.CounterValue,
			float64(rpc.UserCount),
		)
	}
	return nil
}
