package tasks

import (
	"go/token"
	mt "lib/matrixes"
)

type Wpooltask interface {
	CanStartTask() bool
	Start()
}

type WpoolcopyTask struct {
	Response chan mt.Matrix

	MatrixName      string
	GetMatrixByName func(name string) mt.Matrix
}

func (t WpoolcopyTask) CanStartTask() bool {
	return true
}

func (t WpoolcopyTask) Start() {
	t.Response <- t.GetMatrixByName(t.MatrixName).Copy()
}

type WpoolbinTaskStart struct {
	Response chan mt.Matrix

	LeftArg     chan mt.Matrix
	RightArg    chan mt.Matrix
	Op          token.Token
	TaskChannel chan Wpooltask
}

func (t WpoolbinTaskStart) CanStartTask() bool {
	if len(t.LeftArg) == 1 && len(t.RightArg) == 1 {
		return true
	} else {
		return false
	}
}

func (t WpoolbinTaskStart) Start() {
	GetMatrixCalcTasks(t.TaskChannel, t.Response, <-t.LeftArg, <-t.RightArg, t.Op)
}

type WpoolbinTaskEnd struct {
	response chan mt.Matrix

	countSubOps      int
	completeOpStatus chan bool
	res              *mt.Matrix
}

func (t WpoolbinTaskEnd) CanStartTask() bool {
	if len(t.completeOpStatus) == t.countSubOps {
		return true
	} else {
		return false
	}
}

func (t WpoolbinTaskEnd) Start() {
	t.response <- *t.res
}

type WpoolmatrixOpTask struct {
	status chan bool

	x *mt.Matrix // left
	y *mt.Matrix // right
	z *mt.Matrix // res

	startRow int
	startCol int
	endRow   int
	endCol   int
	op       func(x, y *mt.Matrix, i, j int) int
}

func (t WpoolmatrixOpTask) CanStartTask() bool {
	return true
}

func (t WpoolmatrixOpTask) Start() {
	for i := t.startRow; i <= t.endRow; i++ {
		for j := t.startCol; j <= t.endCol; j++ {
			t.z.Grid[i][j] = t.op(t.x, t.y, i, j)
		}
	}
	t.status <- true
}
