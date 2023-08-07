---
title: Tag Normaliser
weight: 200
generated_file: true
---

# Fluentd Plugin to re-tag based on log metadata
## Overview
 More info at https://github.com/kube-logging/fluent-plugin-tag-normaliser

 # Available kubernetes metadata

 | Parameter | Description | Example |
 |-----------|-------------|---------|
 | ${pod_name} | Pod name | understood-butterfly-logging-demo-7dcdcfdcd7-h7p9n |
 | ${container_name} | Container name inside the Pod | logging-demo |
 | ${namespace_name} | Namespace name | default |
 | ${pod_id} | Kubernetes UUID for Pod | 1f50d309-45a6-11e9-b795-025000000001  |
 | ${labels} | Kubernetes Pod labels. This is a nested map. You can access nested attributes via `.`  | {"app":"logging-demo", "pod-template-hash":"7dcdcfdcd7" }  |
 | ${host} | Node hostname the Pod runs on | docker-desktop |
 | ${docker_id} | Docker UUID of the container | 3a38148aa37aa3... |

## Configuration
## Tag Normaliser parameters

### format (string, optional) {#tag normaliser parameters-format}

Re-Tag log messages info at [github](https://github.com/kube-logging/fluent-plugin-tag-normaliser) 

Default: ${namespace_name}.${pod_name}.${container_name}

### match_tag (string, optional) {#tag normaliser parameters-match_tag}

Tag used in match directive.  

Default:  kubernetes.**


 ## Example `Parser` filter configurations
 ```yaml
 apiVersion: logging.banzaicloud.io/v1beta1
 kind: Flow
 metadata:

	name: demo-flow

 spec:

	filters:
	  - tag_normaliser:
	      format: cluster1.${namespace_name}.${pod_name}.${labels.app}
	selectors: {}
	localOutputRefs:
	  - demo-output

 ```

 #### Fluentd Config Result
 ```yaml
 <match kubernetes.**>

	@type tag_normaliser
	@id test_tag_normaliser
	format cluster1.${namespace_name}.${pod_name}.${labels.app}

 </match>
 ```

---
