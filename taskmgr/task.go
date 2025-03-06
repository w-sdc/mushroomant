package taskmgr

import (
	"context"
	"io"
)

type procStatus int

const (
	procStatusRunning procStatus = iota
	procStatusDone
	procStatusError
)

type process struct {
	ctx       context.Context
	cmd       string
	args      []string
	envs      map[string]string
	pid       int
	status    procStatus
	retcode   int
	stdinSrc  io.Reader
	stdoutDst io.WriteCloser
	stderrDst io.WriteCloser
}
