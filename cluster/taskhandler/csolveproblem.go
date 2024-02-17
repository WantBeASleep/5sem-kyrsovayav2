package taskhandler

import (
	"encoding/json"
	"fmt"
	"io"
	gn "lib/generatelib"
	is "lib/infostructs"
	mt "lib/matrixes"
	rq "lib/requests"
	ts "lib/tasks"
	tr "lib/trees"
)

func Handler_csolveproblem(workersPool chan *is.WorkerInfo, deferClusterWorkerTaskPool chan ts.ClusterWorkerTask, reqData io.ReadCloser) []byte {
	var inputData rq.ClientReq
	err := json.NewDecoder(reqData).Decode(&inputData)
	if err != nil {
		panic("Ошибка парса задачи на кластере! csolveproblem")
	}

	matrixes := inputData.Matrixes
	tree := tr.ParseExpr(inputData.Expr)
	tr.UpdateTreeStats(tree, matrixes)
	fmt.Println(tree.(*tr.BinaryOp))

	matrixesAlertReady := map[string]chan bool{}
	for k := range matrixes {
		matrixesAlertReady[k] = make(chan bool, 1)
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

				sender.Send(x, matrixes, &matrixesAlertReady, lastName)
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