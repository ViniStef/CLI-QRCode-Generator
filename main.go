package main

import (
	"QR-Code-CLI/internal/generator"
	"QR-Code-CLI/internal/renderer"
	"fmt"
)

func main() {
	you := generator.QRCodeV1{}
	matrix := you.InitializeMatrix("www.google.com")
	rows := renderer.RenderQR(matrix)
	for i := range rows {
		fmt.Println(rows[i])
	}

}
