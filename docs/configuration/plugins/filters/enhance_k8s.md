---
title: Enhance K8s Metadata
weight: 200
generated_file: true
---

# [Enhance K8s Metadata](https://github.com/SumoLogic/sumologic-kubernetes-fluentd/tree/main/fluent-plugin-enhance-k8s-metadata)
## Overview
 Fluentd Filter plugin to fetch several metadata for a Pod

## Configuration
## EnhanceK8s

### api_groups ([]string, optional) {#enhancek8s-api_groups}

Kubernetes resources api groups

Default: `["apps/v1", "extensions/v1beta1"]`

### bearer_token_file (string, optional) {#enhancek8s-bearer_token_file}

Bearer token path

Default: nil

### ca_file (secret.Secret, optional) {#enhancek8s-ca_file}

Kubernetes API CA file

Default: nil

### cache_refresh (int, optional) {#enhancek8s-cache_refresh}

Cache refresh

Default: 60*60

### cache_refresh_variation (int, optional) {#enhancek8s-cache_refresh_variation}

Cache refresh variation

Default: 60*15

### cache_size (int, optional) {#enhancek8s-cache_size}

Cache size

Default: 1000

### cache_ttl (int, optional) {#enhancek8s-cache_ttl}

Cache TTL

Default: 60*60*2

### client_cert (secret.Secret, optional) {#enhancek8s-client_cert}

Kubernetes API Client certificate

Default: nil

### client_key (secret.Secret, optional) {#enhancek8s-client_key}

// Kubernetes API Client certificate key

Default: nil

### core_api_versions ([]string, optional) {#enhancek8s-core_api_versions}

Kubernetes core API version (for different Kubernetes versions)

Default: ['v1']

### data_type (string, optional) {#enhancek8s-data_type}

Sumologic data type

Default: metrics

### in_namespace_path ([]string, optional) {#enhancek8s-in_namespace_path}

parameters for read/write record

Default: `['$.namespace']`

### in_pod_path ([]string, optional) {#enhancek8s-in_pod_path}

Default: `['$.pod','$.pod_name']`

### kubernetes_url (string, optional) {#enhancek8s-kubernetes_url}

Kubernetes API URL

Default: nil

### ssl_partial_chain (*bool, optional) {#enhancek8s-ssl_partial_chain}

If `ca_file` is for an intermediate CA, or otherwise we do not have the root CA and want to trust the intermediate CA certs we do have, set this to `true` - this corresponds to the openssl s_client -partial_chain flag and X509_V_FLAG_PARTIAL_CHAIN

Default: false

### secret_dir (string, optional) {#enhancek8s-secret_dir}

Service account directory

Default: /var/run/secrets/kubernetes.io/serviceaccount

### verify_ssl (*bool, optional) {#enhancek8s-verify_ssl}

Verify SSL

Default: true




## Example `EnhanceK8s` filter configurations

{{< highlight yaml >}}
apiVersion: logging.banzaicloud.io/v1beta1
kind: Logging
metadata:
  name: demo-flow
spec:
  globalFilters:
    - enhanceK8s: {}
{{</ highlight >}}

Fluentd config result:

{{< highlight xml >}}
<filter **>
  @type enhance_k8s_metadata
  @id test_enhanceK8s
</filter>
{{</ highlight >}}


---
