// Copyright © 2021 Cisco Systems, Inc. and/or its affiliates
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

package kind

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	KindPath  string
	KindImage string
)

// ErrTimeout is returned when a kind invocation is stopped by its deadline
// rather than by kind itself.
var ErrTimeout = errors.New("timed out")

// CommandTimeout bounds a single kind invocation.
//
// kind's own --wait only covers the last of its create actions, waiting for
// control plane readiness. Provisioning the node containers, pulling the node
// image and `kind load docker-image` are bounded by nothing. When several
// clusters are built at once on a busy CI runner those phases can stall for
// longer than the whole test budget, and with no deadline here the only thing
// that ends the stall is `go test` panicking on its own -timeout. That kills
// every other test in the package and skips the deferred cluster cleanup, so
// one stalled cluster takes the rest of the package down with it.
//
// The value is deliberately far above anything a healthy run needs, so that
// it only fires on a genuine stall. Override it with KIND_COMMAND_TIMEOUT.
var CommandTimeout = 10 * time.Minute

// CleanupTimeout bounds the best-effort delete issued after a command times
// out. It is much shorter than CommandTimeout: the runner is already
// struggling by then, and the caller still has to report the failure inside
// the remaining test budget.
var CleanupTimeout = 2 * time.Minute

const commandTimeoutEnv = "KIND_COMMAND_TIMEOUT"

// waitDelay is how long Wait may keep waiting on inherited output pipes after
// the process itself has been killed. A variable so the tests can shorten it.
var waitDelay = 10 * time.Second

func init() {
	KindPath = os.Getenv("KIND_PATH")
	if KindPath == "" {
		KindPath = "../../bin/kind"
		fmt.Println("KIND_PATH is not set, defaulting to ../../bin/kind")
	}
	KindImage = os.Getenv("KIND_IMAGE")
	CommandTimeout = resolveCommandTimeout(os.Getenv(commandTimeoutEnv), CommandTimeout)
}

// resolveCommandTimeout reads the timeout from raw, keeping fallback when raw
// is empty or unparseable. A bad value is not worth failing the run over, but
// it is worth saying out loud.
func resolveCommandTimeout(raw string, fallback time.Duration) time.Duration {
	if raw == "" {
		return fallback
	}
	timeout, err := time.ParseDuration(raw)
	if err != nil {
		fmt.Printf("%s=%q is not a valid duration, keeping %s\n", commandTimeoutEnv, raw, fallback)
		return fallback
	}
	return timeout
}

// runKind runs a kind subcommand under the given deadline. configure wires up
// the command's output and may be nil.
//
// When the deadline is what stopped the process, exec only reports
// "signal: killed", which says nothing about what was stuck. So that case is
// reported as an error wrapping ErrTimeout and naming the subcommand.
func runKind(timeout time.Duration, args []string, configure func(*exec.Cmd)) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, KindPath, args...)
	// kind shells out to docker, and killing kind does not kill its children.
	// Whenever output is wired to something other than a file, exec pipes it
	// and Wait blocks until every writer closes that pipe, so a surviving
	// docker would keep us here well past the deadline. WaitDelay caps that.
	cmd.WaitDelay = waitDelay
	if configure != nil {
		configure(cmd)
	}

	err := cmd.Run()
	if err != nil && errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return fmt.Errorf("kind %s %w after %s (override with %s)", strings.Join(args, " "), ErrTimeout, timeout, commandTimeoutEnv)
	}
	return err
}

func CreateCluster(options CreateClusterOptions) error {
	args := []string{"create", "cluster"}
	if KindImage != "" && options.Image == "" {
		options.Image = KindImage
	}
	args = options.AppendToArgs(args)

	cmderr := &bytes.Buffer{}
	err := runKind(CommandTimeout, args, func(cmd *exec.Cmd) {
		cmd.Stdout = os.Stdout
		cmd.Stderr = io.MultiWriter(os.Stderr, cmderr)
	})

	switch {
	case err == nil:
		return nil
	case errors.Is(err, ErrTimeout):
		// kind removes a half-built cluster when one of its own actions fails,
		// but not when we kill it, so the leftovers have to go explicitly.
		// Otherwise the node containers hold their share of the runner for the
		// rest of the job and make the next cluster likelier to stall too.
		if deleteErr := deleteCluster(CleanupTimeout, DeleteClusterOptions{Name: options.Name}); deleteErr != nil {
			return fmt.Errorf("%w; deleting the partial cluster also failed: %w", err, deleteErr)
		}
		return err
	default:
		return errors.New(cmderr.String())
	}
}

