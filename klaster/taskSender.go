package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"lib"
	"net/http"
	"sync/atomic"
)

const (
	op_for_core int = 1e6
)

type TaskSender struct {
	OpCount int
	
	direction interface{}
	workersPool chan *lib.Worker
	deferTaskQ chan lib.DeferKlasterToWorkerTask
}

func NewTaskSender(workersPool chan *lib.Worker, deferTasksPool chan lib.DeferKlasterToWorkerTask) TaskSender {
	newTaskSender := TaskSender{
		workersPool: workersPool,
		deferTaskQ: deferTasksPool,
	}

	return newTaskSender
}

func (t *TaskSender) NewTask() {
	select {
	case t.direction = <-t.workersPool:
		t.OpCount = t.direction.(*lib.Worker).Cores * op_for_core

	default:
		t.direction = t.deferTaskQ
		t.OpCount = op_for_core
	}
}

func (t *TaskSender) Send(root lib.ASTNode, clientmatrixes map[string]lib.Matrix, matrixesStatus map[string]*atomic.Bool, subTreeName string) {
	necessaryMatrix := lib.GetLeafsNames(root)
	sendMatrix := map[string]lib.Matrix{}

	for k := range necessaryMatrix {
		isReady := matrixesStatus[k].Load()
		if isReady {
			sendMatrix[k] = clientmatrixes[k]
			delete(necessaryMatrix, k)
		}
	}

	newTask := lib.KlasterToWorkerTask{
		Root: root,
		Data: sendMatrix,
	}

	if x, ok := t.direction.(*lib.Worker); len(necessaryMatrix) == 0 && ok {
		newTaskJson, _ := json.Marshal(newTask)
		go func(){
			response, _ := http.Post(fmt.Sprintf("http://localhost:%s/subtreesolver", x.Port), "application/json", bytes.NewReader(newTaskJson))
			defer response.Body.Close()
			var responseMatrix lib.Matrix
			json.NewDecoder(response.Body).Decode(&response)
			clientmatrixes[subTreeName] = responseMatrix
			matrixesStatus[subTreeName].Store(true)
		}()
	} else {
		if x, ok := t.direction.(*lib.Worker); ok {
			t.workersPool <- x
			t.direction = nil
		}

		deferTask := lib.DeferKlasterToWorkerTask{
			Task: newTask,
			IsMatrixReady: len(necessaryMatrix) == 0,
			NecessaryMatrixes: necessaryMatrix,
			ReadyMatrix: matrixesStatus,
		}

		t.deferTaskQ <- deferTask
	}
}