package minesweeper

import "errors"

var (
	ErrInvalidSize   = errors.New("size must be positive number")
	ErrOutOfBounds   = errors.New("coordinates out of bounds")
	ErrInvalidAction = errors.New("invalid action")
)
