# logging-operator

Logging operator for Kubernetes based on Fluentd and Fluent-bit. For more details please follow up with this [post](https://banzaicloud.com/blog/k8s-logging-operator/).

## What is this operator for?

This operator helps you to pack together logging information with your applications. With the help of Custom Resource Definition you can describe the behaviour of your application within its charts. The operator does the rest.


## Project status: Alpha

This is the first version of this operator to showcase our plans of handling logs on Kubernetes. This version includes only basic configuration that will expand quickly. Stay tuned.


## Installing the operator

```
helm repo add  banzai-stable http://kubernetes-charts.banzaicloud.com/branch/master
helm install banzai-stable/logging-operator
```

## Example

The following steps set up an example configuration for sending nginx logs to S3.

### Create Secret

Create a manifest file for the AWS access key:

```
apiVersion: v1
kind: Secret
metadata:
  name: loggings3
type: Opaque
data:
  awsAccessKeyId: <base64encoded>
  awsSecretAccesKey: <base64encoded>
```

Submit the secret with kubectl:

```
kubectl apply -f secret.yaml
```

### Create LoggingOperator resource

Create a manifest that defines that you want to parse the nginx logs with the specified regular expressions on the standard output of pods with the `app: nginx` label, and store them in the given S3 bucket.

```
apiVersion: "logging.banzaicloud.com/v1alpha1"
kind: "LoggingOperator"
metadata:
  name: "nginx-logging"
spec:
  input:
    label:
      app: nginx
  filter:
    - type: parser
      name: parser-nginx
      parameters:
        - name: format
          value: '/^(?<remote>[^ ]*) (?<host>[^ ]*) (?<user>[^ ]*) \[(?<time>[^\]]*)\] "(?<method>\S+)(?: +(?<path>[^\"]*) +\S*)?" (?<code>[^ ]*) (?<size>[^ ]*)(?: "(?<referer>[^\"]*)" "(?<agent>[^\"]*)")?$/'
        - name: timeFormat
          value: "%d/%b/%Y:%H:%M:%S %z"
  output:
    - type: s3
      name: outputS3
      parameters:
        - name: aws_key_id
          valueFrom:
            secretKeyRef:
              name: loggings3
              key: awsAccessKeyId
        - name: aws_sec_key
          valueFrom:
            secretKeyRef:
              name: loggings3
              key: awsSecretAccesKey
        - name: s3_bucket
          value: logging-bucket
        - name: s3_region
          value: ap-northeast-1
```

## All in one example

If you just want to try the logging operator, install the operator and use our `nginx` example:

```
helm install ./deploy/helm/nginx-test
```

## Contributing

If you find this project useful here's how you can help:

- Send a pull request with your new features and bug fixes
- Help new users with issues they may encounter
- Support the development of this project and star this repo!
