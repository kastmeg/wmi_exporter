// +build windows

package collector

import (
	"github.com/StackExchange/wmi"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	subsystem string = "exchange"
	MSE2019   exver  = 2019
	MSE2016   exver  = 2016
	MSE2013   exver  = 2013
)

var (
	argOverrideVersion = kingpin.Flag("collector.exchange.version", "Override exchange version (2016, 2013, ...)").String()
	argFoobar          = kingpin.Flag("collector.exchange.foobar", "foobar").String()
)

type exver int

// ExchangeCollector collects Exchange metrics
type ExchangeCollector struct {
	Version exver

	PercentProcessorTime                    *prom.Desc
	PercentUserTime                         *prom.Desc
	PercentPrivilegedTime                   *prom.Desc
	ProcessorQueueLength                    *prom.Desc
	AvailableMbytes                         *prom.Desc
	PoolPagedBytes                          *prom.Desc
	TransitionPagesRePurposedPerSec         *prom.Desc
	PageReadsPerSec                         *prom.Desc
	PagesPerSec                             *prom.Desc
	PagesInputPerSec                        *prom.Desc
	PagesOutputPerSec                       *prom.Desc
	FreeSpace                               *prom.Desc
	AvgDiskSecPerRead                       *prom.Desc
	AvgDiskSecPerWrite                      *prom.Desc
	AvgDiskSecPerTransfer                   *prom.Desc
	PacketsOutboundErrors                   *prom.Desc
	ConnectionsResetTCPv4                   *prom.Desc
	ConnectionsResetTCPv6                   *prom.Desc
	ApplicationRestarts                     *prom.Desc
	WorkerProcessRestarts                   *prom.Desc
	RequestsCurrent                         *prom.Desc
	RequestWaitTime                         *prom.Desc
	RequestsInApplicationQueue              *prom.Desc
	LDAPReadTime                            *prom.Desc
	LDAPSearchTime                          *prom.Desc
	LDAPSearchesTimedOutPerMin              *prom.Desc
	LongRunningLDAPOperationsPerMin         *prom.Desc
	LDAPSearchesTimedLimitExceededPerMin    *prom.Desc
	RPCAveragedLatency                      *prom.Desc
	RPCRequests                             *prom.Desc
	ActiveUserCount                         *prom.Desc
	ConnectionCount                         *prom.Desc
	RPCOperationsPerSec                     *prom.Desc
	UserCount                               *prom.Desc
	ExternalActiveRemoteDeliveryQueueLength *prom.Desc
	InternalActiveRemoteDeliveryQueueLength *prom.Desc
	ActiveMailboxDeliveryQueueLength        *prom.Desc
	RetryMailboxDeliveryQueueLength         *prom.Desc
	UnreachableQueueLength                  *prom.Desc
	ExternalLargestDeliveryQueueLength      *prom.Desc
	InternalLargestDeliveryQueueLength      *prom.Desc
	PoisonQueueLength                       *prom.Desc
	MailboxServerLocatorAverageLatency      *prom.Desc
	AverageAuthenticationLatency            *prom.Desc
	MailboxServerProxyFailureRate           *prom.Desc
	OutstandingProxyRequests                *prom.Desc
	ProxyRequestsPerSec                     *prom.Desc
	RequestsPerSec                          *prom.Desc
	PingCommandsPending                     *prom.Desc
	SyncCommandsPerSec                      *prom.Desc
	AvailabilityRequestssec                 *prom.Desc
	CurrentUniqueUsers                      *prom.Desc
	ActiveTasks                             *prom.Desc
	CompletedTasks                          *prom.Desc
	QueuedTasks                             *prom.Desc
	IODatabaseReadsAverageLatency           *prom.Desc
	IODatabaseWritesAverageLatency          *prom.Desc
	IOLogWritesAverageLatency               *prom.Desc
	IODatabaseReadsRecoveryAverageLatency   *prom.Desc
	IODatabaseWritesRecoveryAverageLatency  *prom.Desc
}

func init() {
	Factories[subsystem] = NewExchangeCollector
}

// Collect sends the metric values for each metric to the provided prometheus Metric channel
func (c *ExchangeCollector) Collect(ctx *ScrapeContext, ch chan<- prom.Metric) error {
	if desc, err := c.collect(ch); err != nil {
		log.Error("failed collecting Exchange metrics:", desc, err)
		return err
	}
	return nil
}

