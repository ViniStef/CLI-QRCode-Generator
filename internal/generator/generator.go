package generator

import "fmt"

type QRCode interface {
	initializeMatrix()
	addPositionSquares()
	addIndicators()
	render() string
}

type QRCodeV1 struct {
	matrix [][]int
}

func (qrc QRCodeV1) InitializeMatrix() {
	qrc.matrix = make([][]int, 21)
	for i := 0; i < 21; i++ {
		qrc.matrix[i] = make([]int, 21)
	}

	qrc.addPositionSquares()

	for i := 0; i < len(qrc.matrix); i++ {
		fmt.Println(qrc.matrix[i])
	}
}

func (qrc QRCodeV1) addPositionSquares() {
	qrc.matrix = createSquare(qrc.matrix, 0, 0)
	qrc.matrix = createSquare(qrc.matrix, 0, 14)
	qrc.matrix = createSquare(qrc.matrix, 14, 0)
}

func createSquare(matrix [][]int, rowStart, colStart int) [][]int {
	bottomSquare := matrix

	for i := rowStart; i < rowStart+7; i++ {
		for j := colStart; j < colStart+7; j++ {
			if i != rowStart && i != (rowStart+7-1) && j != colStart && j != (colStart+7-1) {
				if i >= rowStart+2 && i <= rowStart+4 && j >= colStart+2 && j <= colStart+4 {
					if i == rowStart+3 && j == colStart+3 {
						bottomSquare[i][j] = 0
					} else {
						bottomSquare[i][j] = 1
					}
				} else {
					bottomSquare[i][j] = 0
				}
			} else {
				bottomSquare[i][j] = 1
			}
		}
	}

	return bottomSquare
}
