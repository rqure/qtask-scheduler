package main

import (
	"context"
	"sync/atomic"

	"github.com/rqure/qlib/pkg/qapp"
	"github.com/rqure/qlib/pkg/qdata"
)

type TaskScheduler interface {
	qapp.Worker

	OnReady(context.Context)
	OnNotReady(context.Context)
}

type taskScheduler struct {
	store   *qdata.Store
	isReady atomic.Bool

	// Main thread variables
	tokens []qdata.NotificationToken
}

func NewTaskScheduler(store *qdata.Store) TaskScheduler {
	return &taskScheduler{
		store:   store,
		isReady: atomic.Bool{},
	}
}

func (me *taskScheduler) Init(context.Context) {

}

func (me *taskScheduler) Deinit(context.Context) {

}

func (me *taskScheduler) DoWork(ctx context.Context) {

}

func (me *taskScheduler) OnReady(context.Context) {
	if !me.isReady.CompareAndSwap(false, true) {
		return // already ready
	}

}

func (me *taskScheduler) OnNotReady(context.Context) {
	if !me.isReady.CompareAndSwap(true, false) {
		return // already not ready
	}

}
