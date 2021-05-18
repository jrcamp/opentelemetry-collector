// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:generate genny -in=$GOFILE -out=gen-$GOFILE gen "TelemetryType=Metrics,Traces,Logs"

package protocols

import (
	"go.opentelemetry.io/collector/consumer/pdata"
	"go.opentelemetry.io/collector/protocols/bytes"
	"go.opentelemetry.io/collector/protocols/models"
)

type TelemetryTypeEncoder struct {
	mod models.TelemetryTypeEncoder
	enc bytes.TelemetryTypeEncoder
}

// EncodeT encodes pdata to bytes.
func (t *TelemetryTypeEncoder) EncodeT(td pdata.TelemetryType) ([]byte, error) {
	out := t.mod.Type()
	if err := t.mod.FromTelemetryType(td, &out); err != nil {
		return nil, err
	}
	return t.enc.EncodeTelemetryType(out)
}
