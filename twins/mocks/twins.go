// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package mocks

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/mainflux/mainflux/twins"
)

var _ twins.TwinRepository = (*twinRepositoryMock)(nil)

type twinRepositoryMock struct {
	mu      sync.Mutex
	counter uint64
	twins   map[string]twins.Twin
}

// NewTwinRepository creates in-memory twin repository.
func NewTwinRepository() twins.TwinRepository {
	return &twinRepositoryMock{
		twins: make(map[string]twins.Twin),
	}
}

func (trm *twinRepositoryMock) Save(ctx context.Context, twin twins.Twin) (string, error) {
	trm.mu.Lock()
	defer trm.mu.Unlock()

	for _, tw := range trm.twins {
		if tw.ID == twin.ID {
			return "", twins.ErrConflict
		}
	}

	trm.twins[key(twin.Owner, twin.ID)] = twin

	return twin.ID, nil
}

func (trm *twinRepositoryMock) Update(ctx context.Context, twin twins.Twin) error {
	trm.mu.Lock()
	defer trm.mu.Unlock()

	dbKey := key(twin.Owner, twin.ID)
	if _, ok := trm.twins[dbKey]; !ok {
		return twins.ErrNotFound
	}

	trm.twins[dbKey] = twin

	return nil
}

func (trm *twinRepositoryMock) RetrieveByID(_ context.Context, id string) (twins.Twin, error) {
	trm.mu.Lock()
	defer trm.mu.Unlock()

	for k, v := range trm.twins {
		if id == v.ID {
			return trm.twins[k], nil
		}
	}

	return twins.Twin{}, twins.ErrNotFound
}

func (trm *twinRepositoryMock) RetrieveByAttribute(ctx context.Context, channel, subtopic string) ([]string, error) {
	var ids []string
	for _, twin := range trm.twins {
		def := twin.Definitions[len(twin.Definitions)-1]
		for _, attr := range def.Attributes {
			if attr.Channel == channel && attr.Subtopic == subtopic {
				ids = append(ids, twin.ID)
				break
			}
		}
	}

	return ids, nil
}

func (trm *twinRepositoryMock) RetrieveByThing(_ context.Context, thingid string) (twins.Twin, error) {
	trm.mu.Lock()
	defer trm.mu.Unlock()

	for _, twin := range trm.twins {
		if twin.ThingID == thingid {
			return twin, nil
		}
	}

	return twins.Twin{}, twins.ErrNotFound

}

func (trm *twinRepositoryMock) RetrieveAll(_ context.Context, owner string, offset uint64, limit uint64, name string, metadata twins.Metadata) (twins.TwinsPage, error) {
	trm.mu.Lock()
	defer trm.mu.Unlock()

	items := make([]twins.Twin, 0)

	if limit <= 0 {
		return twins.TwinsPage{}, nil
	}

	// This obscure way to examine map keys is enforced by the key structure in mocks/commons.go
	prefix := fmt.Sprintf("%s-", owner)
	for k, v := range trm.twins {
		if (uint64)(len(items)) >= limit {
			break
		}
		if !strings.HasPrefix(k, prefix) {
			continue
		}
		suffix := string(v.ID[len(u4Pref):])
		id, _ := strconv.ParseUint(suffix, 10, 64)
		if id > offset && id <= uint64(offset+limit) {
			items = append(items, v)
		}
	}

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})

	page := twins.TwinsPage{
		Twins: items,
		PageMetadata: twins.PageMetadata{
			Total:  trm.counter,
			Offset: offset,
			Limit:  limit,
		},
	}

	return page, nil
}

func (trm *twinRepositoryMock) Remove(ctx context.Context, id string) error {
	trm.mu.Lock()
	defer trm.mu.Unlock()

	for k, v := range trm.twins {
		if id == v.ID {
			delete(trm.twins, k)
		}
	}

	return nil
}
