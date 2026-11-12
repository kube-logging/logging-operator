//go:build integration

// Copyright © 2026 Kube logging authors
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

package fluentbit

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

func TestMultilineExampleRuntime(t *testing.T) {
	var values struct {
		Logging struct {
			Fluentbit v1beta1.FluentbitSpec `json:"fluentbit"`
		} `json:"logging"`
	}
	example, err := os.ReadFile("../../../charts/logging-operator/values-multiline-example.yaml")
	require.NoError(t, err)
	require.NoError(t, yaml.Unmarshal(example, &values))
	spec := values.Logging.Fluentbit

	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "parsers.conf"), []byte(spec.CustomParsers), 0o600))
	config, err := generateConfig(fluentBitConfig{
		Flush: 1, Grace: 1, LogLevel: "warn", CoroStackSize: 24576,
		DefaultParsers: "/fluent-bit/etc/parsers.conf", CustomParsers: "/test/parsers.conf",
		DisableKubernetesFilter: true, FilterMultiline: spec.FilterMultiline,
		Input: fluentbitInputConfig{
			Values: map[string]string{
				"Path": "/test/*.log", "Tag": "kubernetes.*", "Path_Key": "test_path",
				"Read_From_Head": "On",
			},
			MultilineParser: spec.InputTail.MultilineParser,
		},
	})
	require.NoError(t, err)
	config += "\n[OUTPUT]\n    Name stdout\n    Match *\n    Format json_lines\n"
	require.NoError(t, os.WriteFile(filepath.Join(dir, "fluent-bit.conf"), []byte(config), 0o600))

	exceptions := map[string][]string{
		"java":   {"java.lang.RuntimeException: failed", "    at example.Main.run(Main.java:12)"},
		"python": {"Traceback (most recent call last):", "  File \"app.py\", line 1, in <module>", "ValueError: failed"},
	}
	expected := make(map[string][]string)
	for language, lines := range exceptions {
		for _, stream := range []string{"stdout", "stderr"} {
			for _, fragmented := range []bool{false, true} {
				name := fmt.Sprintf("%s-%s-%t.log", language, stream, fragmented)
				content := append([]string{"before"}, lines...)
				content = append(content, "after")
				var input strings.Builder
				for _, line := range content {
					prefix := "2026-01-01T00:00:00.000000000Z " + stream + " "
					if fragmented && len(line) > 1 {
						split := len(line) / 2
						fmt.Fprintf(&input, "%sP %s\n%sF %s\n", prefix, line[:split], prefix, line[split:])
					} else {
						fmt.Fprintf(&input, "%sF %s\n", prefix, line)
					}
				}
				require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(input.String()), 0o600))
				expected[name] = []string{"before", strings.Join(lines, "\n"), "after"}
			}
		}
	}

	run := exec.Command("docker", "run", "-d", "--rm", "-v", dir+":/test",
		"ghcr.io/fluent/fluent-bit:5.1.2", "-c", "/test/fluent-bit.conf")
	id, err := run.CombinedOutput()
	require.NoError(t, err, string(id))
	container := strings.TrimSpace(string(id))
	t.Cleanup(func() {
		output, err := exec.Command("docker", "rm", "-f", container).CombinedOutput()
		if err != nil {
			t.Errorf("remove Fluent Bit container: %v: %s", err, output)
		}
	})

	var output []byte
	ready := func() bool {
		output, err = exec.Command("docker", "logs", container).CombinedOutput()
		if err != nil {
			return false
		}
		return strings.Count(string(output), `"test_path"`) >= len(expected)*3
	}
	deadline := time.Now().Add(20 * time.Second)
	for !ready() && time.Now().Before(deadline) {
		time.Sleep(500 * time.Millisecond)
	}
	require.NoError(t, err, "reading Fluent Bit output: %s", output)
	require.GreaterOrEqual(t, strings.Count(string(output), `"test_path"`), len(expected)*3, "Fluent Bit output:\n%s", output)

	actual := make(map[string][]string)
	for _, line := range strings.Split(string(output), "\n") {
		if !strings.HasPrefix(line, "{") {
			continue
		}
		var record struct {
			Path   string `json:"test_path"`
			Log    string `json:"log"`
			Stream string `json:"stream"`
		}
		require.NoError(t, json.Unmarshal([]byte(line), &record))
		name := filepath.Base(record.Path)
		require.Contains(t, expected, name)
		require.Contains(t, name, "-"+record.Stream+"-")
		actual[name] = append(actual[name], record.Log)
	}
	require.Equal(t, expected, actual)
}
