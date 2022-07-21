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

type DestinationDef struct {
	Name    string
	Drivers []DestinationDriver
}

type DestinationDriver interface {
	__DestinationDriver_union()
	Name() string
}

type DestinationDriverAlts interface {
	SyslogDestinationDriver
	Name() string
}

func NewDestinationDriver[Alt DestinationDriverAlts](alt Alt) DestinationDriver {
	return DestinationDriverAlt[Alt]{
		Alt: alt,
	}
}

type DestinationDriverAlt[Alt DestinationDriverAlts] struct {
	Alt Alt
}

func (DestinationDriverAlt[Alt]) __DestinationDriver_union() {}

func (alt DestinationDriverAlt[Alt]) Name() string {
	return alt.Alt.Name()
}

type SyslogDestinationDriver struct {
	Host           string
	Port           int
	Transport      string
	CADir          string
	CAFile         string
	CloseOnInput   *bool
	Flags          []string
	FlushLines     int
	SoKeepalive    *bool
	Suppress       int
	Template       string
	TemplateEscape *bool
	TLS            *SyslogDestinationDriverTLS
	TSFormat       string
	DiskBuffer     *DiskBufferDef
}

func (SyslogDestinationDriver) Name() string {
	return "syslog"
}

type SyslogDestinationDriverTLS struct {
}

type DiskBufferDef struct {
	DiskBufSize  float64
	Reliable     bool
	Compaction   *bool
	Dir          string
	MemBufLength int64
	MemBufSize   float64
	QOutSize     int64
}
