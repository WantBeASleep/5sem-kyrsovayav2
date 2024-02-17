package taskhandler

import (
	ts "lib/tasks"
	mt "lib/matrixes"
	rq "lib/requests"
	is "lib/infostructs"
)

func DeferTasksPoolHandler(deferClusterWorkerTaskPool chan ts.ClusterWorkerTask, workersPool chan *is.WorkerInfo) {
	for {
		deferTask := <- deferClusterWorkerTaskPool
		if !deferTask.CheckReady() {
			deferClusterWorkerTaskPool <- deferTask
			continue
		}

		worker := <-workersPool
		go func(){
			var resultMatrix mt.Matrix
			rq.SendRequest(worker.Port, "wsolveproblem", deferTask.CWR, &resultMatrix)
			(*deferTask.AllMatrixes)[deferTask.ResultMatrixName] = resultMatrix
			close((*deferTask.MatrixesAlertReady)[deferTask.ResultMatrixName])
		}()
	}
}