// NewExchangeCollector creates a new ExchangeCollector
func NewExchangeCollector() (Collector, error) {
	exCollector := ExchangeCollector{
		PercentProcessorTime: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "cpu_time_total"),
			"Total percent processor time", []string{"mode"}, nil,
		),
		PercentUserTime: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "cpu_time_total"),
			"Total percent user time", []string{"mode"}, nil,
		),
		PercentPrivilegedTime: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "cpu_time_total"),
			"Total percent privilege time", []string{"mode"}, nil,
		),
		ProcessorQueueLength: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", []string{}, nil,
		),
		LDAPReadTime: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		LDAPSearchTime: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		LDAPSearchesTimedOutPerMin: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		LongRunningLDAPOperationsPerMin: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		LDAPSearchesTimedLimitExceededPerMin: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		AvailableMbytes: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		PoolPagedBytes: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		TransitionPagesRePurposedPerSec: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		PageReadsPerSec: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		PagesPerSec: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		PagesInputPerSec: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		PagesOutputPerSec: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		FreeSpace: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		AvgDiskSecPerRead: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		AvgDiskSecPerWrite: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		AvgDiskSecPerTransfer: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		PacketsOutboundErrors: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		ConnectionsResetTCPv4: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		ConnectionsResetTCPv6: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		RPCAveragedLatency: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		RPCRequests: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		ActiveUserCount: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		ConnectionCount: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		RPCOperationsPerSec: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		UserCount: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		ExternalActiveRemoteDeliveryQueueLength: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		InternalActiveRemoteDeliveryQueueLength: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		ActiveMailboxDeliveryQueueLength: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		RetryMailboxDeliveryQueueLength: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		UnreachableQueueLength: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		ExternalLargestDeliveryQueueLength: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		InternalLargestDeliveryQueueLength: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		PoisonQueueLength: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		ApplicationRestarts: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		WorkerProcessRestarts: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		RequestsCurrent: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		RequestWaitTime: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		RequestsInApplicationQueue: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		MailboxServerLocatorAverageLatency: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		AverageAuthenticationLatency: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		MailboxServerProxyFailureRate: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		OutstandingProxyRequests: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		ProxyRequestsPerSec: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		RequestsPerSec: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		PingCommandsPending: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		SyncCommandsPerSec: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		AvailabilityRequestssec: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		CurrentUniqueUsers: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		ActiveTasks: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		CompletedTasks: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		QueuedTasks: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		IODatabaseReadsAverageLatency: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		IODatabaseWritesAverageLatency: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		IOLogWritesAverageLatency: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		IODatabaseReadsRecoveryAverageLatency: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
		IODatabaseWritesRecoveryAverageLatency: prom.NewDesc(
			prom.BuildFQName(Namespace, subsystem, "foo"),
			"(descr)", nil, nil,
		),
	}

	return &exCollector, nil
}

func (c *ExchangeCollector) collect(ch chan<- prom.Metric) (*prom.Desc, error) {
	var dstADAccess []Win32_PerfRawData_MSExchangeADAccess_MSExchangeADAccessProcesses
	if err := wmi.Query(queryAll(dstADAccess), &dstADAccess); err != nil {
		return nil, err
	}
	for _, app := range dstADAccess {
		ch <- prom.MustNewConstMetric(
			c.LDAPReadTime,
			prom.GaugeValue,
			float64(app.LDAPReadTime),
		)
		ch <- prom.MustNewConstMetric(
			c.LDAPSearchTime,
			prom.GaugeValue,
			float64(app.LDAPSearchTime),
		)
		ch <- prom.MustNewConstMetric(
			c.LDAPSearchesTimedOutPerMin,
			prom.GaugeValue,
			float64(app.LongRunningLDAPOperationsPerMin),
		)
		ch <- prom.MustNewConstMetric(
			c.LongRunningLDAPOperationsPerMin,
			prom.GaugeValue,
			float64(app.LDAPSearchesTimedOutPerMin),
		)
		ch <- prom.MustNewConstMetric(
			c.LDAPSearchesTimedLimitExceededPerMin,
			prom.GaugeValue,
			float64(app.LDAPSearchesTimedLimitExceededPerMin),
		)
	}

	return nil, nil
}

// All WMI Class definitions have the wrong type  !!
type placeholder uint64

