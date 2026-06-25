package minesweeper

import (
	"errors"
	"fmt"
	"strings"
)

type ActionKind byte

const (
	ActionToggleFlag ActionKind = 'F'
	ActionRevealCell ActionKind = 'C'
)

func (a ActionKind) String() string {
	switch a {
	case ActionToggleFlag:
		return "F"
	case ActionRevealCell:
		return "C"
	default:
		return " "
	}
}

type Action struct {
	Kind ActionKind
	X, Y int
}

func NewAction(kind ActionKind, x, y int) Action {
	return Action{Kind: kind, X: x, Y: y}
}

func (a Action) String() string {
	return fmt.Sprintf("%c%d;%d.", a.Kind, a.X, a.Y)
}

func ParseAction(s string) (*Action, error) {
	// <Letter><PosX>;<PosY>.
	format := "%c%d;%d"
	if len(format) < 4 {
		return nil, errors.New("failed to parse action with less than 4 characters")
	}
	s = strings.ToUpper(s)
	result := new(Action)
	_, err := fmt.Sscanf(s, format, &result.Kind, &result.X, &result.Y)
	if err != nil {
		return nil, fmt.Errorf("failed to parse action: %v", err)
	}
	return result, nil
}

func ParseActions(s string) ([]*Action, error) {
	s = strings.TrimSpace(s)
	result := make([]*Action, 0)
	for actionStr := range strings.SplitSeq(s, ".") {
		if actionStr != "" {
			resultStr := strings.TrimSpace(actionStr)
			act, err := ParseAction(resultStr)
			if err != nil {
				return nil, err
			}
			result = append(result, act)
		}
	}
	return result, nil
}
