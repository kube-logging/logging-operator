# Deploying logging-operator with terraform

This sample shows how to deploy logging-operator with terraform/OpenTofu.
The sample configures logging-operator to send logs to a S3 bucket and deploys the following:
- logging-operator
- syslog-ng flow
- syslog-ng output for S3
- fluentbit agent configured to parse CRI logs and store the parsed data under ```json.msg```

# Prerequisites

- terraform or opentofu
- kubernetes cluster
- AWS profile used for deploying
- AWS S3 bucket for storing the logs
- AWS IAM access_key with permissions to store data to s3 bucket

## Variables
```
variable "aws_profile_name" {
  description = "AWS credentials profile to use for provisioning"
}

variable "aws_region" {
  description = "AWS region"
}

variable "aws_s3_bucket" {
  description = "AWS S3 bucket used to store logs,"
}

variable "aws_iam_access_key_id" {
  description = "IAM access key for storing logs in S3."
}

variable "aws_iam_access_key_secret" {
  description = "IAM access key secret for storing logs in S3."
}

variable "eks_cluster_name" {
  description = "EKS cluster name."
  type        = string
}



```

## Main file

```
provider aws {
  profile = var.aws_profile_name
  region = var.aws_region
}

locals {
  namespace = "apps"
}

data "aws_eks_cluster" "cluster" {
  name = var.eks_cluster_name
}

data "aws_eks_cluster_auth" "cluster" {
  name = var.eks_cluster_name
}

data "aws_caller_identity" "current" {}

provider "kubernetes" {
  host                   = data.aws_eks_cluster.cluster.endpoint
  cluster_ca_certificate = base64decode(data.aws_eks_cluster.cluster.certificate_authority.0.data)
  token                  = data.aws_eks_cluster_auth.cluster.token
}

provider "helm" {
  kubernetes {
    host                   = data.aws_eks_cluster.cluster.endpoint
    cluster_ca_certificate = base64decode(data.aws_eks_cluster.cluster.certificate_authority.0.data)
    token                  = data.aws_eks_cluster_auth.cluster.token
  }
}

resource "kubernetes_secret" "access-logs" {
  metadata {
    name = "s3-logging-aws-key"
    namespace = local.namespace
  }

  data = {
    access-key = var.aws_iam_access_key_id
    secret-key = var.aws_iam_access_key_secret
  }
}

# https://github.com/kube-logging/logging-operator/tree/master/charts/logging-operator
resource "helm_release" "logging-operator" {
  repository = "oci://ghcr.io/kube-logging/helm-charts/"
  chart      = "logging-operator"
  name       = "logging-operator"
  namespace  = "logging"

  set {
    name = "createCustomResource"
    value = false
  }
}

resource "kubernetes_manifest" "logging" {
  manifest = {
    apiVersion = "logging.banzaicloud.io/v1beta1"
    kind       = "Logging"

    metadata = {
      name = "logging-operator"
    }

    spec = {
      controlNamespace = "logging"
      syslogNG = {
        globalOptions = {
          log_level = "info"
        }
      }
      watchNamespaces = [ local.namespace ]
    }
  }
}

# https://github.com/kube-logging/logging-operator/blob/master/config/crd/bases/logging.banzaicloud.io_syslogngoutputs.yaml
resource "kubernetes_manifest" "syslog-ng-output" {
  manifest = {
    apiVersion = "logging.banzaicloud.io/v1beta1"
    kind       = "SyslogNGOutput"

    metadata = {
      name = "syslog-s3-output"
      namespace = local.namespace
    }

    spec = {
      s3 = {
        url = "https://s3.${var.aws_region}.amazonaws.com"
        bucket = var.aws_s3_bucket
        region = var.aws_region
        access_key = {
          valueFrom = {
            secretKeyRef = {
              name: "s3-logging-aws-key"
              key: "access-key"
            }
          }
        }
        secret_key = {
          valueFrom = {
            secretKeyRef = {
              name: "s3-logging-aws-key"
              key: "secret-key"
            }
          }
        }
        # sample partitioning for AWS Athena
        object_key = "backend/$${R_YEAR}/$${R_MONTH}/logs"
      }
    }
  }
}

resource "kubernetes_manifest" "syslog-ng-flow" {
  manifest = {
    apiVersion = "logging.banzaicloud.io/v1beta1"
    kind       = "SyslogNGFlow"

    metadata = {
      name = "syslog-ng-apps-flow"
      namespace = local.namespace
    }

    spec = {
      match = {
        and = [
          {
            regexp = {
              value = "json.kubernetes.labels.app.kubernetes.io/name"
              pattern = "backend"
              type = "string"
            }
          },
          {
            regexp = {
              value = "json.msg.log_type"
              pattern = "nginx_access"
              type = "string"
            }
          }
        ]
      }

      localOutputRefs = [
        "syslog-s3-output"
      ]
    }
  }

  field_manager {
    force_conflicts = true
  }
}

resource "kubernetes_manifest" "fluentbit-agents" {
  manifest = {
    apiVersion = "logging.banzaicloud.io/v1beta1"
    kind       = "FluentbitAgent"

    metadata = {
      # Use the name of the logging resource
      name = "logging-operator"
    }

    spec = {
      # the indentation needs to be 4 spaces or fluentbit may ignore the directive
      customParsers = <<EOF
[PARSER]
    Name cri-parser
    Format regex
    Regex ^(?<time>[^ ]+) (?<stream>stdout|stderr) (?<logtag>[^ ]*) (?<log>.*)$
    Time_Key time
    Time_Format %Y-%m-%dT%H:%M:%S.%L%z
EOF

      filterKubernetes = {
        Merge_Log_Key: "msg"
      }

      inputTail = {
        Parser = "cri-parser"
      }

      filterModify = [
        {
          rules = [
            {
              Remove = {
                key = "log"
              }
            }
          ]
        }
      ]
    }
  }
}
```
