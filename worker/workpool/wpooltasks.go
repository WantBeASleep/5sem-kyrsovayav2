package workpool

import (
	structs "lib/matrixes"
	"go/token"
)

type task interface {
	CanStartTask() bool
	Start()
}

type copyTask struct {
	response chan structs.Matrix

	matrixName      string
	getMatrixByName func(name string) structs.Matrix
}

func (t copyTask) CanStartTask() bool {
	return true
}

func (t copyTask) Start() {
	t.response <- t.getMatrixByName(t.matrixName).Copy()
}

type binTaskStart struct {
	response chan structs.Matrix

	leftArg     chan structs.Matrix
	rightArg    chan structs.Matrix
	op          token.Token
	taskChannel chan task
}

func (t binTaskStart) CanStartTask() bool {
	if len(t.leftArg) == 1 && len(t.rightArg) == 1 {
		return true
	} else {
		return false
	}
}

func (t binTaskStart) Start() {
	GetMatrixCalcTasks(t.taskChannel, t.response, <-t.leftArg, <-t.rightArg, t.op)
}

type binTaskEnd struct {
	response chan structs.Matrix

	countSubOps      int
	completeOpStatus chan bool
	res              *structs.Matrix
}

func (t binTaskEnd) CanStartTask() bool {
	if len(t.completeOpStatus) == t.countSubOps {
		return true
	} else {
		return false
	}
}

func (t binTaskEnd) Start() {
	t.response <- *t.res
}

type matrixOpTask struct {
	status chan bool

	x *structs.Matrix // left
	y *structs.Matrix // right
	z *structs.Matrix // res

	startRow int
	startCol int
	endRow   int
	endCol   int
	op       func(x, y *structs.Matrix, i, j int) int
}

func (t matrixOpTask) CanStartTask() bool {
	return true
}

func (t matrixOpTask) Start() {
	for i := t.startRow; i <= t.endRow; i++ {
		for j := t.startCol; j <= t.endCol; j++ {
			t.z.Grid[i][j] = t.op(t.x, t.y, i, j)
		}
	}
	t.status <- true
}