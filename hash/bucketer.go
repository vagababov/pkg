/*
Copyright 2020 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// This file contains the utilities to make bucketing decisions.

package hash

import (
	"sync"

	lru "github.com/hashicorp/golang-lru"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
)

// Bucket answers whether key X belongs to the bucket in
// consistent manner (consistent as in consistent hashing).
// Bucket implements reconciler.Bucket interface.
type Bucket struct {
	// The name of this bucket.
	name string
	// Stores the cached lookups. cache is internally thread safe.
	cache *lru.Cache

	// mu guards buckets.
	mu sync.RWMutex
	// All the bucket names. Needed for building hash universe.
	// `name` must be in this set.
	buckets sets.String
}

// Has returns true if this bucket owns the key.
func (b *Bucket) Has(nn types.NamespacedName) bool {
	return b.has(nn.String())
}

func (b *Bucket) has(key string) bool {
	if v, ok := b.cache.Get(key); ok {
		return v.(bool)
	}
	b.mu.RLock()
	defer b.mu.RUnlock()
	l := ChooseSubset(b.buckets, 1 /*single target*/, key)
	ret := (l.UnsortedList()[0] == b.name)
	b.cache.Add(key, ret)
	return ret
}

// Update updates the universe of buckets.
func (b *Bucket) Update(newB sets.String) {
	b.mu.Lock()
	defer b.mu.Unlock()
	// In theory we can iterate over the map and
	// purge only the keys moved to a new shard.
	// But this might be more expensive than re-build
	// the cache as reconciliations happen.
	b.cache.Purge()
	b.buckets = newB
}
