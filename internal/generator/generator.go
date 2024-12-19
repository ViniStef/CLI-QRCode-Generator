package generator

import (
	"fmt"
	"github.com/klauspost/reedsolomon"
	"log"
	"math"
	"strconv"
	"strings"
)

type QRCode interface {
	InitializeMatrix()
	addPositionSquares()
	addIndicators()
	addData()
}

type QRCodeV2 struct {
	Matrix [][]int
}

func (qrc QRCodeV2) GetMatrix() [][]int {
	return qrc.Matrix
}

func (qrc QRCodeV2) InitializeMatrix(url string) [][]int {
	qrc.Matrix = make([][]int, 25)
	for i := 0; i < 25; i++ {
		qrc.Matrix[i] = make([]int, 25)
	}

	qrc.addPositionSquares()
	qrc.Matrix = addAlignmentPattern(qrc.Matrix)
	qrc.Matrix = addFormatStrips(qrc.Matrix)
	qrc.addIndicators(url)
	qrc.Matrix = addBlackPixel(qrc.Matrix)
	qrc.Matrix = addTimingStrips(qrc.Matrix)
	qrc.addData(url)

	for _, val := range qrc.Matrix {
		fmt.Println(val)
	}

	return qrc.Matrix
}

func addBlackPixel(matrix [][]int) [][]int {
	matrix[len(matrix)-8][8] = 2

	return matrix
}

func getAlignmentCoordinates(version int) []int {
	if version <= 1 {
		return nil
	}

	intervals := (version / 7) + 1

	distance := 4*version + 4

	step := int(math.Round(float64(distance) / float64(intervals)))

	if step%2 != 0 {
		step++
	}

	coordinates := make([]int, intervals+1)

	coordinates[0] = 6

	for i := 1; i <= intervals; i++ {
		coordinates[i] = 6 + distance - step*(intervals-i)
	}

	return coordinates
}

func addAlignmentPattern(matrix [][]int) [][]int {
	coordinates := getAlignmentCoordinates(2)
	locAlignmentPatternCenter := coordinates[1]

	for i := locAlignmentPatternCenter - 2; i <= locAlignmentPatternCenter+2; i++ {
		if i == locAlignmentPatternCenter-2 || i == locAlignmentPatternCenter+2 {
			for j := locAlignmentPatternCenter - 2; j <= locAlignmentPatternCenter+2; j++ {
				matrix[i][j] = 2
				matrix[j][i] = 2
			}
		} else {
			for j := locAlignmentPatternCenter - 1; j <= locAlignmentPatternCenter+1; j++ {
				matrix[i][j] = 3
				matrix[j][i] = 3
			}
		}
	}

	matrix[locAlignmentPatternCenter][locAlignmentPatternCenter] = 2

	return matrix
}

func addFormatStrips(matrix [][]int) [][]int {
	for i := 0; i <= 8; i++ {
		if i != 8 {
			matrix[8][len(matrix)-i-1] = 5
		}
		if i == 6 {
			continue
		}
		matrix[8][i] = 5
		matrix[i][8] = 5
	}

	for j := 0; j < 8; j++ {
		matrix[len(matrix)-8+j][8] = 5
		if j == 6 {
			continue
		}
	}

	return matrix

}

func (qrc QRCodeV2) addPositionSquares() {
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

func (qrc QRCodeV2) addIndicators(url string) {
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

	// terminator
	b.WriteString("0000")

	fmt.Println("\x1b[32;1m QRCode linking to \x1b[0m", s)

	padByte236 := "11101100"
	padByte17 := "00010001"
	padBytes := []string{padByte236, padByte17}
	currPad := 0

	currByteAmount := len(b.String()) / 8

	if currByteAmount < 26 {
		for i := currByteAmount; i <= 26; i++ {
			b.WriteString(padBytes[currPad])
			if currPad == 0 {
				currPad = 1
			} else {
				currPad = 0
			}
		}
	} else if currByteAmount > 26 {
		panic("Cannot create a version 2 QRCode with a message bigger than 26 bytes")
	}

	fmt.Println("B string aqui: ", b.String(), "Tamanho da string: ", len(b.String()))

	return b.String()
}

func binaryToBytes(codificationMode, binaryMsg string) []byte {
	var b strings.Builder
	binaryMsgSize := fmt.Sprintf("%08b", len(binaryMsg)/8)

	b.WriteString(codificationMode)
	b.WriteString(binaryMsgSize)
	b.WriteString(binaryMsg)

	s := b.String()

	var msgBytes []byte

	if len(s) == 28*8 {
		for i := 0; i < len(s); i += 8 {
			bitGroup := s[i : i+8]

			value, err := strconv.ParseUint(bitGroup, 2, 8)
			if err != nil {
				panic("Could not convert to base 2 8 bit representation")
			}

			msgBytes = append(msgBytes, byte(value))
		}
	}

	return msgBytes
}

func (qrc QRCodeV2) addData(url string) {
	matrix := qrc.Matrix

	binaryUrl := stringToBinary(url)
	var currBinary int

	msgBytes := binaryToBytes("0100", binaryUrl)

	dataShards := 28
	parityShards := 16
	totalShards := dataShards + parityShards

	encoder, err := reedsolomon.New(dataShards, parityShards)

	if err != nil {
		panic("Could not create reed solomon encoder")
	}

	shards := make([][]byte, totalShards)
	for i := 0; i < dataShards; i++ {
		shards[i] = []byte{msgBytes[i]}
	}
	for i := dataShards; i < totalShards; i++ {
		shards[i] = make([]byte, 1)
	}

	err = encoder.Encode(shards)
	if err != nil {
		log.Fatalf("Error while codifying shards: %v", err)
	}

	_, err = encoder.Verify(shards)
	if err != nil {
		panic("Data shards aren't of equal size")
	}

	var b strings.Builder

	for i := 0; i < dataShards; i++ {
		for j := range shards[i] {
			b.WriteString(fmt.Sprintf("%08b", shards[i][j]))
		}
	}

	for i := dataShards; i < totalShards; i++ {
		for j := range shards[i] {
			b.WriteString(fmt.Sprintf("%08b", shards[i][j]))
		}
	}

	// Skip the 4 first from the data mode and the next 8 from the msg length
	urlAndCorrectionBinary := b.String()[12:]

	currIteration := 0
	startRow := len(matrix) - 1
	startCol := len(matrix) - 1

	goUpwards := true

	for currIteration < (len(matrix)/2) && currBinary < len(urlAndCorrectionBinary) {
		if goUpwards {
			for i := startRow; i >= 0; i-- {
				if i == 6 {
					continue
				}
				for j := startCol; j > startCol-2 && j >= 0; j-- {
					if currBinary == len(urlAndCorrectionBinary) {
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
							matrix[i][j] = int(urlAndCorrectionBinary[currBinary] - '0')
							currBinary++
						}
					} else if startCol < len(matrix)-2 {
						if matrix[i][j] == 0 {
							matrix[i][j] = int(urlAndCorrectionBinary[currBinary] - '0')
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
				for j := startCol; j >= startCol-1 && j >= 0; j-- {
					if currBinary == len(urlAndCorrectionBinary) {
						break
					}
					if j == 6 {
						continue
					}
					if i == 13 && j == 8 {
						continue
					}
					if matrix[i][j] == 0 {
						matrix[i][j] = int(urlAndCorrectionBinary[currBinary] - '0')
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
