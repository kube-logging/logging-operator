// Copyright 2019 The Kubernetes Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This code is based on https://raw.githubusercontent.com/kubernetes-sigs/cluster-api/ee96c31eea0f9b1b2f4cffae5c7da2c274722e76/util/restmapper/cached.go

package k8sutil

import (
	"sync"

	"golang.org/x/time/rate"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func NewCached(config *rest.Config) (meta.RESTMapper, error) {
	c := &Cached{
		limiter: rate.NewLimiter(rate.Limit(1), 2),
		factory: func() (meta.RESTMapper, error) {
			return apiutil.NewDiscoveryRESTMapper(config)
		},
	}
	if err := c.flush(); err != nil {
		return nil, err
	}
	return c, nil
}

type Cached struct {
	sync.Mutex

	limiter *rate.Limiter
	factory func() (meta.RESTMapper, error)
	mapper  meta.RESTMapper
}

func (c *Cached) flush() error {
	c.Lock()
	defer c.Unlock()

	var err error
	if c.mapper == nil || c.limiter.Allow() {
		c.mapper, err = c.factory()
	}
	return err
}

func (c *Cached) shouldFlushOn(err error) bool {
	switch err.(type) { //nolint:errorlint
	case *meta.NoKindMatchError:
		return true
	}
	return false
}

func (c *Cached) onError(err error) bool {
	if !c.shouldFlushOn(err) {
		return false
	}
	if err := c.flush(); err != nil {
		log.Log.Error(err, "failed to reload RESTMapper")
		return false
	}
	return true
}

func (c *Cached) KindFor(resource schema.GroupVersionResource) (schema.GroupVersionKind, error) {
	gvk, err := c.mapper.KindFor(resource)
	if c.onError(err) {
		gvk, err = c.mapper.KindFor(resource)
	}
	return gvk, err
}

func (c *Cached) KindsFor(resource schema.GroupVersionResource) ([]schema.GroupVersionKind, error) {
	gvks, err := c.mapper.KindsFor(resource)
	if c.onError(err) {
		gvks, err = c.mapper.KindsFor(resource)
	}
	return gvks, err
}

func (c *Cached) ResourceFor(input schema.GroupVersionResource) (schema.GroupVersionResource, error) {
	gvr, err := c.mapper.ResourceFor(input)
	if c.onError(err) {
		gvr, err = c.mapper.ResourceFor(input)
	}
	return gvr, err
}

func (c *Cached) ResourcesFor(input schema.GroupVersionResource) ([]schema.GroupVersionResource, error) {
	gvrs, err := c.mapper.ResourcesFor(input)
	if c.onError(err) {
		gvrs, err = c.mapper.ResourcesFor(input)
	}
	return gvrs, err
}

func (c *Cached) RESTMapping(gk schema.GroupKind, versions ...string) (*meta.RESTMapping, error) {
	m, err := c.mapper.RESTMapping(gk, versions...)
	if c.onError(err) {
		m, err = c.mapper.RESTMapping(gk, versions...)
	}
	return m, err
}

func (c *Cached) RESTMappings(gk schema.GroupKind, versions ...string) ([]*meta.RESTMapping, error) {
	ms, err := c.mapper.RESTMappings(gk, versions...)
	if c.onError(err) {
		ms, err = c.mapper.RESTMappings(gk, versions...)
	}
	return ms, err
}

func (c *Cached) ResourceSingularizer(resource string) (singular string, err error) {
	s, err := c.mapper.ResourceSingularizer(resource)
	if c.onError(err) {
		s, err = c.mapper.ResourceSingularizer(resource)
	}
	return s, err
}
