## AxoSyslog

AxoSyslog is the Schema for the AxoSyslogs API

###  (metav1.TypeMeta, required) {#axosyslog-}


### metadata (metav1.ObjectMeta, optional) {#axosyslog-metadata}


### spec (AxoSyslogSpec, optional) {#axosyslog-spec}


### status (AxoSyslogStatus, optional) {#axosyslog-status}



## AxoSyslogSpec

AxoSyslogSpec defines the desired state of AxoSyslog

### controlNamespace (string, required) {#axosyslogspec-controlnamespace}

ControlNamespace is the namespace where the AxoSyslog StatefulSet and related components will be deployed 


### destinations ([]string, optional) {#axosyslogspec-destinations}

Destinations is a list of destinations to be rendered in the AxoSyslog configuration 


### logPaths ([]LogPath, optional) {#axosyslogspec-logpaths}

LogPaths is a list of log paths to be rendered in the AxoSyslog configuration 



## LogPath

LogPath defines a single log path that will be rendered in the AxoSyslog configuration

### destination (string, optional) {#logpath-destination}

name of a destionation to be used in the log path 


### filterx (string, optional) {#logpath-filterx}

filterx block to be rendered within the log path 



## Destination

Destination defines a single destination that will be rendered in the AxoSyslog configuration

### config (string, optional) {#destination-config}

Config is the configuration for the destination 


### name (string, optional) {#destination-name}

Name of the destination 



## AxoSyslogStatus

AxoSyslogStatus defines the observed state of AxoSyslog

### problems ([]string, optional) {#axosyslogstatus-problems}

Problems with the AxoSyslog resource 


### problemsCount (int, optional) {#axosyslogstatus-problemscount}

Count of problems with the AxoSyslog resource 


### sources ([]Source, optional) {#axosyslogstatus-sources}

Sources configured for AxoSyslog 



## Source

Source represents the source of logs for AxoSyslog

### otlp (*OTLPSource, optional) {#source-otlp}

OTLP specific configuration 



## OTLPSource

OTLPSource contains configuration for OpenTelemetry Protocol sources

### endpoint (string, optional) {#otlpsource-endpoint}

Endpoint for the OTLP source 



## AxoSyslogList

AxoSyslogList contains a list of AxoSyslog

###  (metav1.TypeMeta, required) {#axosysloglist-}


### metadata (metav1.ListMeta, required) {#axosysloglist-metadata}


### items ([]AxoSyslog, required) {#axosysloglist-items}



