package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dyoung522/esotools/modTools"
)

func main() {
	mods, err := modTools.FindMods()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("%v\n", strings.Join(mods, "\n"))
}
