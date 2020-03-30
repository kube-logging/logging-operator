---
title: Model
weight: 1500
---

## Goal

Define an opinionated logical fluentd configuration model for processing kubernetes log events using go structs and render the following two representations:

- working fluentd configuration
- a configuration format that can be used for visual representation

## The model

Flow (data pipeline)

The term "data pipeline" is used in fluentd for labeled (https://docs.fluentd.org/quickstart/life-of-a-fluentd-event#labels) 
configuration sections, that should apply only to a subset of events. Also these labeled events will be skipped by the default (non-labeled) plugins. See:
https://docs.fluentd.org/configuration/routing-examples#input-greater-than-filter-greater-than-output-with-label

The flow is identified by a kubernetes namespace and/or a set of labels. The non-labeled non-namespaced events can be considered as the "global" flow that processes all events.

### Components

Each flow has the following fluentd components:

- inside the label section:

    - zero or more sequential filters
    - one or more outputs

- outside the label section (in the global section:

    - a router

1. Event reaches the single input source.
2. The routers (a smart match directive) for the specific flow get the event, examines it for the
 kubernetes namespace and label information and reemits it with a fluentd label when they found a match.
  The global router matches all events.
3. The event arrives to the correct flow based on it's label. Filters are then applied in sequential order and finally the flow's outputs will be the end of the event's journey.

    ```bash
            global router -> global flow [filter, ...] => output[, output, ...]
          / 
    input - router A      -> flow A      [filter, ...] => output[, output, ...]
          \
            router B      -> flow B      [filter, ...] => output[, output, ...]
            .
            .
            .
    ```