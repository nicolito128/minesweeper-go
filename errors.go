package minesweeper

import "errors"

var (
	ErrInvalidSize       = errors.New("size must be a positive number")
	ErrInvalidTotalMines = errors.New("total mines must be a positive number")
	ErrOutOfBounds       = errors.New("coordinates out of bounds")
	ErrInvalidAction     = errors.New("invalid action")
)
