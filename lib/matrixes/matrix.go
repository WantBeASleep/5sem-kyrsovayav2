package matrixes

type Matrix struct {
	Size MatrixSize
	Grid [][]int
}

type MatrixSize struct {
	Rows int
	Cols int
}

func (m Matrix) Copy() Matrix {
	newM := Matrix{
		Size: MatrixSize{
			Rows: m.Size.Rows,
			Cols: m.Size.Cols,
		},
		Grid: make([][]int, m.Size.Rows),
	}

	for i := range newM.Grid {
		newM.Grid[i] = make([]int, m.Size.Cols)
		copy(newM.Grid[i], m.Grid[i])
	}

	return newM
}
