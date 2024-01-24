---
title: FluentdConfig
weight: 200
generated_file: true
---

## FluentdConfig

###  (metav1.TypeMeta, required) {#fluentdconfig-}


### metadata (metav1.ObjectMeta, optional) {#fluentdconfig-metadata}


### spec (FluentdSpec, optional) {#fluentdconfig-spec}


### status (FluentdConfigStatus, optional) {#fluentdconfig-status}



## FluentdConfigStatus

### active (*bool, optional) {#fluentdconfigstatus-active}


### logging (string, optional) {#fluentdconfigstatus-logging}


### problems ([]string, optional) {#fluentdconfigstatus-problems}


### problemsCount (int, optional) {#fluentdconfigstatus-problemscount}



## FluentdConfigList

###  (metav1.TypeMeta, required) {#fluentdconfiglist-}


### metadata (metav1.ListMeta, optional) {#fluentdconfiglist-metadata}


### items ([]FluentdConfig, required) {#fluentdconfiglist-items}



