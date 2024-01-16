---
title: SyslogNGConfig
weight: 200
generated_file: true
---

## SyslogNGConfig

###  (metav1.TypeMeta, required) {#syslogngconfig-}


### metadata (metav1.ObjectMeta, optional) {#syslogngconfig-metadata}


### spec (SyslogNGSpec, optional) {#syslogngconfig-spec}


### status (SyslogNGConfigStatus, optional) {#syslogngconfig-status}



## SyslogNGConfigStatus

### active (*bool, optional) {#syslogngconfigstatus-active}


### logging (string, optional) {#syslogngconfigstatus-logging}


### problems ([]string, optional) {#syslogngconfigstatus-problems}


### problemsCount (int, optional) {#syslogngconfigstatus-problemscount}



## SyslogNGConfigList

###  (metav1.TypeMeta, required) {#syslogngconfiglist-}


### metadata (metav1.ListMeta, optional) {#syslogngconfiglist-metadata}


### items ([]SyslogNGConfig, required) {#syslogngconfiglist-items}



