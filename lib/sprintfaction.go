package spsw

import (
//"fmt"

//"github.com/google/uuid"
)

type SprintfAction struct {
	AbstractAction
	FormatString string
	Arguments    []string
}
