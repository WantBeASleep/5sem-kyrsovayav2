package tasks

import (
	mt "lib/matrixes"
	"go/token"
)

func plusOp(x, y *mt.Matrix, i, j int) int {
	return x.Grid[i][j] + y.Grid[i][j]
}

func multOp(x, y *mt.Matrix, i, j int) int {
	res := 0
	for k := 0; k < x.Size.Cols; k++ {
		res += x.Grid[i][k] * y.Grid[k][j]
	}
	return res
}

const OpForGoroutine int = 1e4
const CoresQty int = 4
const MinToParallel int = 1e4

func GetMatrixCalcTasks(taskQ chan Wpooltask, opResponseChannel chan mt.Matrix, x, y mt.Matrix, op token.Token) {
	switch op {
	case token.ADD:

		//Create new matrix
		res := mt.Matrix{
			Size: mt.MatrixSize{
				Rows: x.Size.Rows,
				Cols: x.Size.Cols,
			},
			Grid: make([][]int, x.Size.Rows),
		}
		for i := range res.Grid {
			res.Grid[i] = make([]int, res.Size.Cols)
		}

		// PARALLELING
		cntOps := x.Size.Rows * x.Size.Cols
		gorOps := cntOps / CoresQty

		interval := gorOps
		if cntOps < MinToParallel {
			interval = cntOps
		} else if gorOps < MinToParallel {
			interval = MinToParallel
		}

		startRow, endRow := 0, 0
		startCol, endCol := 0, 0

		taskQty := cntOps / interval
		if cntOps%interval != 0 {
			taskQty++
		}

		OpsStatusChannelMrBeast := make(chan bool, taskQty)

		for startRow < x.Size.Rows && startCol < x.Size.Cols {
			endCol = (startCol + interval - 1) % x.Size.Cols
			endRow = startRow + (startCol+interval-1)/x.Size.Cols
			if endRow >= x.Size.Rows {
				endRow = x.Size.Rows - 1
				endCol = x.Size.Cols - 1
			}

			newTask := WpoolmatrixOpTask{
				status:   OpsStatusChannelMrBeast,
				x:        &x,
				y:        &y,
				z:        &res,
				startRow: startRow,
				startCol: startCol,
				endRow:   endRow,
				endCol:   endCol,
				op:       plusOp,
			}
			taskQ <- newTask

			startCol = (endCol + 1) % x.Size.Cols
			startRow = endRow + (endCol+1)/x.Size.Cols
		}

		finalOpTask := WpoolbinTaskEnd{
			response:         opResponseChannel,
			countSubOps:      taskQty,
			completeOpStatus: OpsStatusChannelMrBeast,
			res:              &res,
		}
		taskQ <- finalOpTask

	case token.MUL:
		res := mt.Matrix{
			Size: mt.MatrixSize{
				Rows: x.Size.Rows,
				Cols: y.Size.Cols,
			},
			Grid: make([][]int, x.Size.Rows),
		}
		for i := range res.Grid {
			res.Grid[i] = make([]int, res.Size.Cols)
		}

		qtyInX := x.Size.Rows * x.Size.Cols
		qtyInY := y.Size.Rows * y.Size.Cols

		gorOpsX := qtyInX / CoresQty
		gorOpsY := qtyInY / CoresQty

		rowsInterval := x.Size.Rows
		colsInterval := y.Size.Cols
		
		if qtyInX > MinToParallel {
			rowsInterval = gorOpsX / x.Size.Cols
			if gorOpsX % x.Size.Cols != 0 {
				rowsInterval++
			}
		}

		if qtyInY > MinToParallel {
			colsInterval = gorOpsY / y.Size.Rows 
			if gorOpsY % y.Size.Rows != 0 {
				colsInterval++
			}
		}

		taskRowsQty := x.Size.Rows / rowsInterval
		if qtyInX % rowsInterval != 0 {
			taskRowsQty++
		}
		taskColQty := y.Size.Cols / colsInterval
		if y.Size.Cols%colsInterval != 0 {
			taskColQty++
		}

		taskQty := taskRowsQty * taskColQty
		OpsStatusChannelMrBeast := make(chan bool, taskQty)

		for i := 0; i < x.Size.Rows; i += rowsInterval {
			for j := 0; j < y.Size.Cols; j += colsInterval {
				endRow := i + rowsInterval - 1
				if endRow > x.Size.Rows {
					endRow = x.Size.Rows - 1
				}
				endCol := j + colsInterval - 1
				if endCol > y.Size.Cols {
					endCol = y.Size.Cols - 1
				}

				newTask := WpoolmatrixOpTask{
					status:   OpsStatusChannelMrBeast,
					x:        &x,
					y:        &y,
					z:        &res,
					startRow: i,
					startCol: j,
					endRow:   endRow,
					endCol:   endCol,
					op:       multOp,
				}
				taskQ <- newTask
			}
		}

		finalOpTask := WpoolbinTaskEnd{
			response:         opResponseChannel,
			countSubOps:      taskQty,
			completeOpStatus: OpsStatusChannelMrBeast,
			res:              &res,
		}
		taskQ <- finalOpTask
	}

}