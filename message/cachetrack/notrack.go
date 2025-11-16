package cachetrack

// NoTracker is an implementation of CacheTracker for objects that are
// immutable and therefore don't need caching.
type NoTracker struct{}

var _ CacheTracker = (*NoTracker)(nil)

// Dirty returns whether the cache is dirty (i.e., invalid).
func (t *NoTracker) Dirty() bool { return false }

// MarkClean marks the cache as clean.
func (t *NoTracker) MarkClean() {}

// OnDirty registers a function to be called when the cache becomes dirty.
func (t *NoTracker) OnDirty(_ func(string)) {}

// MarkDirty marks the cache as dirty.
func (t *NoTracker) MarkDirty(string) {
	panic("cannot call MarkDirty on an immutable object")
}
