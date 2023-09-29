# Multi-tenancy overview

Logging operator supports sevaral multi-tenant scenarios

## Multi tenancy behaviour and basic primitives

### Structural components

All container logs are aggregated into fluentd/syslogng (statefulset) by fluentbit (daemonset) without filtering or routing.

### Logical components

> Note: the notion of _flow_, _clusterflow_, _output_ and _clusteroutput_ are used interchangeably for fluentd and syslog-ng flows and outputs as well.

- Logs can only be manipulated through higher level abstractions, called _flows_ and _outputs_ which are applied at the aggregator level.
- `Flow` and `Output` are namespaced resources and can only work with logs in their own namespace, which is enforced in the aggregator.
- `ClusterFlow` can work with logs from multiple namespaces as long as the logs are available in the aggregator.
- `ClusterOutputs` can be used from a _clusterflow_ or from a namespaced _flow_.
- All the above resources are combined in a single aggregator defined by the `Logging` resource.

### Isolation options with **logging ref** and **watch namespaces**

The `loggingRef` field in each resource makes it possible to create multiple logging "domains" with different configuration.

By default, resources with an empty `loggingRef` are processed by all `logging` resources, except if the `logging` defines a
`watchNamespaces` field. In that case _flow_ and _output_ resources are only combined into a logging configuration if
they are all in the defined list of watched namespaces.

_Clusterflow_ and _clusteroutput_ resources (which are actually namespaced resources with multi-namespace capabilities) are by default
processed only in the logging control namespace, which is the namespace where the fluentd/syslogng statefulset and fluentbit daemonset are created.
Additionally there is an option to process these from all namespaces.

### Limitations

1. The collector (fluentbit) collects all logs and sends them to the aggregator (fluentd/syslog-ng). There is no option to filter or route log records based on any criteria there. 

    In case someone defines multiple tenants, they have to either
    - set up a routing logging that collects and forwards all logs to separate standalone logging aggregators
      (a logging resource that has no fluentbit collector defined).
    - set up multiple logging resources, that all have collectors, thus processing all log messages on their own.

2. WatchNamespaces is a static list.
3. Cannot define more than one fluentbit agent or cannot use a different agent for log collection.

## Extensions to the multi-tenancy toolset

### FluentbitAgent (from version 4.2)

One logging resource can have multiple FluentbitAgent resources instead of a single one, so that it supports new use cases
like the rolling upgrade of the collector and nodegroup specific collector configurations.

### WatchNamespaceSelector (from version 4.3)

`WatchNamespaceSelector` is available besides `watchNamespaces` in the logging resource to allow a dynamic, label based selection of
namespaces as well.

### LoggingRoute (from version 4.4)

Introducing a new resource called `LoggingRoute` which is a bridge between logging domains. Read more about it under [Hard multi-tenancy](#hard-multi-tenancy-with-a-logging-route)

## Multi-tenancy scenarios

### Soft multi-tenancy with a single logging resource

Users define their own namespaced resources which are merged into a single config and run by a single aggregator in a separate namespace.

Multi-namespace resources are typically defined by privileged users like system administrators. Multi-namespace resources
are good for:
- global rules (_clusterflow_ + _clusteroutput_ defined for multiple namespaces)
- sharing configuration (_clusteroutput_ is available from namespaced Flows)

Pros:
- default behaviour, easy to set up
- shared resources

Cons:
- noisy/misbehaving neighbours can cause performance issues and configuration errors
- no real visibility into configuration issues

### Soft multi-tenancy with watch namespaces

There are multiple logging resources configured in different control namespaces and the isolation boundary is controlled by
`watchNamespaces` and `watchNamespaceSelector` fields in the logging resource.

This way every logging resource represents a tenant. A tenant can define its own _flow_ and _output_ resources in the designated
namespaces. Tenants can also define _clusterflows_ and _clusteroutputs_ but has to be careful as those can process all
logs, even outside the tenant's namespaces.

Pros:
- multiple aggregators isolate tenants in terms of resource usage and configuration issues
- this setup does not need to manage loggingRefs on the resources

Cons:
- requires multiple collector daemonset running and processing the same set of logs on the nodes
- using _clusterflows_ and _clusteroutputs_ are discouraged or need to be tightly controlled
- multiple mis-configured logging resources might try to manage the same _flows_ and _outputs_

### Soft multi-tenancy with logging ref

Extends the previous scenario.

There are multiple logging resources configured in different control namespaces and the isolation boundary is also controlled by
`loggingRef` field in the logging resource (in addition to the watch namespaces).

This option is similar to the above, but extends it to more tightly control the relationship of _flow_ and _output_ resources to _loggings_.

Pros:
- prevents _flow_ and _output_ resources accidentally be managed by multiple _logging_ resources

Cons:
- in case users are in control of their own _flows_ and _outputs_, they need to make sure the `loggingRef` field is properly set

### Hard multi-tenancy with a logging route

> Warning: Experimental feature

With the introduction of a `LoggingRoute` resource, it is now possible to route logs based on namespaces to different aggregators.

A _logging route_ connects a collector (fluentbit) of one logging resource with the aggregators of other logging resources.

The collector uses the _watch namespaces_ configuration in the target logging resources to send only the logs relevant to them.
This way, target loggings can effectively work in isolation on their own logs and nothing more.

> Note: This method can be combined with a _logging ref_ as well, similar to how the [Soft multi-tenancy with logging ref](#soft-multi-tenancy-with-logging-ref) 

Pros:
- the log routing happens on the collector level so there is no possibility for a tenant to accidentally mess with logs outside
of its logging domain

Cons:
- one collector agent is now responsible to handle multiple output queues and failure scenarios, which fluentbit does not handle well by default
