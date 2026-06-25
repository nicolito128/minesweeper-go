package prng

import (
	"math/rand"
	"sync"
)

type SeededSource struct {
	mu    sync.Mutex
	state int64
}

var _ (rand.Source) = (*SeededSource)(nil)

func NewSeededSource(seed int64) *SeededSource {
	return &SeededSource{
		state: seed,
	}
}

func (s *SeededSource) Seed(seed int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state = seed
}

func (s *SeededSource) Int63() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	// change the state predictably but deterministically
	s.state = s.state*1103515245 + 12345

	// clear the sign bit so it is non-negative
	return s.state & 0x7FFFFFFFFFFFFFFF
}
