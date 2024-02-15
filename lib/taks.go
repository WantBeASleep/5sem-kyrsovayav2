package lib

import (
	"sync/atomic"
)

type ClientTask struct {
	Expr string
	Data map[string]Matrix
}

type KlasterToWorkerTask struct {
	Root ASTNode
	Data map[string]Matrix
}

type DeferKlasterToWorkerTask struct {
	Task KlasterToWorkerTask
	IsMatrixReady bool
	NecessaryMatrixes map[string]bool
	ReadyMatrix map[string]*atomic.Bool
}

func (t DeferKlasterToWorkerTask) CheckIsReady() bool {
	for k := range t.NecessaryMatrixes {
		isReady := t.ReadyMatrix[k].Load()
		if isReady {
			delete(t.NecessaryMatrixes, k)
		}
	}
	return len(t.NecessaryMatrixes) == 0
}