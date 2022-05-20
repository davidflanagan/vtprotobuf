package vtproto

import "sync"

var (
	internedStrings map[string]string = make(map[string]string, 0)
	internMu        sync.Mutex
)

// Intern is the default implementation of string interning used by generated unmarshaling
// code when you annotate a string typed field with the [(vtproto.intern) = true] option.
// If you want to use a different implementation (which must have the same signature as
// this function) you can use protoc options like these:
//
//   --go-vtproto_opt=internPkg=github.com/josharian/intern
//   --go-vtproto_opt=internFunc=Bytes
//
// Note that this is an intentionally bare-bones implementation. You might want to use a
// different implementation with different concurrency handling or one that caps the
// maximum number of interned strings (to prevent DOS attacks, for example)
func Intern(bytes []byte) string {
	internMu.Lock()
	// The Go compiler optimizes the string(bytes) as the map key below
	// so that no memory allocation is done during the map lookup
	s, ok := internedStrings[string(bytes)]
	if ok {
		internMu.Unlock()
		return s
	}

	// If the string does not already exist in the map, this is where we
	// allocate the string and intern it. Any subsequent lookups of this
	// sequence of byte will be allocation-free.
	s = string(bytes)
	internedStrings[s] = s
	internMu.Unlock()
	return s
}
