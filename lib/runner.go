package spsw

import (
//"github.com/davecgh/go-spew/spew"
//log "github.com/sirupsen/logrus"
)

type Runner struct {
	BackendAddr string
}

func NewRunner(backendAddr string) *Runner {
	return &Runner{BackendAddr: backendAddr}
}
