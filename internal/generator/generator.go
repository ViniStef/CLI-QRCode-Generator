package generator

import (
	"fmt"
	"strings"
)

type QRCode interface {
	InitializeMatrix()
	addPositionSquares()
	addIndicators()
	addData()
}

type QRCodeV1 struct {
	Matrix [][]int
}

func (qrc QRCodeV1) GetMatrix() [][]int {
	return qrc.Matrix
}

func (qrc QRCodeV1) InitializeMatrix(url string) [][]int {
	qrc.Matrix = make([][]int, 25)
	for i := 0; i < 25; i++ {
		qrc.Matrix[i] = make([]int, 25)
	}

	qrc.addPositionSquares()
	qrc.addIndicators(url)
	qrc.Matrix = addBlackPixel(qrc.Matrix)
	qrc.Matrix = addTimingStrips(qrc.Matrix)
	qrc.addData(url)

	return qrc.Matrix
}

func addBlackPixel(matrix [][]int) [][]int {
	matrix[len(matrix)-8][8] = 2

	return matrix
}

func (qrc QRCodeV1) addPositionSquares() {
	qrc.Matrix = createSquare(qrc.Matrix, 0, 0, "topLeft")
	qrc.Matrix = createSquare(qrc.Matrix, 0, 18, "topRight")
	qrc.Matrix = createSquare(qrc.Matrix, 18, 0, "botLeft")
}

func createSquare(matrix [][]int, rowStart, colStart int, position string) [][]int {
	square := matrix

	switch position {
	case "topLeft":
		for i := 0; i < 8; i++ {
			square[rowStart+7][i] = 3
			square[i][rowStart+7] = 3
		}
	case "topRight":
		for i := 0; i < 8; i++ {
			square[rowStart+7][i+colStart-1] = 3
			square[i][colStart-1] = 3
		}
	case "botLeft":
		for i := 0; i < 8; i++ {
			square[rowStart-1][i] = 3
			square[i+rowStart-1][colStart+7] = 3
		}
	}

	for i := rowStart; i < rowStart+7; i++ {
		for j := colStart; j < colStart+7; j++ {
			if i != rowStart && i != (rowStart+7-1) && j != colStart && j != (colStart+7-1) {
				if i >= rowStart+2 && i <= rowStart+4 && j >= colStart+2 && j <= colStart+4 {
					square[i][j] = 2
				} else {
					square[i][j] = 3
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
			matrix[i][6] = 3
			matrix[6][i] = 3
		}
	}

	return matrix
}

func (qrc QRCodeV1) addIndicators(url string) {
	matrix := qrc.Matrix
	matrix[len(matrix)-1][len(matrix)-1] = 0
	matrix[len(matrix)-1][len(matrix)-2] = 1
	matrix[len(matrix)-2][len(matrix)-1] = 0
	matrix[len(matrix)-2][len(matrix)-2] = 0

	qrc.Matrix = matrix

	binaryMsgSize := fmt.Sprintf("%08b", len(url))

	startBinaryRow := len(matrix) - 3
	finishBinaryRow := len(matrix) - 6
	countCurrBinary := 0

	for i := startBinaryRow; i >= finishBinaryRow; i-- {
		for j := len(matrix) - 1; j > len(matrix)-3; j-- {
			qrc.Matrix[i][j] = int(binaryMsgSize[countCurrBinary] - '0')
			countCurrBinary++
		}
	}
}

func stringToBinary(s string) string {
	var b strings.Builder

	for _, c := range s {
		b.WriteString(fmt.Sprintf("%08b", c))
	}

	fmt.Println("\x1b[32;1m QRCode linking to \x1b[0m", s)

	return b.String()
}

func (qrc QRCodeV1) addData(url string) {
	matrix := qrc.Matrix

	binaryUrl := stringToBinary(url)
	var currBinary int

	currIteration := 0
	startRow := len(matrix) - 1
	startCol := len(matrix) - 1

	goUpwards := true

	for currIteration < (len(matrix)/2) && currBinary < len(binaryUrl) {
		if goUpwards {
			for i := startRow; i >= 0; i-- {
				if i == 6 {
					continue
				}
				for j := startCol; j > startCol-2 && j >= 0; j-- {
					if currBinary == len(binaryUrl) {
						break
					}
					if j == 6 {
						continue
					}
					if i == 13 && j == 8 {
						continue
					}
					if i < len(matrix)-6 {
						if matrix[i][j] == 0 {
							matrix[i][j] = int(binaryUrl[currBinary] - '0')
							currBinary++
						}
					} else if startCol < len(matrix)-2 {
						if matrix[i][j] == 0 {
							matrix[i][j] = int(binaryUrl[currBinary] - '0')
							currBinary++
						}
					}
				}
			}

			goUpwards = false
		} else {
			for i := 0; i <= startRow; i++ {
				if i == 6 {
					continue
				}
				for j := startCol - 1; j <= startCol && j >= 0; j++ {
					if currBinary == len(binaryUrl) {
						break
					}
					if j == 6 {
						continue
					}
					if i == 13 && j == 8 {
						continue
					}
					if matrix[i][j] == 0 {
						matrix[i][j] = int(binaryUrl[currBinary] - '0')
						currBinary++
					}
				}
			}

			goUpwards = true
		}

		startCol = startCol - 2
		currIteration++
	}
}
