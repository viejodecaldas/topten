package main

import (
	"fmt"
	"github.com/viejodecaldas/topten/files"
	"os"
)

func main() {
	input := files.InputParams()

	err := files.ProcessOrders(input.ActorsFile, input.EventsFile, input.ReposFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "something went wrong while processing files. Error: %s", err.Error())
	}
}
