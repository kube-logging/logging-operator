---
title: Null
weight: 200
generated_file: true
---

# Null output plugin for Fluentd
## Overview


For details, see [https://docs.fluentd.org/output/null](https://docs.fluentd.org/output/null).

## Example output configurations

```yaml
spec:
  nullout:
    never_flush: false
```


## Configuration
## NullOutputConfig

### never_flush (*bool, optional) {#nulloutputconfig-never_flush}

The parameter for testing to simulate the output plugin that never succeeds to flush. 



