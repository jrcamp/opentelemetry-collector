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
	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/collector/consumer/pdata"
)

type mockBytesTelemetryType struct {
	mock.Mock
}

func (m *mockBytesTelemetryType) EncodeTelemetryType(model interface{}) ([]byte, error) {
	args := m.Called(model)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *mockBytesTelemetryType) DecodeTelemetryType(bytes []byte) (interface{}, error) {
	args := m.Called(bytes)
	return args.Get(0), args.Error(1)
}

type mockModelTelemetryType struct {
	mock.Mock
	typ interface{}
}

func (m *mockModelTelemetryType) ToTelemetryType(src interface{}) (pdata.TelemetryType, error) {
	args := m.Called(src)
	return args.Get(0).(pdata.TelemetryType), args.Error(1)
}

func (m *mockModelTelemetryType) FromTelemetryType(md pdata.TelemetryType, out interface{}) error {
	args := m.Called(md, out)
	return args.Error(0)
}

func (m *mockModelTelemetryType) Type() interface{} {
	return m.typ
}
