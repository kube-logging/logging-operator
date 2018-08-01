package v1alpha1

import (
	"encoding/json"
	"fmt"
	"k8s.io/apimachinery/pkg/util/yaml"
	"testing"
)

var rawData = []byte(`
apiVersion: "logging.banzaicloud.com/v1alpha1"
kind: "LoggingOperator"
metadata:
  name: "nginx-logging"
spec:
  input:
    label:
      app-label: nginx
  filter:
    - type: parse
      format: '/^(?<remote>[^ ]*) (?<host>[^ ]*) (?<user>[^ ]*) \[(?<time>[^\]]*)\] "(?<method>\S+)(?: +(?<path>[^\"]*) +\S*)?" (?<code>[^ ]*) (?<size>[^ ]*)(?: "(?<referer>[^\"]*)" "(?<agent>[^\"]*)")?$/'
      timeFormat: "%d/%b/%Y:%H:%M:%S %z"
  output:
    - s3:
        name: adasd
        params:
          - name: aws_key_id
            valueFrom:
              secretKeyRef:
                name: loggingS3
                key: AWS_ACCESS_KEY_ID
          - name: aws_sec_key
            valueFrom:
              secretKeyRef:
                name: loggingS3
                key: AWS_SECRET_ACCESS_KEY
          - name: s3_bucket
            value: logging-bucket
          - name: s3_region
            value: ap-northeast-1
`)

func TestExampleCRD(t *testing.T) {
	output := LoggingOperator{}
	jsonData, err := yaml.ToJSON(rawData)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(jsonData)
	err = json.Unmarshal(jsonData, &output)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%#v", output)
}
