package tasks

import (
	mt "lib/matrixes"
	rq "lib/requests"
	tr "lib/trees"
)

type ClusterWorkerTask struct {
	CWR                rq.ClusterWorkerReq
	ResultMatrixName   string
	AllMatrixes        *map[string]mt.Matrix
	MatrixesAlertReady *map[string]chan bool
}

func (task *ClusterWorkerTask) CheckReady() bool {
	necessaryMatrixes := tr.GetLeafsNames(task.CWR.Root)
	for name := range necessaryMatrixes {
		if _, isOpen := <-(*task.MatrixesAlertReady)[name]; isOpen {
			return false
		}
	}
	return true
}
