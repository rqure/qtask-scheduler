package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/rqure/qlib/pkg/qapp"
	"github.com/rqure/qlib/pkg/qapp/qworkers"
	"github.com/rqure/qlib/pkg/qdata"
	"github.com/rqure/qlib/pkg/qlog"
)

type TaskScheduler interface {
	qapp.Worker

	OnReady(context.Context)
	OnNotReady(context.Context)

	OnEntityCreated(qworkers.EntityCreateArgs)
	OnEntityDeleted(qworkers.EntityDeletedArgs)
}

type taskScheduler struct {
	store     *qdata.Store
	isReady   bool
	mu        sync.RWMutex
	tokens    []qdata.NotificationToken
	scheduler *gocron.Scheduler
}

func NewTaskScheduler(store *qdata.Store) TaskScheduler {
	return &taskScheduler{
		store:     store,
		isReady:   false,
		mu:        sync.RWMutex{},
		tokens:    make([]qdata.NotificationToken, 0),
		scheduler: gocron.NewScheduler(time.UTC),
	}
}

func (me *taskScheduler) Init(context.Context) {

}

func (me *taskScheduler) Deinit(context.Context) {

}

func (me *taskScheduler) DoWork(ctx context.Context) {

}

func (me *taskScheduler) OnReady(ctx context.Context) {
	me.mu.Lock()
	defer me.mu.Unlock()

	if me.isReady {
		return // already ready
	}

	me.isReady = true

	for _, token := range me.tokens {
		token.Unbind(ctx)
	}

	me.tokens = make([]qdata.NotificationToken, 0)
}

func (me *taskScheduler) OnNotReady(context.Context) {
	me.mu.Lock()
	defer me.mu.Unlock()

	if !me.isReady {
		return // already not ready
	}

	me.isReady = false
}

func (me *taskScheduler) OnEntityCreated(args qworkers.EntityCreateArgs) {
	me.mu.RLock()
	defer me.mu.RUnlock()

	if !me.isReady {
		return // not ready, ignore
	}

	// Handle entity creation logic here
}

func (me *taskScheduler) OnEntityDeleted(args qworkers.EntityDeletedArgs) {
	me.mu.RLock()
	defer me.mu.RUnlock()

	if !me.isReady {
		return // not ready, ignore
	}

	// Handle entity deletion logic here
}

func (me *taskScheduler) registerAllTasks(ctx context.Context) error {
	tasks, err := me.store.Find(ctx, "Task", []qdata.FieldType{})
	if err != nil {
		return fmt.Errorf("failed to find tasks: %w", err)
	}

	reqs := make([]*qdata.Request, len(tasks))
	for _, task := range tasks {
		task.Field("ConfigurationStatus").Value.FromChoice(1) // Valid
		if err := me.registerTask(ctx, task); err != nil {
			task.Field("ConfigurationStatus").Value.FromChoice(0) // Invalid
			qlog.Warn("failed to register task %s: %v", task.EntityId, err)
		}
		reqs = append(reqs, task.Field("ConfigurationStatus").AsWriteRequest())
	}

	return nil
}

func (me *taskScheduler) unregisterAllTasks(ctx context.Context) error {
	// Unregister all tasks here
	return nil
}

func (me *taskScheduler) registerTask(ctx context.Context, task *qdata.Entity) error {
	return nil
}

func (me *taskScheduler) unregisterTask(ctx context.Context, task *qdata.Entity) error {
	return nil
}

func (me *taskScheduler) registerCronTask(ctx context.Context, task *qdata.Entity) error {
	return nil
}

func (me *taskScheduler) unregisterCronTask(ctx context.Context, task *qdata.Entity) error {
	return nil
}

func (me *taskScheduler) registerNotificationTask(ctx context.Context, task *qdata.Entity) error {
	return nil
}

func (me *taskScheduler) unregisterNotificationTask(ctx context.Context, task *qdata.Entity) error {
	return nil
}
