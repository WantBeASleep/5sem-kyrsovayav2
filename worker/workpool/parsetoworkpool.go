package workpool

import(
	ts "lib/tasks"
	tr "lib/trees"
	mt "lib/matrixes"
)

func ParseTreeToTasks(taskQ chan ts.Wpooltask, root tr.ASTNode, matrixesData map[string]mt.Matrix) chan mt.Matrix {
	matrixDisplay := func(name string) mt.Matrix {
		return matrixesData[name]
	}

	var dfs func(node tr.ASTNode, responseChannel chan mt.Matrix)
	dfs = func(node tr.ASTNode, responseChannel chan mt.Matrix) {
		switch x := node.(type) {
		case *tr.BinaryOp:
			leftMatrix := make(chan mt.Matrix, 1)
			rightMatrix := make(chan mt.Matrix, 1)

			dfs(x.Left, leftMatrix)
			dfs(x.Right, rightMatrix)

			newTask := ts.WpoolbinTaskStart{
				Response:    responseChannel,
				LeftArg:     leftMatrix,
				RightArg:    rightMatrix,
				Op:          x.Op,
				TaskChannel: taskQ,
			}

			taskQ <- newTask

		case *tr.MatrixLeaf:
			newTask := ts.WpoolcopyTask{
				Response:        responseChannel,
				MatrixName:      x.MatrixName,
				GetMatrixByName: matrixDisplay,
			}

			taskQ <- newTask

		}
	}

	result := make(chan mt.Matrix, 1)
	dfs(root, result)
	return result
}