package prng

import (
	"math/rand"
	"sync"
)

type SeededSource64 struct {
	mu    sync.Mutex
	state uint64
}

var _ (rand.Source64) = (*SeededSource64)(nil)

func NewSeededSource64(seed uint64) *SeededSource64 {
	return &SeededSource64{
		state: seed,
	}
}

// The underlying used seed.
func (s *SeededSource64) State() uint64 {
	return s.state
}

func (s *SeededSource64) Seed(seed int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state = uint64(seed)
}

func (s *SeededSource64) Seed64(seed uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state = seed
}

// Look for rand.Source64 interface
func (s *SeededSource64) Uint64() uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.state = s.state*6364136223846793005 + 1442695040888963407
	return s.state
}

func (s *SeededSource64) Int63() int64 {
	return int64(s.Uint64() & 0x7FFFFFFFFFFFFFFF)
}
