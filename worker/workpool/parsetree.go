package workpool

import (
	structs "lib/matrixes"
	tree "lib/trees"
	
)

func ParseTreeToTasks(taskQ chan task, root tree.ASTNode, matrixesData map[string]structs.Matrix) chan structs.Matrix {
	matrixDisplay := func(name string) structs.Matrix {
		return matrixesData[name]
	}

	var dfs func(node tree.ASTNode, responseChannel chan structs.Matrix)
	dfs = func(node tree.ASTNode, responseChannel chan structs.Matrix) {
		switch x := node.(type) {
		case *tree.BinaryOp:
			leftMatrix := make(chan structs.Matrix, 1)
			rightMatrix := make(chan structs.Matrix, 1)

			dfs(x.Left, leftMatrix)
			dfs(x.Right, rightMatrix)

			newTask := binTaskStart{
				response:    responseChannel,
				leftArg:     leftMatrix,
				rightArg:    rightMatrix,
				op:          x.Op,
				taskChannel: taskQ,
			}

			taskQ <- newTask

		case *tree.MatrixLeaf:
			newTask := copyTask{
				response:        responseChannel,
				matrixName:      x.MatrixName,
				getMatrixByName: matrixDisplay,
			}

			taskQ <- newTask

		}
	}

	result := make(chan structs.Matrix, 1)
	dfs(root, result)
	return result
}