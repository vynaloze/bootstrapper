package main

import (
	"bootstrapper/blueprint"
	"fmt"
)

func main() {
	err := blueprint.Bootstrap()
	if err != nil {
		fmt.Println(err)
	}
}
