package main

import (
	"fmt"
	"os"

	"github.com/dyoung522/esotools/modTools"
)

func main() {
	modlist, err := modTools.FindMods()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	mods, errs := modTools.ReadMods(&modlist)
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
