package workpool

import (
	ts "lib/tasks"
	rq "lib/requests"
	is "lib/infostructs"
	mt "lib/matrixes"
)

var taskQ chan ts.Wpooltask

func StartPool(workerinfo *is.WorkerInfo) {
	taskQ = make(chan ts.Wpooltask, 1000000)
	for range workerinfo.Cores {
		go worker()
	}
}

func worker() {
	for {
		newTask := <-taskQ
		if !newTask.CanStartTask() {
			taskQ <- newTask
		} else {
			newTask.Start()
		}
	}
}

func DoExpr(newTask rq.ClusterWorkerReq) mt.Matrix {
	return <-ParseTreeToTasks(taskQ, newTask.Root, newTask.Matrixes)
}