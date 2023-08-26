package storage

import "errors"

var (
	ErrSegmentExists     = errors.New("segment exists")
	ErrNothingDelete     = errors.New("no such segment found")
	ErrNothingDeleteUser = errors.New("no such user found")
)
