
# TLS Configuration

To configure TLS for Fluentd and Fluentbit the operator needs TLS certificates
set via the Fluentd and Fluentbit Custom Resources respectively. This can be
done in two ways:

## Generic Opaque secret (default)

Create a secret like this:

```
apiVersion: v1
data:
  caCert: ...
  clientCert: ...
  clientKey: ...
  serverCert: ...
  serverKey: ...
kind: Secret
metadata:
  name: something-something-tls
type: Opaque
```

Note that we are providing three certificates in the same secret, one for
Fluentd (`serverCert`), one for Fluentbit (`clientCert`), and the CA
certificate (`caCert`).

Then in your custom resource configure like this:

```
apiVersion: logging.banzaicloud.com/v1alpha1
kind: Fluentd/Fluentbit
metadata:
  name: my-fluent-thing
spec:
  ...
  tls:
    enabled: true
    secretName: something-something-tls
    sharedKey: changeme
```


## `kubernetes.io/tls`

The alternative is if your certificates are in secrets of type `kubernetes.io/tls`, e.g.

```
apiVersion: v1
data:
  ca.crt: LS0tLS1...
  tls.crt: LS0tLS1...
  tls.key: LS0tLS1...
kind: Secret
metadata:
  name: something-something-tls
type: kubernetes.io/tls
```

Then configure your custom resources like this:

```
apiVersion: logging.banzaicloud.com/v1alpha1
kind: Fluentd/Fluentbit
metadata:
  name: my-fluent-thing
spec:
  ...
  tls:
    enabled: true
    secretName: something-something-tls
    secretType: tls
    sharedKey: changeme
```

Note: in this case we can use the same secret for both Fluentbit and Fluentd,
or create separate secrets for each.

Note: the secret's data include the CA certificate, which is in-line with the
structure created by [jetstack/cert-manager](https://github.com/jetstack/cert-manager/).

## Usage with the helm chart

For the generic Opaque secret just set `tls.enabled=True` and optionally provide the `tls.secretName` value to use your own certificates (instead of the automatically generated ones from the chart).

For `kubernetes.io/tls` install `logging-operator-fluent` with a `values.yaml` like this:

```
tls:
  enabled: true

fluentbit:
  tlsSecret: something-something-tls

fluentd:
  tlsSecret: otherthing-otherthing-tls
```

For more information see the helm chart's [README.md](https://github.com/banzaicloud/logging-operator/blob/master/charts/logging-operator-fluent/README.md).
