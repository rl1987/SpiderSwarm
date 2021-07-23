package main

import (
	"github.com/google/uuid"
)

type NullAction struct {
	AbstractAction
}

func NewNullAction() *NullAction {
	return &NullAction{
		AbstractAction: AbstractAction{
			UUID: uuid.New().String(),
		},
	}
}
