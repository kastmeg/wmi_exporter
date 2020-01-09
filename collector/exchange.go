// +build windows

package collector

import (
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	prom "github.com/prometheus/client_golang/prometheus"
)

const (
	Subsystem  string = "exchange"
	MSE2019 exver  = 2019
	MSE2016 exver  = 2016
	MSE2013 exver  = 2013
)

var (
	argOverrideVersion = kingpin.Flag("collector.exchange.version", "Override exchange version (2016, 2013, ...)")
	argFoobar = kingpin.Flag("collector.exchange.foobar", "Foorbar baz...)")
	argBarfoo = kingpin.Flag("collector.exchange.barfoo", "Barfoo fight")
)

type exver int

// ExchangeCollector is a prom.collector for Exchange metrics
type ExchangeCollector struct {
	// MajorVersion (defaults to Ex2016 for now)
	Version exver

	// Win32_PerfRawData_PerfOS_Processor
	PercentProcessorTime  *prom.Desc
	PercentUserTime       *prom.Desc
	PercentPrivilegedTime *prom.Desc

	// Win32_PerfFormattedData_PerfOS_System
	ProcessorQueueLength *prom.Desc

	// Win32_PerfRawData_PerfOS_Memory
	AvailableMbytes                 *prom.Desc
	PoolPagedBytes                  *prom.Desc
	TransitionPagesRePurposedPerSec *prom.Desc
	PageReadsPerSec                 *prom.Desc
	PagesPerSec                     *prom.Desc
	PagesInputPerSec                *prom.Desc
	PagesOutputPerSec               *prom.Desc

	// Win32_LogicalDisk
	FreeSpace *prom.Desc

	// Win32_PerfRawData_PerfDisk_LogicalDisk
	AvgDiskSecPerRead     *prom.Desc
	AvgDiskSecPerWrite    *prom.Desc
	AvgDiskSecPerTransfer *prom.Desc

	// Win32_PerfRawData_Tcpip_NetworkInterface
	PacketsOutboundErrors *prom.Desc

	// Win32_PerfRawData_Tcpip_TCPv4
	ConnectionsResetTCPv4 *prom.Desc

	// Win32_PerfRawData_Tcpip_TCPv6
	ConnectionsResetTCPv6 *prom.Desc

	// Win32_PerfRawData_ASPNET_ASPNET
	ApplicationRestarts   *prom.Desc
	WorkerProcessRestarts *prom.Desc
	RequestsCurrent       *prom.Desc
	RequestWaitTime       *prom.Desc

	// Win32_PerfRawData_ASPNET_ASPNETApplications
	RequestsInApplicationQueue *prom.Desc

	// Win32_PerfFormattedData_MSExchangeADAccess_MSExchangeADAccessDomainControllers
	// Win32_PerfRawData_MSExchangeADAccess_MSExchangeADAccessProcesses
	LDAPReadTime                         *prom.Desc
	LDAPSearchTime                       *prom.Desc
	LDAPSearchesTimedOutPerMin           *prom.Desc
	LongRunningLDAPOperationsPerMin      *prom.Desc
	LDAPSearchesTimedLimitExceededPerMin *prom.Desc


	// Win32_PerfRawData_MSExchangeRpcClientAccess_MSExchangeRpcClientAccess
	RPCAveragedLatency  *prom.Desc
	RPCRequests         *prom.Desc
	ActiveUserCount     *prom.Desc
	ConnectionCount     *prom.Desc
	RPCOperationsPerSec *prom.Desc
	UserCount           *prom.Desc

	// Win32_PerfRawData_MSExchangeTransportQueues_MSExchangeTransportQueues
	ExternalActiveRemoteDeliveryQueueLength *prom.Desc
	InternalActiveRemoteDeliveryQueueLength *prom.Desc
	ActiveMailboxDeliveryQueueLength        *prom.Desc
	RetryMailboxDeliveryQueueLength         *prom.Desc
	UnreachableQueueLength                  *prom.Desc
	ExternalLargestDeliveryQueueLength      *prom.Desc
	InternalLargestDeliveryQueueLength      *prom.Desc
	PoisonQueueLength                       *prom.Desc

	// Win32_PerfRawData_MSExchangeHttpProxy_MSExchangeHttpProxy
	MailboxServerLocatorAverageLatency *prom.Desc
	AverageAuthenticationLatency       *prom.Desc
	MailboxServerProxyFailureRate      *prom.Desc
	OutstandingProxyRequests           *prom.Desc
	ProxyRequestsPerSec                *prom.Desc
	RequestsPerSec                     *prom.Desc

	// Win32_PerfRawData_MSExchangeActiveSync_MSExchangeActiveSync
	RequestsPerSec      *prom.Desc
	PingCommandsPending *prom.Desc
	SyncCommandsPerSec  *prom.Desc

	// Win32_PerfRawData_MSExchangeAvailabilityService_MSExchangeAvailabilityService
	AvailabilityRequestssec *prom.Desc

	// Win32_PerfRawData_MSExchangeOWA_MSExchangeOWA
	CurrentUniqueUsers *prom.Desc
	RequestsPerSec     *prom.Desc

	// Win32_PerfRawData_MSExchangeAutodiscover_MSExchangeAutodiscover
	RequestsPerSec *prom.Desc

	// Win32_PerfRawData_MSExchangeWorkloadManagementWorkloads_MSExchangeWorkloadManagementWorkloads
	ActiveTasks    *prom.Desc
	CompletedTasks *prom.Desc
	QueuedTasks    *prom.Desc

	// Win32_PerfRawData_ESE_MSExchangeDatabaseInstances
	IODatabaseReadsAverageLatency          *prom.Desc
	IODatabaseWritesAverageLatency         *prom.Desc
	IOLogWritesAverageLatency              *prom.Desc
	IODatabaseReadsRecoveryAverageLatency  *prom.Desc
	IODatabaseWritesRecoveryAverageLatency *prom.Desc
}

