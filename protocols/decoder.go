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
	"github.com/cheekybits/genny/generic"
	"go.opentelemetry.io/collector/consumer/pdata"
	"go.opentelemetry.io/collector/protocols/bytes"
	"go.opentelemetry.io/collector/protocols/models"
)

type TelemetryType generic.Type

type TelemetryTypeDecoder struct {
	mod models.TelemetryTypeDecoder
	enc bytes.TelemetryTypeDecoder
}

// DecodeTelemetryType decodes bytes to pdata.
func (t *TelemetryTypeDecoder) DecodeTelemetryType(data []byte) (pdata.TelemetryType, error) {
	model, err := t.enc.DecodeTelemetryType(data)
	if err != nil {
		return pdata.NewTelemetryType(), err
	}
	return t.mod.ToTelemetryType(model)
}
