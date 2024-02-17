package taskhandler

import (
	is "lib/infostructs"
	"encoding/json"
	rq "lib/requests"
	ts "lib/tasks"
	tr "lib/trees"
	mt "lib/matrixes"
	gn "lib/generatelib"
	"io"
)

func Handler_csolveproblem(workersPool chan *is.WorkerInfo, deferClusterWorkerTaskPool chan ts.ClusterWorkerTask, reqData io.ReadCloser) []byte {
	var inputData rq.ClientReq
	err := json.NewDecoder(reqData).Decode(&inputData)
	if err != nil {
		panic("Ошибка парса задачи на кластере! csolveproblem")
	}

	tree := tr.ParseExpr(inputData.Expr)
	matrixes := inputData.Matrixes
	matrixesAlertReady := map[string]chan bool{}
	for k := range matrixes {
		matrixesAlertReady[k] = make(chan bool, 1)
		matrixesAlertReady[k] <- true
		close(matrixesAlertReady[k])
	}

	sender := NewTaskSender(workersPool, deferClusterWorkerTaskPool)
	sender.NewTask()

	lastName := ""
	var dfs func(node tr.ASTNode, prevnode tr.ASTNode)
	dfs = func(node, prevnode tr.ASTNode) {
		if x, ok := node.(*tr.BinaryOp); ok {
			dfs(x.Left, node)
			dfs(x.Right, node)

			if (x.GetCountOp() >= sender.OpCount ) || (x == tree.(*tr.BinaryOp)) {
				newLeaf := &tr.MatrixLeaf{
					MatrixName: gn.GetRandString(100),
					Size: x.GetMatrixSize(),
				}
				lastName = newLeaf.MatrixName
				matrixesAlertReady[lastName] = make(chan bool, 1)
				matrixes[lastName] = mt.Matrix{
					Size: x.GetMatrixSize(),
				}

				if x != tree.(*tr.BinaryOp) {
					if prevnode.(*tr.BinaryOp).Left == x {
						prevnode.(*tr.BinaryOp).Left = newLeaf
					} else {
						prevnode.(*tr.BinaryOp).Right = newLeaf
					}
					tr.UpdateTreeStats(tree, matrixes)
				}

				sender.Send(tree, matrixes, &matrixesAlertReady, lastName)
				sender.NewTask()
			}
		}
	}
	dfs(tree, tree)
	sender.WorkerDrop()

	<-matrixesAlertReady[lastName]
	
	result, err := json.Marshal(matrixes[lastName])
	if err != nil {
		panic("Ошибка парса результата поддерева на кластере")
	}
	return result
}