func init() {
	Factories[Subsystem] = NewExchangeCollector
}

// NewExchangeCollector creates a new ExchangeCollector 
func NewExchangeCollector() (ExchangeCollector, error) {
	exCollector := ExchangeCollector{
		PercentProcessorTime: prom.NewDesc(
			prom.BuildFQName(Namespace, Subsystem, "cpu_time_total"),
			"Total percent processor time", []string{"mode"}, nil,
		),
		PercentUserTime: prom.NewDesc(
			prom.BuildFQName(Namespace, Subsystem, "cpu_time_total"),
			"Total percent user time", []string{"mode"}, nil,
		),
		PercentPrivilegedTime: prom.NewDesc(
			prom.BuildFQName(Namespace, Subsystem, "cpu_time_total"),
			"Total percent privilege time",	[]string{"mode"}, nil,
		),
		ProcessorQueueLength: prom.NewDesc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		LDAPReadTime: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		LDAPSearchTime: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		LDAPSearchesTimedOutPerMin: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		LongRunningLDAPOperationsPerMin: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		LDAPSearchesTimedLimitExceededPerMin: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		AvailableMbytes: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		PoolPagedBytes: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		TransitionPagesRePurposedPerSec: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		PageReadsPerSec: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		PagesPerSec: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		PagesInputPerSec: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		PagesOutputPerSec: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		FreeSpace: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		AvgDiskSecPerRead: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		AvgDiskSecPerWrite: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		AvgDiskSecPerTransfer: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		PacketsOutboundErrors: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		ConnectionsResetTCPv4: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		ConnectionsResetTCPv6: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		RPCAveragedLatency: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		RPCRequests: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		ActiveUserCount: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		ConnectionCount: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		RPCOperationsPerSec: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		UserCount: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		ExternalActiveRemoteDeliveryQueueLength: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		InternalActiveRemoteDeliveryQueueLength: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		ActiveMailboxDeliveryQueueLength: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		RetryMailboxDeliveryQueueLength: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		UnreachableQueueLength: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		ExternalLargestDeliveryQueueLength: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		InternalLargestDeliveryQueueLength: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		PoisonQueueLength: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		ApplicationRestarts: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		WorkerProcessRestarts: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		RequestsCurrent: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		RequestWaitTime: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		RequestsInApplicationQueue: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		MailboxServerLocatorAverageLatency: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		AverageAuthenticationLatency: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		MailboxServerProxyFailureRate: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		OutstandingProxyRequests: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		ProxyRequestsPerSec: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		RequestsPerSec: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		RequestsPerSec: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		PingCommandsPending: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		SyncCommandsPerSec: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		AvailabilityRequestssec: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		CurrentUniqueUsers: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		RequestsPerSec: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		RequestsPerSec: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		ActiveTasks: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		CompletedTasks: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		QueuedTasks: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		IODatabaseReadsAverageLatency: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		IODatabaseWritesAverageLatency: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		IOLogWritesAverageLatency: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		IODatabaseReadsRecoveryAverageLatency: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
		IODatabaseWritesRecoveryAverageLatency: prom.Desc(
			prom.BuildFQName(Namespace, Subsystem, "foo"),
			"(descr)", []string{}, nil,
		},
	}

	return &exCollector, nil
}

