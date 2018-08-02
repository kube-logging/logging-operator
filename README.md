# logging-operator

Logging operator for Kubernetes based on Fluentd and Fluent-bit

## What is this operator for?

This operator helps you to pack together logging information with your applications. With the help of Custom Resource Definition you can describe the behaviour of your application within it's charts. The operator does the rest.


## Projact status: Alpha

This is the first version of this operator to showcase our plans of handling logs on Kubernetes. This version includes only basic configuration that will expand quickly. Stay tuned.


## Installing the operator

```
helm install banzai-stable/logging-operator
```

## Example

The following example shows what you need to provide for proper log shipping.

### Create Secret

Create a secret file with valid input

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

Create secret with kubectl

```
kubectl apply -f secret.yaml
```

### CustomResourceDefinitions

This is an example how-to define log definition for an application

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
    - type: parse
      format: |
        /^(?<remote>[^ ]*) (?<host>[^ ]*) (?<user>[^ ]*) \[(?<time>[^\]]*)\] "(?<method>\S+)(?: +(?<path>[^\"]*) +\S*)?" (?<code>[^ ]*) (?<size>[^ ]*)(?: "(?<referer>[^\"]*)" "(?<agent>[^\"]*)")?$/

      timeFormat: '%d/%b/%Y:%H:%M:%S %z'
  output:
    - s3:
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

## AllInOne example

If you just want to try out use our `nginx` example

```
helm install ./nginx-example
```
