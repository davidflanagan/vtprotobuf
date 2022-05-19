package vtproto

import "sync"

var (
	internedStrings map[string]string = make(map[string]string, 0)
	internMu        sync.RWMutex
)

func Intern(bytes []byte) string {
	internMu.RLock()
	s, ok := internedStrings[string(bytes)]
	internMu.RUnlock()
	if ok {
		return s
	}
	internMu.Lock()
	s, ok = internedStrings[string(bytes)]
	if !ok {
		s = string(bytes)
		internedStrings[s] = s
	}
	internMu.Unlock()
	return s
}