// All WMI Class definitions have the wrong type  !!
type placeholder uint64

//
// Exchange WMI Classes
//
type Win32_PerfRawData_MSExchangeADAccess_MSExchangeADAccessProcesses struct {
	LDAPReadTime                         placeholder
	LDAPSearchTime                       placeholder
	LDAPSearchesTimedOutPerMin           placeholder
	LongRunningLDAPOperationsPerMin      placeholder
	LDAPSearchesTimedLimitExceededPerMin placeholder
}

type Win32_PerfRawData_MSExchangeRpcClientAccess_MSExchangeRpcClientAccess struct {
	RPCAveragedLatency  placeholder
	RPCRequests         placeholder
	ActiveUserCount     placeholder
	ConnectionCount     placeholder
	RPCOperationsPerSec placeholder
	UserCount           placeholder
}

type Win32_PerfRawData_MSExchangeTransportQueues_MSExchangeTransportQueues struct {
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
	MailboxServerLocatorAverageLatency placeholder
	AverageAuthenticationLatency       placeholder
	MailboxServerProxyFailureRate      placeholder
	OutstandingProxyRequests           placeholder
	ProxyRequestsPerSec                placeholder
	RequestsPerSec                     placeholder
}

type Win32_PerfRawData_MSExchangeActiveSync_MSExchangeActiveSync struct {
	RequestsPerSec      placeholder
	PingCommandsPending placeholder
	SyncCommandsPerSec  placeholder
}

type Win32_PerfRawData_MSExchangeAvailabilityService_MSExchangeAvailabilityService struct {
	AvailabilityRequestssec placeholder
}

type Win32_PerfRawData_MSExchangeOWA_MSExchangeOWA struct {
	CurrentUniqueUsers placeholder
	RequestsPerSec     placeholder
}

type Win32_PerfRawData_MSExchangeAutodiscover_MSExchangeAutodiscover struct {
	RequestsPerSec placeholder
}

type Win32_PerfRawData_MSExchangeWorkloadManagementWorkloads_MSExchangeWorkloadManagementWorkloads struct {
	ActiveTasks    placeholder
	CompletedTasks placeholder
	QueuedTasks    placeholder
}

type Win32_PerfRawData_ESE_MSExchangeDatabaseInstances struct {
	IODatabaseReadsAverageLatency          placeholder
	IODatabaseWritesAverageLatency         placeholder
	IOLogWritesAverageLatency              placeholder
	IODatabaseReadsRecoveryAverageLatency  placeholder
	IODatabaseWritesRecoveryAverageLatency placeholder
}

//
// General WMI Classes
//
type Win32_PerfRawData_PerfOS_Processor struct {
	PercentProcessorTime  placeholder
	PercentUserTime       placeholder
	PercentPrivilegedTime placeholder
}

type Win32_PerfFormattedData_PerfOS_System struct {
	ProcessorQueueLength placeholder
}

type Win32_PerfRawData_PerfOS_Memory struct {
	AvailableMbytes                 placeholder
	PoolPagedBytes                  placeholder
	TransitionPagesRePurposedPerSec placeholder
	PageReadsPerSec                 placeholder
	PagesPerSec                     placeholder
	PagesInputPerSec                placeholder
	PagesOutputPerSec               placeholder
}

type Win32_LogicalDisk struct {
	FreeSpace placeholder
}

type Win32_PerfRawData_PerfDisk_LogicalDisk struct {
	AvgDiskSecPerRead     placeholder
	AvgDiskSecPerWrite    placeholder
	AvgDiskSecPerTransfer placeholder
}

type Win32_PerfRawData_Tcpip_NetworkInterface struct {
	PacketsOutboundErrors placeholder
}

type Win32_PerfRawData_Tcpip_TCPv4 struct {
	ConnectionsResetTCPv4 placeholder
}

type Win32_PerfRawData_Tcpip_TCPv6 struct {
	ConnectionsResetTCPv6 placeholder
}

type Win32_PerfRawData_ASPNET_ASPNET struct {
	ApplicationRestarts   placeholder
	WorkerProcessRestarts placeholder
	RequestsCurrent       placeholder
	RequestWaitTime       placeholder
}

type Win32_PerfRawData_ASPNET_ASPNETApplications struct {
	RequestsInApplicationQueue placeholder
}

