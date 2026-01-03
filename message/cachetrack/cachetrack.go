// Package cachetrack provides a CacheTracker interface and a Tracker
// implementation that track when a cache has been invalidated.
package cachetrack

import "github.com/rothskeller/packet/message/msgifc"

// CacheTracker is an interface satisfied by an object with cache tracking.
type CacheTracker = msgifc.CacheTracker

// Tracker is a cache tracker that keeps track of whether a cache is dirty.
type Tracker struct {
	dirty    bool
	handlers []func(string)
}

var _ CacheTracker = (*Tracker)(nil)

// Dirty returns whether the cache is dirty (i.e., invalid).
func (t *Tracker) Dirty() bool { return t.dirty }

// MarkClean marks the cache as clean.
func (t *Tracker) MarkClean() { t.dirty = false }

// OnDirty registers a function to be called when the cache becomes dirty.
func (t *Tracker) OnDirty(fn func(string)) { t.handlers = append(t.handlers, fn) }

// MarkDirty marks the cache as dirty.
func (t *Tracker) MarkDirty(reason string) {
	if !t.dirty {
		t.dirty = true
		for _, fn := range t.handlers {
			fn(reason)
		}
	}
}
