package main

import (
	"fmt"

	"github.com/digvijaysingh13/imgpro/bmp"
	"github.com/digvijaysingh13/imgpro/util"
)

func main() {
	filename := "/home/digvijaysingh/Downloads/tiger.bmp"
	data, err := util.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	img := bmp.NewBmp(data)
	createDiffRgb(&img)
}

func createDiffRgb(img *bmp.Bmp) {
	// save gray
	gray := img.Copy()
	bmp.MakeGray(gray)
	err := util.WriteFile("sample_red1.bmp", gray.Data)
	if err == nil {
		fmt.Println("gray file created")
	} else {
		fmt.Println(err)
	}

	//
}
