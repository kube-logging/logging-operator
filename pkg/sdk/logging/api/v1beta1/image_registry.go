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

package v1beta1

import "strings"

var defaultImageRegistry string

// SetDefaultImageRegistry sets the registry used for the default images of the components the
// operator deploys, for example to pull them from a private mirror. It only applies where a
// resource doesn't set the image repository itself. An empty registry keeps the default images.
func SetDefaultImageRegistry(registry string) {
	defaultImageRegistry = strings.TrimSuffix(registry, "/")
}

// ReplaceImageRegistry replaces the registry of the image with the given one. An image without
// a registry, like fluent/fluent-bit, gets the registry added in front. An empty registry returns
// the image unchanged.
func ReplaceImageRegistry(registry, image string) string {
	registry = strings.TrimSuffix(registry, "/")
	if registry == "" {
		return image
	}
	if host, rest, found := strings.Cut(image, "/"); found && (strings.ContainsAny(host, ".:") || host == "localhost") {
		image = rest
	}
	return registry + "/" + image
}

func defaultImageRepository(repository string) string {
	return ReplaceImageRegistry(defaultImageRegistry, repository)
}
