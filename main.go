package main

import (
	"QR-Code-CLI/internal/generator"
	"QR-Code-CLI/internal/renderer"
	"fmt"
)

func main() {
	you := generator.QRCodeV2{}
	matrix := you.InitializeMatrix("www.youtube.com/veritasium")
	rows := renderer.RenderQR(matrix)
	for i := range rows {
		fmt.Println(rows[i])
	}

}
