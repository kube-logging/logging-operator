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
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// ErrTimeout reports that an invocation was stopped by its deadline rather than
// by kind itself.
var ErrTimeout = errors.New("timed out")

const (
	pathEnv           = "KIND_PATH"
	imageEnv          = "KIND_IMAGE"
	commandTimeoutEnv = "KIND_COMMAND_TIMEOUT"

	defaultPath = "../../bin/kind"

	// timeoutFraction is the share of the enclosing -timeout that one invocation
	// may consume. The cap is derived from the surrounding budget rather than
	// from healthy timings because a stall only matters on a degraded runner,
	// which is exactly when every step legitimately takes several times longer.
	timeoutFraction = 0.8

	// fallbackTimeout applies when the binary has no deadline to derive from.
	fallbackTimeout = 15 * time.Minute

	defaultCleanupTimeout = 2 * time.Minute
	defaultWaitDelay      = 10 * time.Second
)

// Kind runs kind subcommands against one installation of the binary.
type Kind struct {
	Path  string
	Image string

	// CommandTimeout bounds a single invocation. Zero derives it from the test
	// binary's own -timeout on first use.
	CommandTimeout time.Duration

	// CleanupTimeout bounds the delete issued after a timeout. It is short
	// because the runner is already struggling by then and the caller still has
	// to report the failure within what is left of the budget.
	CleanupTimeout time.Duration

	// waitDelay caps how long Wait blocks on inherited pipes after the process
	// has been killed.
	waitDelay   time.Duration
	resolveOnce sync.Once
}

// New returns a Kind configured from KIND_PATH, KIND_IMAGE and
// KIND_COMMAND_TIMEOUT.
func New() *Kind {
	path := os.Getenv(pathEnv)
	if path == "" {
		path = defaultPath
		fmt.Printf("%s is not set, defaulting to %s\n", pathEnv, defaultPath)
	}

	return &Kind{
		Path:  path,
		Image: os.Getenv(imageEnv),
		// Left to be derived: the enclosing -timeout only becomes readable once
		// the testing package has parsed flags, which is after initialisation.
		CommandTimeout: resolveCommandTimeout(os.Getenv(commandTimeoutEnv), 0),
		CleanupTimeout: defaultCleanupTimeout,
		waitDelay:      defaultWaitDelay,
	}
}

func (k *Kind) CreateCluster(options CreateClusterOptions) error {
	if k.Image != "" && options.Image == "" {
		options.Image = k.Image
	}
	args := options.AppendToArgs([]string{"create", "cluster"})

	cmderr := &bytes.Buffer{}
	err := k.run(k.commandTimeout(), args, func(cmd *exec.Cmd) {
		cmd.Stdout = os.Stdout
		cmd.Stderr = io.MultiWriter(os.Stderr, cmderr)
	})

	switch {
	case err == nil:
		return nil

	case errors.Is(err, ErrTimeout):
		// kind removes a half-built cluster when one of its own actions fails,
		// but not when we kill it, so the leftovers have to go explicitly.
		cleanupErr := k.deleteCluster(k.CleanupTimeout, DeleteClusterOptions{
			Name:       options.Name,
			Kubeconfig: options.Kubeconfig,
		})
		if cleanupErr != nil {
			return fmt.Errorf("%w; deleting the partial cluster also failed: %w", err, cleanupErr)
		}
		return err

	case cmderr.Len() == 0:
		// An empty stderr means kind never ran, and the buffer on its own would
		// produce an error with no message.
		return err

	default:
		// common.isClusterAlreadyExistsError matches against kind's own wording.
		return errors.New(cmderr.String())
	}
}

func (k *Kind) DeleteCluster(options DeleteClusterOptions) error {
	return k.deleteCluster(k.commandTimeout(), options)
}

func (k *Kind) deleteCluster(timeout time.Duration, options DeleteClusterOptions) error {
	cmderr := &bytes.Buffer{}
	err := k.run(timeout, options.AppendToArgs([]string{"delete", "cluster"}), func(cmd *exec.Cmd) {
		cmd.Stderr = io.MultiWriter(os.Stderr, cmderr)
	})
	if err != nil && cmderr.Len() > 0 {
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(cmderr.String()))
	}
	return err
}

