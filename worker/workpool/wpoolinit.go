package workpool

import (
	structs "lib/matrixes"
	tree "lib/trees"
)

var taskQ chan task

func worker(id int) {
	for {
		newTask := <-taskQ
		if !newTask.CanStartTask() {
			taskQ <- newTask
		} else {
			newTask.Start()
		}
	}
}

const CntOfWorkers int = 4

func StartPool() {
	taskQ = make(chan task, 1000000)
	for i := range CntOfWorkers {
		go worker(i)
	}
}

func Solve(root tree.ASTNode, data map[string]structs.Matrix) structs.Matrix {
	return <-ParseTreeToTasks(taskQ, root, data)
}