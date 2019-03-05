# Developer's Guide

## Setting up local development environment


### Prerequisites

These steps are required to build the logging-operator and run on your computer.

### Install operator-sdk

Please follow the official guide for the **operator-sdk**: 
https://github.com/operator-framework/operator-sdk#quick-start

### Set-up the `kubernetes` context

Set up the kubernetes environment where you want create resources

#### Docker-for-mac

```
kubectl config use-context docker-for-desktop
```

#### Minikube

```
kubectl config use-context minikiube
```

### Install using operator-sdk local

```
operator-sdk up local
```

## Building docker image from the operator

```
$ docker build -t banzaicloud/logging-operator:local
```

### Using Helm to install logging-operator (with custom image)

Add banzaicloud-stable repo (or download the chart)

```
helm repo add banzaicloud-stable http://kubernetes-charts.banzaicloud.com/branch/master
helm repo update
```

Install the Helm deployment with custom (local) image

```
helm install banzaicloud-stable/logging-operator --set image.tag="local"
```

Verify installation

```
helm list
```

### Contribution

1. When contributing please check the issues and pull-requests weather your problem has been already addressed.
2. Open an issue and/or pull request describing your contribution
3. Please follow the issue and pull-request templates instructions
