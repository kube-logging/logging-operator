// Copyright © 2019 Banzai Cloud
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

package model

type SourceDef struct {
	Name    string
	Drivers []SourceDriver
}

type SourceDriver interface {
	__SourceDriver_union()
	Name() string
}

type SourceDriverAlts interface {
	NetworkSourceDriver
	Name() string
}

func NewSourceDriver[Alt SourceDriverAlts](alt Alt) SourceDriver {
	return SourceDriverAlt[Alt]{
		Alt: alt,
	}
}

type SourceDriverAlt[Alt SourceDriverAlts] struct {
	Alt Alt
}

func (SourceDriverAlt[Alt]) __SourceDriver_union() {}

func (alt SourceDriverAlt[Alt]) Name() string {
	return alt.Alt.Name()
}

type NetworkSourceDriver struct {
	Flags     []string
	IP        string
	Port      uint16
	Transport string
}

func (NetworkSourceDriver) Name() string {
	return "network"
}
