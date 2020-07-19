---
title: Google Cloud Storage
weight: 200
generated_file: true
---

### GCSOutput
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| project | string | Yes | - | Project identifier for GCS<br> |
| keyfile | string | No | - | Path of GCS service account credentials JSON file<br> |
| credentials_json | *secret.Secret | No | - | GCS service account credentials in JSON format<br>[Secret](../secret/)<br> |
| client_retries | int | No | - | Number of times to retry requests on server error<br> |
| client_timeout | int | No | - | Default timeout to use in requests<br> |
| bucket | string | Yes | - | Name of a GCS bucket<br> |
| object_key_format | string | No |  %{path}%{time_slice}_%{index}.%{file_extension} | Format of GCS object keys <br> |
| path | string | No | - | Path prefix of the files on GCS<br> |
| store_as | string | No |  gzip | Archive format on GCS: gzip json text <br> |
| transcoding | bool | No | - | Enable the decompressive form of transcoding<br> |
| auto_create_bucket | bool | No |  true | Create GCS bucket if it does not exists <br> |
| hex_random_length | int | No |  4 | Max length of `%{hex_random}` placeholder(4-16) <br> |
| overwrite | bool | No |  false | Overwrite already existing path <br> |
| acl | string | No | - | Permission for the object in GCS: auth_read owner_full owner_read private project_private public_read<br> |
| storage_class | string | No | - | Storage class of the file: dra nearline coldline multi_regional regional standard<br> |
| encryption_key | string | No | - | Customer-supplied, AES-256 encryption key<br> |
| object_metadata | []ObjectMetadata | No | - | User provided web-safe keys and arbitrary string values that will returned with requests for the file as "x-goog-meta-" response headers.<br>[Object Metadata](#objectmetadata)<br> |
| format | *Format | No | - | [Format](../format/)<br> |
| buffer | *Buffer | No | - | [Buffer](../buffer/)<br> |
### ObjectMetadata
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| key | string | Yes | - | Key<br> |
| value | string | Yes | - | Value<br> |
