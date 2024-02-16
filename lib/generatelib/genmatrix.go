package generatelib

import (
	mt "lib/matrixes"
	"math/rand"
)

func GenerateRandMatrix(m, n, k int) mt.Matrix {
	newMatrix := mt.Matrix{
		Size: mt.MatrixSize{
			Rows: m,
			Cols: n,
		},
		Grid: make([][]int, m),
	}

	for i := range m {
		newMatrix.Grid[i] = make([]int, n)
		for j := range n {
			newMatrix.Grid[i][j] = rand.Intn(k)
		}
	}

	return newMatrix
}
