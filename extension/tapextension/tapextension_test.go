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

package tapextension

import (
	"context"
	"net"
	"net/http"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config/confignet"
	"go.opentelemetry.io/collector/testutil"
)

func TesttapExtensionUsage(t *testing.T) {
	config := Config{
		TCPAddr: confignet.TCPAddr{
			Endpoint: testutil.GetAvailableLocalAddress(t),
		},
	}

	tapExt := newServer(config, zap.NewNop())
	require.NotNil(t, tapExt)

	require.NoError(t, tapExt.Start(context.Background(), componenttest.NewNopHost()))
	defer tapExt.Shutdown(context.Background())

	// Give a chance for the server goroutine to run.
	runtime.Gosched()

	_, tapPort, err := net.SplitHostPort(config.TCPAddr.Endpoint)
	require.NoError(t, err)

	client := &http.Client{}
	resp, err := client.Get("http://localhost:" + tapPort + "/debug/tracez")
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TesttapExtensionPortAlreadyInUse(t *testing.T) {
	endpoint := testutil.GetAvailableLocalAddress(t)
	ln, err := net.Listen("tcp", endpoint)
	require.NoError(t, err)
	defer ln.Close()

	config := Config{
		TCPAddr: confignet.TCPAddr{
			Endpoint: endpoint,
		},
	}
	tapExt := newServer(config, zap.NewNop())
	require.NotNil(t, tapExt)

	require.Error(t, tapExt.Start(context.Background(), componenttest.NewNopHost()))
}

func TesttapMultipleStarts(t *testing.T) {
	config := Config{
		TCPAddr: confignet.TCPAddr{
			Endpoint: testutil.GetAvailableLocalAddress(t),
		},
	}

	tapExt := newServer(config, zap.NewNop())
	require.NotNil(t, tapExt)

	require.NoError(t, tapExt.Start(context.Background(), componenttest.NewNopHost()))
	defer tapExt.Shutdown(context.Background())

	// Try to start it again, it will fail since it is on the same endpoint.
	require.Error(t, tapExt.Start(context.Background(), componenttest.NewNopHost()))
}

func TesttapMultipleShutdowns(t *testing.T) {
	config := Config{
		TCPAddr: confignet.TCPAddr{
			Endpoint: testutil.GetAvailableLocalAddress(t),
		},
	}

	tapExt := newServer(config, zap.NewNop())
	require.NotNil(t, tapExt)

	require.NoError(t, tapExt.Start(context.Background(), componenttest.NewNopHost()))
	require.NoError(t, tapExt.Shutdown(context.Background()))
	require.NoError(t, tapExt.Shutdown(context.Background()))
}

func TesttapShutdownWithoutStart(t *testing.T) {
	config := Config{
		TCPAddr: confignet.TCPAddr{
			Endpoint: testutil.GetAvailableLocalAddress(t),
		},
	}

	tapExt := newServer(config, zap.NewNop())
	require.NotNil(t, tapExt)

	require.NoError(t, tapExt.Shutdown(context.Background()))
}
