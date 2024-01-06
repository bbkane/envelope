package main

import (
	"fmt"

	"go.bbkane.com/warg/command"
)

func hello(ctx command.Context) error {
	fmt.Println("Hello!!")
	return nil
}
