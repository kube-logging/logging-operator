# [Throttle Filter](https://github.com/rubrikinc/fluent-plugin-throttle)
## Overview
 A sentry plugin to throttle logs. Logs are grouped by a configurable key. When a group exceeds a configuration rate, logs are dropped for this group.

## Configuration
### Throttle
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| group_key | string | No |  kubernetes.container_name | Used to group logs. Groups are rate limited independently <br> |
| group_bucket_period_s | int | No |  60 | This is the period of of time over which group_bucket_limit applies <br> |
| group_bucket_limit | int | No |  6000 | Maximum number logs allowed per groups over the period of group_bucket_period_s <br> |
| group_drop_logs | bool | No |  true | When a group reaches its limit, logs will be dropped from further processing if this value is true <br> |
| group_reset_rate_s | int | No |  group_bucket_limit/group_bucket_period_s | After a group has exceeded its bucket limit, logs are dropped until the rate per second falls below or equal to group_reset_rate_s. <br> |
| group_warning_delay_s | int | No |  10 seconds | When a group reaches its limit and as long as it is not reset, a warning message with the current log rate of the group is emitted repeatedly. This is the delay between every repetition. <br> |
 #### Example `Throttle` filter configurations
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
  outputRefs:
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