type CreateClusterOptions struct {
	GlobalOptions
	Config     string
	Image      string
	Kubeconfig string
	Name       string
	Retain     bool
	Wait       string
}

func (options CreateClusterOptions) AppendToArgs(args []string) []string {
	args = options.GlobalOptions.AppendToArgs(args)
	if options.Config != "" {
		args = append(args, "--config", options.Config)
	}
	if options.Image != "" {
		args = append(args, "--image", options.Image)
	}
	if options.Kubeconfig != "" {
		args = append(args, "--kubeconfig", options.Kubeconfig)
	}
	if options.Name != "" {
		args = append(args, "--name", options.Name)
	}
	if options.Retain {
		args = append(args, "--retain")
	}
	if options.Wait != "" {
		args = append(args, "--wait", options.Wait)
	}
	return args
}

func DeleteCluster(options DeleteClusterOptions) error {
	return deleteCluster(CommandTimeout, options)
}

func deleteCluster(timeout time.Duration, options DeleteClusterOptions) error {
	args := []string{"delete", "cluster"}
	args = options.AppendToArgs(args)
	return runKind(timeout, args, nil)
}

type DeleteClusterOptions struct {
	GlobalOptions
	Kubeconfig string
	Name       string
}

func (options DeleteClusterOptions) AppendToArgs(args []string) []string {
	args = options.GlobalOptions.AppendToArgs(args)
	if options.Kubeconfig != "" {
		args = append(args, "--kubeconfig", options.Kubeconfig)
	}
	if options.Name != "" {
		args = append(args, "--name", options.Name)
	}
	return args
}

func GetKubeconfig(options GetKubeconfigOptions) ([]byte, error) {
	args := []string{"get", "kubeconfig"}
	args = options.AppendToArgs(args)

	stdout := &bytes.Buffer{}
	err := runKind(CommandTimeout, args, func(cmd *exec.Cmd) {
		cmd.Stdout = stdout
		cmd.Stderr = os.Stderr
	})
	return stdout.Bytes(), err
}

type GetKubeconfigOptions struct {
	GlobalOptions
	Internal bool
	Name     string
}

func (options GetKubeconfigOptions) AppendToArgs(args []string) []string {
	args = options.GlobalOptions.AppendToArgs(args)
	if options.Internal {
		args = append(args, "--internal")
	}
	if options.Name != "" {
		args = append(args, "--name", options.Name)
	}
	return args
}

func LoadDockerImage(images []string, options LoadDockerImageOptions) error {
	if len(images) == 0 {
		return nil
	}

	args := []string{"load", "docker-image"}
	args = options.AppendToArgs(args)
	args = append(args, images...)

	output := &bytes.Buffer{}
	err := runKind(CommandTimeout, args, func(cmd *exec.Cmd) {
		cmd.Stdout = output
		cmd.Stderr = output
	})
	if err != nil {
		return fmt.Errorf("kind load failed: %w\nOutput: %s", err, output.String())
	}

	return nil
}

type LoadDockerImageOptions struct {
	GlobalOptions
	Name  string
	Nodes []string
}

func (options LoadDockerImageOptions) AppendToArgs(args []string) []string {
	args = options.GlobalOptions.AppendToArgs(args)
	if options.Name != "" {
		args = append(args, "--name", options.Name)
	}
	if len(options.Nodes) > 0 {
		args = append(args, "--nodes", strings.Join(options.Nodes, ","))
	}
	return args
}

type GlobalOptions struct {
	LogLevel  string
	Quiet     bool
	Verbosity string
}

func (options GlobalOptions) AppendToArgs(args []string) []string {
	if options.LogLevel != "" {
		args = append(args, "--loglevel", options.LogLevel)
	}
	if options.Quiet {
		args = append(args, "--quiet")
	}
	if options.Verbosity != "" {
		args = append(args, "--verbosity", options.Verbosity)
	}
	return args
}
