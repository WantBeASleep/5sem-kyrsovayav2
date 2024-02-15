package main

import (
	"lib"
	"sync/atomic"
)

func exprDisturbToWorkes(workersPool chan *lib.Worker, deferTasksPool chan lib.DeferKlasterToWorkerTask, clientTask lib.ClientTask) {
	Tree := lib.ParseExpr(clientTask.Expr)
	Matrixes := clientTask.Data
	lib.UpdateTreeStats(Tree, Matrixes)

	ReadyMatrixes := map[string]*atomic.Bool{}
	for k := range Matrixes {
		ReadyMatrixes[k].Store(true)
	}

	Sender := NewTaskSender(workersPool, deferTasksPool)
	Sender.NewTask()

	var dfs func(node lib.ASTNode, prevNode lib.ASTNode)
	dfs = func(node lib.ASTNode, prevNode lib.ASTNode) {
		if x, ok := node.(*lib.BinaryOp); ok {
			dfs(x.Left, node)
			dfs(x.Right, node)

			if x.GetCountOp() >= Sender.OpCount {
				ReplaceLeaf := &lib.MatrixLeaf{
					MatrixName: lib.GetRandString(100),
					Size: x.GetMatrixSize(),
				}

				if prevNode.(*lib.BinaryOp).Left == x {
					prevNode.(*lib.BinaryOp).Left = ReplaceLeaf
				} else {
					prevNode.(*lib.BinaryOp).Right = ReplaceLeaf
				}

				ReadyMatrixes[ReplaceLeaf.MatrixName].Store(false)

				
			}

		}
	}
}