// Copyright 2020 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"net/http"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/shurcooL/vfsgen"
)

//go:generate go run main.go
func main() {
	crds := http.Dir(filepath.Join(getRepoRoot(), "config/crd/bases"))

	err := vfsgen.Generate(crds, vfsgen.Options{
		Filename:     filepath.Join(getRepoRoot(), "pkg/sdk/static/gen/crds/generated.go"),
		PackageName:  "crds",
		VariableName: "Root",
		FileModTime:  timePointer(time.Time{}),
	})
	if err != nil {
		panic(fmt.Sprintf("failed to generate crds vfs: %+v", err))
	}

	rbac := http.Dir(filepath.Join(getRepoRoot(), "config/rbac"))

	err = vfsgen.Generate(rbac, vfsgen.Options{
		Filename:     filepath.Join(getRepoRoot(), "pkg/sdk/static/gen/rbac/generated.go"),
		PackageName:  "rbac",
		VariableName: "Root",
		FileModTime:  timePointer(time.Time{}),
	})
	if err != nil {
		panic(fmt.Sprintf("failed to generate rbac vfs: %+v", err))
	}
}

// getRepoRoot returns the full path to the root of the repo
func getRepoRoot() string {
	_, filename, _, _ := runtime.Caller(0)

	dir := filepath.Dir(filename)

	return filepath.Dir(path.Join(dir, "../.."))
}

func timePointer(t time.Time) *time.Time {
	return &t
}
