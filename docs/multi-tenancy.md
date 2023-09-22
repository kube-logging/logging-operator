# Multi-tenancy overview

Logging operator supports different multi-tenant scenarios

## Multi tenancy primitives

Available in logging operator from 2019:
- Namespaced resources (Flow, Output) and multi-namespace resources (ClusterFlows, ClusterOutputs) are merged into a single forwarder configuration.
- Namespaced Flows can refer to ClusterOutputs.
- Namespaced resources are limited by `watchNamespaces` in the logging spec.
- Multi-namespace resources are limited to the control namespace, or optionally can be watched in all namespaces (this feature was added later).
- The logging ref makes it possible to isolate multiple logging systems.

## Options

### Soft multi-tenancy with a single logging resource

Users define their own namespaced resources which are merged into a single config and run by a single forwarder.
Multi-namespace resources are typically defined by higher privilege users like system administrators for global rules or 
to reuse a single ouptut configuration between all teams.

Pros:
- default behaviour, easy to set up
- shares resources

Cons:
- noisy neighbour

### Soft multi-tenancy with isolated forwarders `watchNamespaces` isolation

There are multiple logging resources configured in different control namespaces where the namespaced
resources (Flow/Output) are isolated by `watchNamespaces`
