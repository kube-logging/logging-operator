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
      name: parse-nginx
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
