package domain

import "errors"

var (
	ErrInvalidState = errors.New("invalid state")
)

var (
	StateLocked  = State{"Locked"}
	StatePrivate = State{"Private"}
	StateShared  = State{"Shared"}
	StateTrashed = State{"Trashed"}
	StateDeleted = State{"Deleted"}
)

type State struct {
	value string
}

func CreateState(value string) (State, error) {
	if value == StateLocked.value {
		return StateLocked, nil
	}
	if value == StatePrivate.value {
		return StatePrivate, nil
	}
	if value == StateShared.value {
		return StateShared, nil
	}
	if value == StateTrashed.value {
		return StateTrashed, nil
	}
	if value == StateDeleted.value {
		return StateDeleted, nil
	}

	return State{}, ErrInvalidState
}

func (state State) IsZero() bool {
	return state == State{}
}

func (state State) Value() string {
	return state.value
}
