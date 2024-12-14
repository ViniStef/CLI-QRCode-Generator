package main

import "QR-Code-CLI/internal/generator"

func main() {
	//renderer.RenderQR()
	you := generator.QRCodeV1{}
	you.InitializeMatrix()
}
