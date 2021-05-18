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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/collector/consumer/pdata"
)

func TestTelemetryTypeDecoder_DecodeTelemetryTypeError(t *testing.T) {
	model := &mockModelTelemetryType{typ: ""}
	bytes := &mockBytesTelemetryType{}

	d := &TelemetryTypeDecoder{
		mod: model,
		enc: bytes,
	}

	//expectedTelemetryType := pdata.NewTelemetryType()
	expectedBytes := []byte{1, 2, 3}
	expectedModel := struct{}{}
	expectedError := errors.New("decode failed")

	bytes.On("DecodeTelemetryType", expectedBytes).Return(expectedModel, expectedError)
	//model.On("ToTelemetryType", expectedModel).Return(expectedTelemetryType, nil)

	_, err := d.DecodeTelemetryType(expectedBytes)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	//assert.Equal(t, expectedTelemetryType, actualTelemetryType)
}

func TestTelemetryTypeDecoder_DecodeTelemetryType(t *testing.T) {
	model := &mockModelTelemetryType{typ: ""}
	bytes := &mockBytesTelemetryType{}

	d := &TelemetryTypeDecoder{
		mod: model,
		enc: bytes,
	}

	expectedTelemetryType := pdata.NewTelemetryType()
	expectedBytes := []byte{1, 2, 3}
	expectedModel := struct{}{}

	bytes.On("DecodeTelemetryType", expectedBytes).Return(expectedModel, nil)
	model.On("ToTelemetryType", expectedModel).Return(expectedTelemetryType, nil)

	actualTelemetryType, err := d.DecodeTelemetryType(expectedBytes)

	assert.NoError(t, err)
	assert.Equal(t, expectedTelemetryType, actualTelemetryType)
}