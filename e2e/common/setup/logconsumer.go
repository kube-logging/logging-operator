// Copyright © 2021 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package setup

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/banzaicloud/logging-operator/e2e/common/cond"
)

func LogConsumer(t *testing.T, c client.Client, opts ...LogConsumerOption) LogConsumerResult {
	options := LogConsumerOptions{
		Name:      "log-consumer",
		Namespace: "default",
	}
	for _, opt := range opts {
		opt.ApplyToLogConsumerOptions(&options)
	}

	const logConsumerSource = `
	var log = console.log;
	function noop(){}
	var outputSocket = null;
	var outputSrv = net.createServer((socket) => {
		if (outputSocket) outputSocket.destroy();
		outputSocket = socket;
		socket.on('close', () => { outputSocket = null; });
	}).listen(8081, '0.0.0.0', () => log('output server is listening on 8081'));
	var inputSrv = http.createServer((req, res) => {
		var data = '';
		req.on('data', chunk => { data += chunk; });
		req.on('end', () => {
			log('got request');
			function ack() {
				res.writeHead(200);
				res.end();
			}
			if (outputSocket) outputSocket.write(data, 'utf8', ack);
			else ack();
		});
	});
	function on(cb) {
		cb ??= noop;
		if (inputSrv.listening) cb();
		else inputSrv.listen(8080, '0.0.0.0', () => {
			log('input server is listening on 8080');
			cb();
		});
	}
	function off(cb) {
		cb ??= noop;
		if (inputSrv.listening) inputSrv.close(() => {
			log('input server stopped listening');
			cb();
		});
		else cb();
	}
	http.createServer((req, res) => {
		let handler = {
			'/on': on,
			'/off': off,
		}[req.url];
		if (handler) handler(() => {
			res.writeHead(200);
			res.end();
		});
		else {
			log('invalid path', req.url);
			res.writeHead(404);
			res.end();
		}
	}).listen(8082, '0.0.0.0', () => log('control server is listening on 8082'));
	on();
	`
	logConsumerPod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      options.Name,
			Namespace: options.Namespace,
			Labels: map[string]string{
				"app.kubernetes.io/name": options.Name,
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "server",
					Image: "node:latest",
					Command: []string{
						"node",
					},
					Args: []string{
						"-e",
						logConsumerSource,
					},
				},
			},
		},
	}
	logConsumerSvc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      logConsumerPod.Name,
			Namespace: logConsumerPod.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: logConsumerPod.Labels,
			Ports: []corev1.ServicePort{
				{
					Name:     "input",
					Protocol: corev1.ProtocolTCP,
					Port:     8080,
				},
				{
					Name:     "output",
					Protocol: corev1.ProtocolTCP,
					Port:     8081,
				},
				{
					Name:     "control",
					Protocol: corev1.ProtocolTCP,
					Port:     8082,
				},
			},
		},
	}
	ctx := context.Background()
	require.NoError(t, c.Create(ctx, &logConsumerPod))
	require.NoError(t, c.Create(ctx, &logConsumerSvc))
	require.Eventually(t, cond.PodShouldBeRunning(t, c, client.ObjectKeyFromObject(&logConsumerPod)), 5*time.Minute, 5*time.Second)
	return LogConsumerResult{
		PodKey:     client.ObjectKeyFromObject(&logConsumerPod),
		ServiceKey: client.ObjectKeyFromObject(&logConsumerSvc),
	}
}

type LogConsumerOption interface {
	ApplyToLogConsumerOptions(*LogConsumerOptions)
}

type LogConsumerOptionFunc func(*LogConsumerOptions)

func (fn LogConsumerOptionFunc) ApplyToLogConsumerOptions(options *LogConsumerOptions) {
	fn(options)
}

type LogConsumerOptions struct {
	Name      string
	Namespace string
}

type LogConsumerResult struct {
	PodKey     client.ObjectKey
	ServiceKey client.ObjectKey
}

func (r LogConsumerResult) InputURL() string {
	return "http://" + r.ServiceKey.Name + "." + r.ServiceKey.Namespace + ".svc:8080"
}
