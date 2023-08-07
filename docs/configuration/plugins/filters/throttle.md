# [Throttle Filter](https://github.com/rubrikinc/fluent-plugin-throttle)
## Overview
 A sentry plugin to throttle logs. Logs are grouped by a configurable key. When a group exceeds a configuration rate, logs are dropped for this group.

## Configuration
## Throttle

### group_key (string, optional) {#throttle-group_key}

Used to group logs. Groups are rate limited independently  

Default:  kubernetes.container_name

### group_bucket_period_s (int, optional) {#throttle-group_bucket_period_s}

This is the period of of time over which group_bucket_limit applies  

Default:  60

### group_bucket_limit (int, optional) {#throttle-group_bucket_limit}

Maximum number logs allowed per groups over the period of group_bucket_period_s  

Default:  6000

### group_drop_logs (bool, optional) {#throttle-group_drop_logs}

When a group reaches its limit, logs will be dropped from further processing if this value is true  

Default:  true

### group_reset_rate_s (int, optional) {#throttle-group_reset_rate_s}

After a group has exceeded its bucket limit, logs are dropped until the rate per second falls below or equal to group_reset_rate_s.  

Default:  group_bucket_limit/group_bucket_period_s

### group_warning_delay_s (int, optional) {#throttle-group_warning_delay_s}

When a group reaches its limit and as long as it is not reset, a warning message with the current log rate of the group is emitted repeatedly. This is the delay between every repetition.  

Default:  10 seconds


 ## Example `Throttle` filter configurations
 ```yaml
 apiVersion: logging.banzaicloud.io/v1beta1
 kind: Flow
 metadata:

	name: demo-flow

 spec:

	filters:
	  - throttle:
	      group_key: "$.kubernetes.container_name"
	selectors: {}
	localOutputRefs:
	  - demo-output

 ```

 #### Fluentd Config Result
 ```yaml
 <filter **>

	@type throttle
	@id test_throttle
	group_key $.kubernetes.container_name

 </filter>
 ```

---
