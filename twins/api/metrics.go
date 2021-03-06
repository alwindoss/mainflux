// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// +build !test

package api

import (
	"context"
	"time"

	"github.com/go-kit/kit/metrics"
	"github.com/mainflux/mainflux/broker"
	"github.com/mainflux/mainflux/twins"
)

var _ twins.Service = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     twins.Service
}

// MetricsMiddleware instruments core service by tracking request count and
// latency.
func MetricsMiddleware(svc twins.Service, counter metrics.Counter, latency metrics.Histogram) twins.Service {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}

func (ms *metricsMiddleware) AddTwin(ctx context.Context, token string, twin twins.Twin, def twins.Definition) (saved twins.Twin, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "add_twin").Add(1)
		ms.latency.With("method", "add_twin").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.AddTwin(ctx, token, twin, def)
}

func (ms *metricsMiddleware) UpdateTwin(ctx context.Context, token string, twin twins.Twin, def twins.Definition) (err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "update_twin").Add(1)
		ms.latency.With("method", "update_twin").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.UpdateTwin(ctx, token, twin, def)
}

func (ms *metricsMiddleware) ViewTwin(ctx context.Context, token, id string) (viewed twins.Twin, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "view_twin").Add(1)
		ms.latency.With("method", "view_twin").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ViewTwin(ctx, token, id)
}

func (ms *metricsMiddleware) ListTwins(ctx context.Context, token string, offset uint64, limit uint64, name string, metadata twins.Metadata) (tw twins.TwinsPage, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_twins").Add(1)
		ms.latency.With("method", "list_twins").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ListTwins(ctx, token, offset, limit, name, metadata)
}

func (ms *metricsMiddleware) SaveStates(msg *broker.Message) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "save_states").Add(1)
		ms.latency.With("method", "save_states").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.SaveStates(msg)
}

func (ms *metricsMiddleware) ListStates(ctx context.Context, token string, offset uint64, limit uint64, id string) (st twins.StatesPage, err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_states").Add(1)
		ms.latency.With("method", "list_states").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ListStates(ctx, token, offset, limit, id)
}

func (ms *metricsMiddleware) ViewTwinByThing(ctx context.Context, token, thingid string) (twins.Twin, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "view_twin_by_thing").Add(1)
		ms.latency.With("method", "view_twin_by_thing").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ViewTwinByThing(ctx, token, thingid)
}

func (ms *metricsMiddleware) RemoveTwin(ctx context.Context, token, id string) (err error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "remove_twin").Add(1)
		ms.latency.With("method", "remove_twin").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.RemoveTwin(ctx, token, id)
}
