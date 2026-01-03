package msgifc

// Subject is the interface satisfied by all subject types.
type Subject interface {
	CacheTracker
	// EncodedSubject returns the encoded subject line.
	EncodedSubject() string
	// Clone returns a copy of the subject.
	Clone() Subject
}