//
// Exchange WMI Classes
//
type Win32_PerfRawData_MSExchangeADAccess_MSExchangeADAccessProcesses struct {
	Name string

	LDAPReadTime                         uint64
	LDAPSearchTime                       uint64
	LDAPSearchesTimedOutPerMin           uint64
	LongRunningLDAPOperationsPerMin      uint64
	LDAPSearchesTimedLimitExceededPerMin uint64
}

type Win32_PerfRawData_MSExchangeRpcClientAccess_MSExchangeRpcClientAccess struct {
	Name string

	RPCAveragedLatency  placeholder
	RPCRequests         placeholder
	ActiveUserCount     placeholder
	ConnectionCount     placeholder
	RPCOperationsPerSec placeholder
	UserCount           placeholder
}

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
type Win32_PerfRawData_MSExchangeHttpProxy_MSExchangeHttpProxy struct {
	Name string

	MailboxServerLocatorAverageLatency placeholder
	AverageAuthenticationLatency       placeholder
	MailboxServerProxyFailureRate      placeholder
	OutstandingProxyRequests           placeholder
	ProxyRequestsPerSec                placeholder
	RequestsPerSec                     placeholder
}

type Win32_PerfRawData_MSExchangeActiveSync_MSExchangeActiveSync struct {
	Name string

	RequestsPerSec      placeholder
	PingCommandsPending placeholder
	SyncCommandsPerSec  placeholder
}

type Win32_PerfRawData_MSExchangeAvailabilityService_MSExchangeAvailabilityService struct {
	Name string

	AvailabilityRequestssec placeholder
}

type Win32_PerfRawData_MSExchangeOWA_MSExchangeOWA struct {
	Name string

	CurrentUniqueUsers placeholder
	RequestsPerSec     placeholder
}

type Win32_PerfRawData_MSExchangeAutodiscover_MSExchangeAutodiscover struct {
	Name string

	RequestsPerSec placeholder
}

type Win32_PerfRawData_MSExchangeWorkloadManagementWorkloads_MSExchangeWorkloadManagementWorkloads struct {
	Name string

	ActiveTasks    placeholder
	CompletedTasks placeholder
	QueuedTasks    placeholder
}

type Win32_PerfRawData_ESE_MSExchangeDatabaseInstances struct {
	Name string

	IODatabaseReadsAverageLatency          placeholder
	IODatabaseWritesAverageLatency         placeholder
	IOLogWritesAverageLatency              placeholder
	IODatabaseReadsRecoveryAverageLatency  placeholder
	IODatabaseWritesRecoveryAverageLatency placeholder
}

//
// General WMI Classes
//
/*
type Win32_PerfRawData_PerfOS_Processor struct {
	Name string

	PercentProcessorTime  placeholder
	PercentUserTime       placeholder
	PercentPrivilegedTime placeholder
}

type Win32_PerfFormattedData_PerfOS_System struct {
	Name string

	ProcessorQueueLength placeholder
}

type Win32_PerfRawData_PerfOS_Memory struct {
	Name string

	AvailableMbytes                 placeholder
	PoolPagedBytes                  placeholder
	TransitionPagesRePurposedPerSec placeholder
	PageReadsPerSec                 placeholder
	PagesPerSec                     placeholder
	PagesInputPerSec                placeholder
	PagesOutputPerSec               placeholder
}

type Win32_LogicalDisk struct {
	Name string

	FreeSpace placeholder
}

type Win32_PerfRawData_PerfDisk_LogicalDisk struct {
	Name string

	AvgDiskSecPerRead     placeholder
	AvgDiskSecPerWrite    placeholder
	AvgDiskSecPerTransfer placeholder
}

type Win32_PerfRawData_Tcpip_NetworkInterface struct {
	Name string

	PacketsOutboundErrors placeholder
}

type Win32_PerfRawData_Tcpip_TCPv4 struct {
	Name string

	ConnectionsResetTCPv4 placeholder
}

type Win32_PerfRawData_Tcpip_TCPv6 struct {
	Name string

	ConnectionsResetTCPv6 placeholder
}

type Win32_PerfRawData_ASPNET_ASPNET struct {
	Name string

	ApplicationRestarts   placeholder
	WorkerProcessRestarts placeholder
	RequestsCurrent       placeholder
	RequestWaitTime       placeholder
}

type Win32_PerfRawData_ASPNET_ASPNETApplications struct {
	Name string

	RequestsInApplicationQueue placeholder
}
*/
