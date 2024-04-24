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
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

var KindPath string
var KindImage string

func init() {
	KindPath = os.Getenv("KIND_PATH")
	if KindPath == "" {
		KindPath = "../../bin/kind"
		fmt.Println("KIND_PATH is not set, defaulting to ../../bin/kind")
	}
	KindImage = os.Getenv("KIND_IMAGE")
}

func CreateCluster(options CreateClusterOptions) error {
	args := []string{"create", "cluster"}
	if KindImage != "" && options.Image == "" {
		options.Image = KindImage
	}
	args = options.AppendToArgs(args)
	cmd := exec.Command(KindPath, args...)
	cmderr := &bytes.Buffer{}
	cmd.Stdout = os.Stdout
	cmd.Stderr = io.MultiWriter(os.Stderr, cmderr)
	if err := cmd.Run(); err != nil {
		return errors.New(cmderr.String())
	}
	return nil
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
	args := []string{"delete", "cluster"}
	args = options.AppendToArgs(args)
	return exec.Command(KindPath, args...).Run()
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
	return exec.Command(KindPath, args...).Output()
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
	_, err := exec.Command(KindPath, args...).Output()
	return err
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
