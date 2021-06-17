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

package service

import (
	"flag"
	"fmt"
	"net"
	"net/http"

	"go.uber.org/zap"

	"go.opentelemetry.io/collector/component"
)

var statusAddr *string

// Flags adds flags related to basic building of the collector server to the given flagset.
func statusFlags(flags *flag.FlagSet) {
	statusAddr = flags.String("status-addr", "localhost:55679",
		fmt.Sprintf("Flag to specify HTTP status listen address."))
}

type registerDebug interface {
	RegisterDebug(mux *http.ServeMux)
}

type status struct {
	logger *zap.Logger
	server http.Server
	stopCh chan struct{}
	addr   string
	host   component.Host
}

func newStatusServer(logger *zap.Logger, addr string, host component.Host) *status {
	return &status{logger: logger, stopCh: make(chan struct{}), addr: addr, host: host}
}

func (s *status) start(handler http.Handler) error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	s.logger.Info("Starting HTTP status server", zap.Any("addr", s.addr))
	s.server = http.Server{Handler: handler}
	go func() {
		defer close(s.stopCh)

		if err := s.server.Serve(ln); err != nil && err != http.ErrServerClosed {
			s.host.ReportFatalError(err)
		}
	}()

	return nil
}

func (s *status) shutdown() error {
	err := s.server.Close()
	if s.stopCh != nil {
		<-s.stopCh
	}
	return err
}
