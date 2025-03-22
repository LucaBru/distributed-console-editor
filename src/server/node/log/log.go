package rlog

import (
	"editor-service/protos/logpb"

	"github.com/google/uuid"
)

type Log[T any, C any] struct {
	StateCh chan<- *LogResult[T]
	Cmd     C
}

type LogResult[T any] struct {
	Msg T
	Err error
}

type ShareLog = Log[*uuid.UUID, *logpb.Share]
type DeleteLog = Log[*struct{}, *logpb.Delete]
type EditLog = Log[*struct{}, *logpb.Edit]
