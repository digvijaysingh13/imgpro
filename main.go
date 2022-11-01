package main

import (
	"fmt"

	"github.com/digvijaysingh13/imgpro/util"
)

func main() {
	fmt.Println("hello go in main.")

	b := []byte{02, 04, 06, 8}

	fmt.Println(util.BytesToUnsignInt(&b))
	fmt.Println(0x08060402)
}
