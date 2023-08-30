package envexpander

import "sync"

// CachedVariablePos is a scoped cache for ExtractAllVariables.
type CachedVariablePos struct {
	cache map[string][]VariablePos
	lock  *sync.RWMutex
}

// NewCachedVariablePos creates a new CachedVariablePos.
func NewCachedVariablePos() CachedVariablePos {
	return CachedVariablePos{
		cache: make(map[string][]VariablePos),
		lock:  &sync.RWMutex{},
	}
}

// read returns the cached variable positions.
func (c CachedVariablePos) read(value string) ([]VariablePos, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	cached, ok := c.cache[value]
	return cached, ok
}

// write writes the cached variable positions.
func (c CachedVariablePos) write(value string, variables []VariablePos) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.cache[value] = variables
}

// MarkVariablePositions returns the cached variable positions.
func (c CachedVariablePos) MarkVariablePositions(value string) []VariablePos {
	if cached, ok := c.read(value); ok {
		return cached
	}

	v := MarkVariablePositions(value)
	c.write(value, v)

	return v
}
