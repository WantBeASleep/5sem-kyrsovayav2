package main

import (
	"lib"
	"net/http"
	"fmt"
	"encoding/json"
	"os"
)

var klasterStatus *lib.Klaster
var workersPool chan *lib.Worker
var deferTasksPool chan lib.DeferKlasterToWorkerTask


const WORKTASK_PER_CORE int = 1e6
const DEFAULT_WORKTASK_SIZE int = 1e6

func addWorkerHandle(w http.ResponseWriter, r *http.Request) {
	var newWorker lib.Worker
	json.NewDecoder(r.Body).Decode(&newWorker)
	addWorker(klasterStatus, &newWorker)
}


type taskHelper struct {
	direction interface{}
	taskSize int
}

func (t *taskHelper) formNewTask() {
	select {
	case t.direction = <-workersPool:
		t.taskSize = t.direction.(*lib.Worker).Cores * WORKTASK_PER_CORE

	default:
		t.direction = taskPool
		t.taskSize = DEFAULT_WORKTASK_SIZE
	}
}

func (t *taskHelper) SetDirToTaskPool() {
	switch x := t.direction.(type) {
	case *lib.Worker:
		workersPool <- x
	}
	t.direction = taskPool
}

func(t *taskHelper) AddTask(clientdata map[string]lib.Matrix, reqMatrixes map[string]bool, root lib.ASTNode, newListName string) {
	switch x := t.direction.(type) {
	case *lib.Worker:
		newTask := lib.KlasterToWorkerTask{
			Root: root,
			Data: map[string]lib.Matrix{},
		}

		for k := range reqMatrixes {
			newTask.Data[k] = clientdata[k]
		}

		reqJson, _ := json.Marshal(newTask)
		go func(){
			response, _ := http.Post(fmt.Sprintf("http://localhost:%s/rofl", x.Port), "application/json", bytes.NewReader(reqJson))
			defer response.Body.Close()
			var ans lib.Matrix
			json.NewDecoder(response.Body).Decode(&ans)
			clientdata[newListName] = ans
		}()

	case chan lib.KlasterToWorkerTask:

	}
}

func exprHandle(w http.ResponseWriter, r *http.Request) {
	var clientReq lib.ClientTask
	json.NewDecoder(r.Body).Decode(&clientReq)
	exprDisturbToWorkes(workersPool, deferTasksPool, clientReq)

	fset := token.NewFileSet()
	defAst, _ := parser.ParseExprFrom(fset, "", clientReq.Expr, 0)
	astTree := formMyAst(defAst, clientReq.Data)

	

	taskManager := taskHelper{}
	taskManager.formNewTask()

	var getLeafsNames func(node lib.ASTNode, namesSet map[string]bool)
	getLeafsNames = func(node lib.ASTNode, namesSet map[string]bool) {
		switch x := node.(type) {
		case *lib.BinaryOp:
			getLeafsNames(x.Left, namesSet)
			getLeafsNames(x.Right, namesSet)
		case *lib.MatrixLeaf:
			namesSet[x.MatrixName] = true
		}
	}

	var dfs func(node lib.ASTNode) (int)
	dfs = func(node lib.ASTNode) (int) {
		switch x := node.(type) {
		case *lib.BinaryOp:
			childsDistibuted := dfs(x.Left) + dfs(x.Right)
			
			SubTreeOp := x.SubTreeCountOperations - childsDistibuted
			if SubTreeOp > taskManager.taskSize {

				subTreeMatrixes := map[string]bool{}
				getLeafsNames(node, subTreeMatrixes)

				allMatrixReady := true
				for k := range subTreeMatrixes{
					_, k_old := clientReq.Data[k]
					allMatrixReady = allMatrixReady && k_old
				}

				if !allMatrixReady {
					taskManager.SetDirToTaskPool()
				}

				SubTree := node
				node = &lib.MatrixLeaf{
					MatrixName: fmt.Sprint(rand.Int()),
					Size: SubTree.GetMatrixSize(),
				}

				taskManager.AddTask(clientReq.Data, subTreeMatrixes, SubTree, node.(*lib.MatrixLeaf).MatrixName)
			}

			return childsDistibuted

		case *lib.MatrixLeaf:
			return 0
		}
		return 0
	}
}



func main() {
	klasterStatus, workersPool, deferTasksPool = klasterInit(os.Args)
	
	http.HandleFunc("/addworker", addWorkerHandle)
	http.HandleFunc("/solve", exprHandle)

	
	fmt.Println("Klaster started on port", klasterStatus.Port)
	http.ListenAndServe(fmt.Sprintf(":%s", klasterStatus.Port), nil)
}