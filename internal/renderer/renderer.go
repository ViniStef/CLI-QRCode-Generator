package renderer

import "fmt"

func RenderQR() {
	black := "  "

	board := [21][21]string{}

	for i := 0; i < 21; i++ {
		for j := 0; j < 21; j++ {
			board[i][j] = black
		}
	}

	for i := 0; i < len(board); i++ {
		fmt.Println(board[i])
	}
}
