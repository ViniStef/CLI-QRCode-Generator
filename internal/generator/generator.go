package generator

import "fmt"

type QRCode interface {
	initializeMatrix()
	addPositionSquares()
	addIndicators()
	addData()
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
	qrc.addIndicators("www.google.com")
	qrc.addData()

	for i := 0; i < len(qrc.matrix); i++ {
		fmt.Println(qrc.matrix[i])
	}
}

func (qrc QRCodeV1) addPositionSquares() {
	qrc.matrix = createSquare(qrc.matrix, 0, 0, "topLeft")
	qrc.matrix = createSquare(qrc.matrix, 0, 14, "topRight")
	qrc.matrix = createSquare(qrc.matrix, 14, 0, "botLeft")
}

func createSquare(matrix [][]int, rowStart, colStart int, position string) [][]int {
	square := matrix

	switch position {
	case "topLeft":
		for i := 0; i < 8; i++ {
			square[rowStart+7][i] = 1
			square[i][rowStart+7] = 1
		}
	case "topRight":
		for i := 0; i < 8; i++ {
			square[rowStart+7][i+colStart-1] = 1
			square[i][colStart-1] = 1
		}
	case "botLeft":
		for i := 0; i < 8; i++ {
			square[rowStart-1][i] = 1
			square[i+rowStart-1][colStart+7] = 1
		}
	}

	for i := rowStart; i < rowStart+7; i++ {
		for j := colStart; j < colStart+7; j++ {
			if i != rowStart && i != (rowStart+7-1) && j != colStart && j != (colStart+7-1) {
				if i >= rowStart+2 && i <= rowStart+4 && j >= colStart+2 && j <= colStart+4 {
					if i == rowStart+3 && j == colStart+3 {
						square[i][j] = 1
					} else {
						square[i][j] = 2
					}
				} else {
					square[i][j] = 1
				}
			} else {
				square[i][j] = 2
			}
		}
	}

	return square
}

func (qrc QRCodeV1) addIndicators(msg string) {
	matrix := qrc.matrix
	matrix[len(matrix)-1][len(matrix)-1] = 0
	matrix[len(matrix)-1][len(matrix)-2] = 1
	matrix[len(matrix)-2][len(matrix)-1] = 0
	matrix[len(matrix)-2][len(matrix)-2] = 0

	qrc.matrix = matrix

	binaryMsgSize := fmt.Sprintf("%08b", len(msg))

	startBinaryRow := len(matrix) - 3
	finishBinaryRow := len(matrix) - 6
	countCurrBinary := 0

	for i := startBinaryRow; i >= finishBinaryRow; i-- {
		for j := len(matrix) - 1; j > len(matrix)-3; j-- {
			qrc.matrix[i][j] = int(binaryMsgSize[countCurrBinary] - '0')
			countCurrBinary++
		}
	}
}

func (qrc QRCodeV1) addData() {
	matrix := qrc.matrix

	currIteration := 0
	startRow := len(matrix) - 1
	startCol := len(matrix) - 1

	goUpwards := true

	for currIteration < (len(matrix) / 2) {
		if goUpwards {
			for i := startRow; i >= 0; i-- {
				for j := startCol; j > startCol-2 && j >= 0; j-- {
					if i < len(matrix)-6 {
						if matrix[i][j] == 0 {
							matrix[i][j] = 3
						}
					} else if startCol < len(matrix)-2 {
						if matrix[i][j] == 0 {
							matrix[i][j] = 3
						}
					}
				}
			}

			goUpwards = false
		} else {
			testVar := 4
			for i := 0; i <= startRow; i++ {
				for j := startCol - 1; j <= startCol && j >= 0; j++ {
					if matrix[i][j] == 0 {
						matrix[i][j] = testVar
						if testVar == 4 {
							testVar = 5
						} else {
							testVar = 4
						}
					}
				}
			}

			goUpwards = true
		}

		startCol = startCol - 2

		currIteration++
	}
}