func (k *Kind) GetKubeconfig(options GetKubeconfigOptions) ([]byte, error) {
	args := options.AppendToArgs([]string{"get", "kubeconfig"})

	stdout := &bytes.Buffer{}
	err := k.run(k.commandTimeout(), args, func(cmd *exec.Cmd) {
		cmd.Stdout = stdout
		cmd.Stderr = os.Stderr
	})
	return stdout.Bytes(), err
}

func (k *Kind) LoadDockerImage(images []string, options LoadDockerImageOptions) error {
	if len(images) == 0 {
		return nil
	}
	args := append(options.AppendToArgs([]string{"load", "docker-image"}), images...)

	output := &bytes.Buffer{}
	err := k.run(k.commandTimeout(), args, func(cmd *exec.Cmd) {
		cmd.Stdout = output
		cmd.Stderr = output
	})
	if err != nil {
		return fmt.Errorf("kind load failed: %w\nOutput: %s", err, output.String())
	}
	return nil
}

// run executes a kind subcommand under the given deadline. configure wires up
// the command's output and may be nil.
func (k *Kind) run(timeout time.Duration, args []string, configure func(*exec.Cmd)) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, k.Path, args...)
	// kind shells out to docker and killing kind does not kill its children, so
	// without this Wait would block until every writer closed the pipe it
	// inherited.
	cmd.WaitDelay = k.waitDelay
	if configure != nil {
		configure(cmd)
	}

	err := cmd.Run()
	if err != nil && errors.Is(ctx.Err(), context.DeadlineExceeded) {
		// exec reports only "signal: killed", which does not say what stalled.
		return fmt.Errorf("kind %s %w after %s (override with %s)",
			strings.Join(args, " "), ErrTimeout, timeout, commandTimeoutEnv)
	}
	return err
}

func (k *Kind) commandTimeout() time.Duration {
	k.resolveOnce.Do(func() {
		if k.CommandTimeout > 0 {
			return
		}
		k.CommandTimeout = deriveCommandTimeout(enclosingTestTimeout())
		fmt.Printf("kind command timeout: %s (override with %s)\n", k.CommandTimeout, commandTimeoutEnv)
	})
	return k.CommandTimeout
}

// deriveCommandTimeout returns a cap strictly below enclosing, so a stall is
// cut while there is still budget left to clean up and name it. It does not
// bound the package as a whole: several invocations can each stay under the cap
// and still add up past the enclosing deadline.
func deriveCommandTimeout(enclosing time.Duration) time.Duration {
	if enclosing <= 0 {
		return fallbackTimeout
	}
	return time.Duration(float64(enclosing) * timeoutFraction)
}

// enclosingTestTimeout reports the -timeout this binary was started with, or
// zero if there is none to read.
func enclosingTestTimeout() time.Duration {
	f := flag.Lookup("test.timeout")
	if f == nil {
		return 0
	}
	getter, ok := f.Value.(flag.Getter)
	if !ok {
		return 0
	}
	timeout, ok := getter.Get().(time.Duration)
	if !ok {
		return 0
	}
	return timeout
}

// resolveCommandTimeout keeps fallback when raw is empty, unparseable or not
// positive. Non-positive is rejected rather than read as "no timeout" because a
// zero deadline has already expired, so it would fail every invocation at once.
func resolveCommandTimeout(raw string, fallback time.Duration) time.Duration {
	if raw == "" {
		return fallback
	}

	timeout, err := time.ParseDuration(raw)
	switch {
	case err != nil:
		fmt.Printf("%s=%q is not a valid duration, %s\n", commandTimeoutEnv, raw, keeping(fallback))
		return fallback
	case timeout <= 0:
		fmt.Printf("%s=%q is not a positive duration, %s\n", commandTimeoutEnv, raw, keeping(fallback))
		return fallback
	default:
		return timeout
	}
}

// keeping describes what a rejected value falls back to. A zero fallback is not
// a timeout of 0s, it is the signal to derive one from -timeout later.
func keeping(fallback time.Duration) string {
	if fallback <= 0 {
		return "deriving it from -timeout instead"
	}
	return fmt.Sprintf("keeping %s", fallback)
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
