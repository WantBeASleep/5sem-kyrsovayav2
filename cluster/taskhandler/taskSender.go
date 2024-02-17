package taskhandler

import (
	is "lib/infostructs"
	ts "lib/tasks"
	rq "lib/requests"
	tr "lib/trees"
	mt "lib/matrixes"
)

const (
	op_for_core int = 1e6
)

type TaskSender struct {
	OpCount int
	
	direction interface{}
	workersPool chan *is.WorkerInfo
	deferClusterWorkerTaskPool chan ts.ClusterWorkerTask
}

func NewTaskSender(workersPool chan *is.WorkerInfo, deferClusterWorkerTaskPool chan ts.ClusterWorkerTask) TaskSender {
	newTaskSender := TaskSender{
		workersPool: workersPool,
		deferClusterWorkerTaskPool: deferClusterWorkerTaskPool,
	}

	return newTaskSender
}

func (t *TaskSender) NewTask() {
	select {
	case t.direction = <-t.workersPool:
		t.OpCount = t.direction.(*is.WorkerInfo).Cores * op_for_core

	default:
		t.direction = t.deferClusterWorkerTaskPool
		t.OpCount = op_for_core
	}
}

func (t *TaskSender) WorkerDrop() {
	if x, ok := t.direction.(*is.WorkerInfo); ok {
		t.workersPool <- x
	}
}

func (t *TaskSender) Send(root tr.ASTNode, matrixes map[string]mt.Matrix, matrixesAlertReady *map[string]chan bool, resultMatrixName string) {
	necessaryMatrixes := tr.GetLeafsNames(root)
	sendMatrix := map[string]mt.Matrix{}

	for k := range necessaryMatrixes {
		if _, isNOTready := <-(*matrixesAlertReady)[k]; !isNOTready {
			sendMatrix[k] = matrixes[k]
			delete(necessaryMatrixes, k)
		}
	}

	newReq := rq.ClusterWorkerReq{
		Root: root,
		Matrixes: sendMatrix,
	}

	if x, ok := t.direction.(*is.WorkerInfo); len(necessaryMatrixes) == 0 && ok {
		wport := x.Port
		t.direction = nil
		go func(){
			var resultMatrix mt.Matrix
			rq.SendRequest(wport, "wsolveproblem", newReq, &resultMatrix)
			matrixes[resultMatrixName] = resultMatrix
			(*matrixesAlertReady)[resultMatrixName] <- true
			close((*matrixesAlertReady)[resultMatrixName])
		}()
	} else {
		t.WorkerDrop()
		newTask := ts.ClusterWorkerTask{
			CWR: newReq,
			ResultMatrixName: resultMatrixName,
			AllMatrixes: &matrixes,
			MatrixesAlertReady: matrixesAlertReady,
		}
		t.deferClusterWorkerTaskPool <- newTask
	}
}