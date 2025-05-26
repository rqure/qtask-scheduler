package main

import (
	"github.com/rqure/qlib/pkg/qapp"
	"github.com/rqure/qlib/pkg/qapp/qworkers"
	"github.com/rqure/qlib/pkg/qdata/qstore"
)

func main() {
	store := qstore.New()

	// qlib workers
	storeWorker := qworkers.NewStore(store)
	readinessWorker := qworkers.NewReadiness()
	readinessWorker.AddCriteria(qworkers.NewStoreConnectedCriteria(storeWorker, readinessWorker))

	// qtask workers
	taskScheduler := NewTaskScheduler(store)
	readinessWorker.BecameReady().Connect(taskScheduler.OnReady)
	readinessWorker.BecameNotReady().Connect(taskScheduler.OnNotReady)

	app := qapp.NewApplication("qtask-scheduler")
	app.AddWorker(storeWorker)
	app.AddWorker(readinessWorker)
	app.Execute()
}
