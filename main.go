package main

import (
	"fmt"
	"os"

	"github.com/dyoung522/esotools/modTools"
)

func main() {
	mods, errs := modTools.GetMods()
	_ = mods

	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Println(e)
		}
		os.Exit(1)
	}

	// fmt.Println(mods.Print())
	fmt.Println(mods.PrintAll())
}
