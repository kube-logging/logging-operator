// Copyright © 2023 Kube logging authors
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

package common

import (
	"fmt"
	"os"
	"sync/atomic"
	"testing"

	"emperror.dev/errors"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
)

var sequence uint32

func RequireNoError(t *testing.T, err error) {
	if err != nil {
		assert.Fail(t, fmt.Sprintf("Received unexpected error:\n%#v %+v", err, errors.GetDetails(err)))
		t.FailNow()
	}
}

func Initialize(t *testing.T) {
	localSeq := atomic.AddUint32(&sequence, 1)
	shards := cast.ToUint32(os.Getenv("SHARDS"))
	shard := cast.ToUint32(os.Getenv("SHARD"))
	if shards > 0 {
		if localSeq%shards != shard {
			t.Skipf("skipping %s as sequence %d not in shard %d", t.Name(), localSeq, shard)
		}
	}
	t.Parallel()
}
