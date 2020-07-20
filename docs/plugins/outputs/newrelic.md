---
title: NewRelic
weight: 200
generated_file: true
---

# New Relic Logs plugin for Fluentd
## Overview
**newrelic** output plugin send log data to New Relic Logs

 #### Example output configurations
 ```
 spec:
   newrelic:
     license_key:
       valueFrom:
         secretKeyRef:
           name: logging-newrelic
           key: licenseKey
 ```

## Configuration
### Output Config
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| api_key | *secret.Secret | No | - | New Relic API Insert key<br>[Secret](../secret/)<br> |
| license_key | *secret.Secret | No | - | New Relic License Key (recommended)<br>[Secret](../secret/"<br>LicenseKey *secret.Secret `json:"license_key)`<br> |
| base_uri | string | No | https://log-api.newrelic.com/log/v1 | New Relic ingestion endpoint<br>[Secret](../secret/)<br> |
