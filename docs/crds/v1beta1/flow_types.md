### FlowSpec
#### FlowSpec is the Kubernetes spec for Flows

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| selectors | map[string]string | Yes | - |  |
| filters | []Filter | No | - |  |
| loggingRef | string | No | - |  |
| outputRefs | []string | Yes | - |  |
### Filter
#### Filter definition for FlowSpec

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| stdout | *filter.StdOutFilterConfig | No | - |  |
| parser | *filter.ParserConfig | No | - |  |
| tag_normaliser | *filter.TagNormaliser | No | - |  |
| dedot | *filter.DedotFilterConfig | No | - |  |
| record_transformer | *filter.RecordTransformer | No | - |  |
| record_modifier | *filter.RecordModifier | No | - |  |
| geoip | *filter.GeoIP | No | - |  |
| concat | *filter.Concat | No | - |  |
| detectExceptions | *filter.DetectExceptions | No | - |  |
| grep | *filter.GrepConfig | No | - |  |
| prometheus | *filter.PrometheusConfig | No | - |  |
### FlowStatus
#### FlowStatus defines the observed state of Flow

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
### Flow
#### Flow Kubernetes object

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
|  | metav1.TypeMeta | Yes | - |  |
| metadata | metav1.ObjectMeta | No | - |  |
| spec | FlowSpec | No | - |  |
| status | FlowStatus | No | - |  |
### FlowList
#### FlowList contains a list of Flow

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
|  | metav1.TypeMeta | Yes | - |  |
| metadata | metav1.ListMeta | No | - |  |
| items | []Flow | Yes | - |  |
