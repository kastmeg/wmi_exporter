# exchange collector

The exchange collector exposes metrics about the MS Exchange server

|||
-|-
Metric name prefix  | `exchange`
Classes 			| [Win32_PerfRawData_MSExchangeADAccess_MSExchangeADAccessProcesses](https://docs.microsoft.com/en-us/exchange/)<br/> [Win32_PerfRawData_MSExchangeTransportQueues_MSExchangeTransportueues](https://docs.microsoft.com/en-us/exchange/)<br/> [Win32_PerfRawData_ESE_MSExchangeDatabaseInstances](https://docs.microsoft.com/en-us/exchange/)<br/> [Win32_PerfRawData_MSExchangeHttpProxy_MSExchangeHttpProxy](https://docs.microsoft.com/en-us/exchange/)<br/> [Win32_PerfRawData_MSExchangeActiveSync_MSExchangeActiveSync](https://docs.microsoft.com/en-us/exchange/)<br/> [Win32_PerfRawData_MSExchangeAvailabilityService_MSExchangeAvailabilityService](https://docs.microsoft.com/en-us/exchange/)<br/> [Win32_PerfRawData_MSExchangeOWA_MSExchangeOWA](https://docs.microsoft.com/en-us/exchange/)<br/> [Win32_PerfRawData_MSExchangeAutodiscover_MSExchangeAutodiscover](https://docs.microsoft.com/en-us/exchange/)<br/> [Win32_PerfRawData_MSExchangeWorkloadManagementWorkloads_MSExchangeWorkloadManagementWorkloads](https://docs.microsoft.com/en-us/exchange/)<br/> [Win32_PerfRawData_MSExchangeRpcClientAccess_MSExchangeRpcClientAccess](https://docs.microsoft.com/en-us/exchange/)<br/>
Enabled by default? | No

## Flags
Since the official WMI class names are extremely long, the following shorthands are used instead.
`--collectors.exchange.class-list` lists these shorthands along with the real WMI class names

* `ad_access_procs`
* `transport_queues`
* `database_instances`
* `http_proxy`
* `active_sync`
* `availability_service`
* `owa`
* `auto_descover`
* `management_workloads`
* `rpc_client_access`

### `--collectors.exchange.list`
List all available MS Exchange WMI class long- and short-names

### `--collectors.exchange.classes-enabled`
Comma-separated list of exchange WMI classes from which the exporter should collect data. 
If no classes are given, all classes will be queried.

## Metrics

Name | Description | Type | Labels
-----|-------------|------|-------
`wmi_exchange_ldap_read_time` | LDAP Read Time | gauge | name
`wmi_exchange_ldap_search_time` | LDAP Search Time | gauge | name
`wmi_exchange_ldap_timeout_errors_per_sec` | LDAP timeout errors per second | gauge | name
`wmi_exchange_long_running_ldap_operations_per_min` | Long Running LDAP operations pr minute | gauge | name
`wmi_exchange_ldap_searches_time_limit_exceeded_per_min` | LDAP searches time limit exceeded per minute | gauge | name
`wmi_exchange_external_active_remote_delivery_queue_len` | External Active Remote Delivery Queue length | gauge | name
`wmi_exchange_internal_active_remote_delivery_queue_len` | Internal Active Remote Delivery Queue length | gauge | name
`wmi_exchange_active_mailbox_delivery_queue_len` | Active Mailbox Delivery Queue length | gauge| name
`wmi_exchange_retry_mailbox_delivery_queue_len` | Retry Mailbox Delivery Queue length | gauge | name
`wmi_exchange_unreachable_queue_len` | Unreachable Queue length | gauge | name
`wmi_exchange_external_largest_delivery_queue_len` | External Largest Delivery Queue length | gauge | name
`wmi_exchange_inernal_largest_delivery_queue_len` | Internal Largest Delivery Queue length | gauge | name
`wmi_exchange_poison_queue_len` | Poison Queue length | gauge | name
`wmi_exchange_io_database_reads_average_latency` | Average database read latency | gauge | name
`wmi_exchange_io_database_writes_average_latency` | Average database write latency | gauge | name
`wmi_exchange_io_log_writes_average_latency` | Average Log Writes Latency | gauge | name
`wmi_exchange_io_database_reads_recovery_average_latency` | Database reads recovery avrage latency | gauge | name
`wmi_exchange_io_database_writes_recovery_average_latency` | Database writes recovery average latency | gauge | name
`wmi_exchange_mailbox_server_locator_average_latency` | Exchange HTTP Proxy Mailbox Server Locator latency (avg) | gauge | name
`wmi_exchange_average_authentication_latency` | Exchange HTTP Proxy Authentication Latency (avg) | gauge | name
`wmi_exchange_average_client_access_server_processing_latency` | Exchange HTTP Proxy Client Access Server Processing Latency (avg) | gauge | name
`wmi_exchange_mailbox_server_proxy_failure_rate` | Exchange HTTP Proxy Mailbox Server Proxy Failure Rate | gauge | name
`wmi_exchange_outstanding_proxy_requests` | Exchange HTTP Proxy outstanding proxy requests | gauge | name
`wmi_exchange_proxy_requests_per_sec` | Exchange HTTP Proxy requests/s | gauge | name
`wmi_exchange_active_sync_requests_per_sec` | Active Sync requests/s  | gauge | name
`wmi_exchange_ping_commands_pending` | Pending Active Sync ping-commands | gauge | name
`wmi_exchange_sync_commands_per_sec` | Active Sync sync-commands/s | gauge | name
`wmi_exchange_availability_requests_per_sec` | Availability Service / Availability requests/s | gauge | name
`wmi_exchange_current_unique_users` | Outlook Web Access current unique users | gauge | name
`wmi_exchange_owa_requests_per_sec` | Outlook Web Access requests/s | gauge | name
`wmi_exchange_autodiscover_requests_per_sec` | Autodiscovery requests/s | gauge | name
`wmi_exchange_active_tasks` | Active Workload Management Tasks | gauge | name
`wmi_exchange_completed_tasks` | Completed Workload Management Tasks | counter | name
`wmi_exchange_queued_tasks` | Queued Workload Management Tasks | counter | name
`wmi_exchange_rpc_averaged_latency` | RPC Client Access averaged latency | counter | name
`wmi_exchange_rpc_requests` | RPC Client Access requests | counter | name
`wmi_exchange_active_user_count` | RPC Client Access active user count | counter | name
`wmi_exchange_connection_count` | RPC Client Access connection count | counter | name
`wmi_exchange_rpc_operations_per_sec` | RPC Client Access operations per sec | gauge | name
`wmi_exchange_user_count` | RPC Client Access user count | counter | name

### Example metric
_This collector does not yet have explained examples, we would appreciate your help adding them!_

## Useful queries
_This collector does not yet have any useful queries added, we would appreciate your help adding them!_

## Alerting examples
_This collector does not yet have alerting examples, we would appreciate your help adding them!_

