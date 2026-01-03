package msgifc

// CacheTracker is an interface satisfied by an object with cache tracking.
type CacheTracker interface {
	// Dirty returns whether the cache is dirty (i.e., invalid).
	Dirty() bool
	// MarkDirty marks the cache as dirty.
	MarkDirty(string)
	// MarkClean marks the cache as clean.
	MarkClean()
	// OnDirty registers a function to be called when the cache becomes
	// dirty.
	OnDirty(func(string))
}
