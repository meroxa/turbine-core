// Copyright © 2022 Meroxa, Inc.
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

package record

import (
	"github.com/conduitio/conduit-commons/opencdc"
	procproto "github.com/conduitio/conduit-connector-protocol/proto/opencdc/v1"
)

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	var cTypes [1]struct{}
	_ = cTypes[int(opencdc.OperationCreate)-int(procproto.Operation_OPERATION_CREATE)]
	_ = cTypes[int(opencdc.OperationUpdate)-int(procproto.Operation_OPERATION_UPDATE)]
	_ = cTypes[int(opencdc.OperationDelete)-int(procproto.Operation_OPERATION_DELETE)]
	_ = cTypes[int(opencdc.OperationSnapshot)-int(procproto.Operation_OPERATION_SNAPSHOT)]
}
