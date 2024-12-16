package renderer

import "strings"

func RenderQR(matrix [][]int) []string {
	renderMatrix := make([][]string, len(matrix))
	rowSlice := make([]string, 21)

	for i := range matrix {
		renderMatrix[i] = make([]string, len(matrix[i]))
		for j := range matrix[i] {
			switch matrix[i][j] {
			case 0, 3:
				renderMatrix[i][j] = "\u001B[37m██\u001B[0m"
			case 1, 2:
				renderMatrix[i][j] = "\u001B[30m██\u001B[0m"
			default:
				renderMatrix[i][j] = "  "
			}
		}
	}

	for _, row := range renderMatrix {
		var b strings.Builder
		for r := range row {
			b.WriteString(row[r])
		}
		rowSlice = append(rowSlice, b.String())
	}

	return rowSlice

}
