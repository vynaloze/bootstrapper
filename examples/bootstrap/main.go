package main

import (
	"bootstrapper/blueprint"
	"fmt"
)

func main() {
	opts := blueprint.BootstrapOpts{}

	err := blueprint.Bootstrap(opts)
	if err != nil {
		fmt.Println(err)
	}
}
