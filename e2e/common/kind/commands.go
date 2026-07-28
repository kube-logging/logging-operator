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

var (
	KindPath  string
	KindImage string
)

// ErrTimeout is returned when a kind invocation is stopped by its deadline
// rather than by kind itself.
var ErrTimeout = errors.New("timed out")

// CommandTimeout bounds a single kind invocation. Zero means it is derived
// from the test binary's own -timeout on first use; see commandTimeout.
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
// Set it directly, or with KIND_COMMAND_TIMEOUT, to pin a value.
var CommandTimeout time.Duration

// timeoutFraction is how much of the enclosing -timeout a single kind call may
// consume before we call it stalled.
//
// The deadline has to be derived from the budget it sits inside, not guessed
// from how long a healthy run takes. Healthy timings are the wrong reference:
// this only fires when the runner is degraded, and that is exactly when every
// step legitimately takes several times longer than healthy. A cap chosen from
// healthy numbers ends up inside the range where real work still completes,
// and then it fails runs that would have passed.
//
// Anything still running this late has left too little of the budget for the
// rest of the test to finish, so cutting it can only replace an anonymous
// whole-binary panic with a named failure and a cleaned-up cluster. It cannot
// turn a passing run red.
const timeoutFraction = 0.8

// fallbackTimeout applies when the enclosing deadline cannot be read, which
// means the binary was started with -timeout 0 and has no deadline of its own.
const fallbackTimeout = 15 * time.Minute

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
	// Not resolved here: the enclosing -timeout is only readable once the
	// testing package has parsed flags, which happens after init.
	CommandTimeout = resolveCommandTimeout(os.Getenv(commandTimeoutEnv), 0)
}

// commandTimeout is the deadline for one kind call, worked out once.
func commandTimeout() time.Duration {
	resolveOnce.Do(func() {
		if CommandTimeout > 0 {
			return // pinned by KIND_COMMAND_TIMEOUT or by a caller
		}
		CommandTimeout = deriveCommandTimeout(enclosingTestTimeout())
		fmt.Printf("kind command timeout: %s (override with %s)\n", CommandTimeout, commandTimeoutEnv)
	})
	return CommandTimeout
}

var resolveOnce sync.Once

// deriveCommandTimeout turns the enclosing test deadline into a per-command
// one. It must always come out strictly below enclosing, so that cutting a
// command short can never be the reason a package misses its own deadline.
func deriveCommandTimeout(enclosing time.Duration) time.Duration {
	if enclosing <= 0 {
		return fallbackTimeout
	}
	return time.Duration(float64(enclosing) * timeoutFraction)
}

// enclosingTestTimeout reports the -timeout this test binary was started with,
// or zero if there is none to read.
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

// resolveCommandTimeout reads the timeout from raw, keeping fallback when raw
// is empty, unparseable, or not positive. A bad value is not worth failing the
// run over, but it is worth saying out loud.
//
// Non-positive is rejected rather than read as "no timeout". A zero deadline
// is already expired, so KIND_COMMAND_TIMEOUT=0 would fail every kind call
// instantly — the opposite of what anyone writing that would mean by it. Pass
// a large duration if you want an effectively unbounded run.
func resolveCommandTimeout(raw string, fallback time.Duration) time.Duration {
	if raw == "" {
		return fallback
	}
	timeout, err := time.ParseDuration(raw)
	if err != nil {
		fmt.Printf("%s=%q is not a valid duration, keeping %s\n", commandTimeoutEnv, raw, fallback)
		return fallback
	}
	if timeout <= 0 {
		fmt.Printf("%s=%q is not a positive duration, keeping %s\n", commandTimeoutEnv, raw, fallback)
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
	err := runKind(commandTimeout(), args, func(cmd *exec.Cmd) {
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
		// kind's own diagnosis is on stderr, and isClusterAlreadyExistsError
		// matches against it, so that is what callers need to see. But stderr
		// is empty when the failure was starting the process at all — a
		// KIND_PATH pointing at nothing, say — and an error with an empty
		// message tells nobody anything.
		if cmderr.Len() == 0 {
			return err
		}
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
	return deleteCluster(commandTimeout(), options)
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
	err := runKind(commandTimeout(), args, func(cmd *exec.Cmd) {
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
	err := runKind(commandTimeout(), args, func(cmd *exec.Cmd) {
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
