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

package tap

import (
	"context"
	"sync"

	"go.uber.org/atomic"
	"go.uber.org/zap"

	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/consumer/pdata"
)

var (
	_ consumer.Metrics = (*Tap)(nil)
	_ consumer.Traces  = (*Tap)(nil)
	_ consumer.Logs    = (*Tap)(nil)
)

type Tap struct {
	sync.Mutex
	Logger *zap.Logger
	MC     consumer.Metrics
	TC     consumer.Traces
	LC     consumer.Logs

	mcs    []consumer.Metrics
	numMcs atomic.Int32
}

func (t *Tap) ConsumeLogs(ctx context.Context, ld pdata.Logs) error {
	// TODO
	panic("implement me")
}

func (t *Tap) ConsumeTraces(ctx context.Context, td pdata.Traces) error {
	// TODO
	panic("implement me")
}

func (t *Tap) ConsumeMetrics(ctx context.Context, md pdata.Metrics) error {
	if t.numMcs.Load() == 0 {
		return t.MC.ConsumeMetrics(ctx, md)
	}

	t.Lock()

	for i, mc := range t.mcs {
		if err := mc.ConsumeMetrics(ctx, md); err != nil {
			t.Logger.Error("send to consumer from tap failed", zap.Int("consumer", i), zap.Error(err))
		}
	}

	t.Unlock()

	return t.MC.ConsumeMetrics(ctx, md)
}

func (t *Tap) RegisterMetricsConsumer(ctx context.Context, mc consumer.Metrics) {
	t.Lock()
	defer t.Unlock()

	t.numMcs.Inc()
	t.mcs = append(t.mcs, mc)
}

func (t *Tap) UnregisterMetricsConsumer(ctx context.Context, mc consumer.Metrics) bool {
	t.Lock()
	defer t.Unlock()

	t.numMcs.Dec()

	// TODO: do we care if multiple registers were called on the same mc and they get unregistered in different orders?
	for i, metricsConsumer := range t.mcs {
		if metricsConsumer == mc {
			t.mcs = append(t.mcs[:i], t.mcs[i+1:]...)
			return true
		}
	}

	return false
}
