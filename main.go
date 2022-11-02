package main

import (
	"fmt"

	"github.com/digvijaysingh13/imgpro/bmp"
	"github.com/digvijaysingh13/imgpro/util"
)

func main() {
	filename := "/home/digvijaysingh/Desktop/sample_file.bmp"
	data, err := util.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	img := bmp.NewBmp(data)
	createDiffRgb(&img)
	img.PrintHeader()
}

func createDiffRgb(img *bmp.Bmp) {
	// save red
	red := img.Copy()
	bmp.MakeRed(red)
	err := util.WriteFile("sample_red1.bmp", red.Data)
	if err == nil {
		println("red file created")
	} else {
		println(err)
	}
	// save green
	green := img.Copy()
	bmp.MakeGreen(green)
	err = util.WriteFile("sample_green1.bmp", green.Data)
	if err == nil {
		println("green file created")
	} else {
		println(err)
	}
	// save blue
	blue := img.Copy()
	bmp.MakeBlue(blue)
	err = util.WriteFile("sample_blue1.bmp", blue.Data)
	if err == nil {
		println("blue file created")
	} else {
		println(err)
	}
}
