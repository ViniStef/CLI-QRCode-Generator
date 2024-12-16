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
	qrc.matrix = addBlackPixel(qrc.matrix)
	qrc.matrix = addTimingStrips(qrc.matrix)
	qrc.addData()

	for i := 0; i < len(qrc.matrix); i++ {
		fmt.Println(qrc.matrix[i])
	}
}

func addBlackPixel(matrix [][]int) [][]int {
	matrix[13][8] = 2

	return matrix
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

func addTimingStrips(matrix [][]int) [][]int {
	for i := 8; i < len(matrix)-8; i++ {
		if i%2 == 0 {
			matrix[i][6] = 2
			matrix[6][i] = 2
		} else {
			matrix[i][6] = 1
			matrix[6][i] = 1
		}
	}

	return matrix
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
				if i == 6 {
					continue
				}
				for j := startCol; j > startCol-2 && j >= 0; j-- {
					if j == 6 {
						continue
					}
					if i == 13 && j == 8 {
						continue
					}
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
				if i == 6 {
					continue
				}

				for j := startCol - 1; j <= startCol && j >= 0; j++ {
					if j == 6 {
						continue
					}
					if i == 13 && j == 8 {
						continue
					}
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
