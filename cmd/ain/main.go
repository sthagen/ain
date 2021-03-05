package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jonaslu/ain/internal/pkg/call"
	"github.com/jonaslu/ain/internal/pkg/disk"
	"github.com/jonaslu/ain/internal/pkg/parse"
)

func printInternalErrorAndExit(err error) {
	formattedError := fmt.Errorf("An error occurred: %v", err.Error())
	fmt.Fprintln(os.Stderr, formattedError.Error())
	os.Exit(1)
}

func main() {
	gotPipe, err := disk.IsConnectedToPipe()
	if err != nil {
		printInternalErrorAndExit(err)
	}

	if !gotPipe {
		if len(os.Args) < 2 {
			printInternalErrorAndExit(errors.New("Missing file name\nUsage ain <template.ain> or connect it to a pipe"))
		}
	}

	template, err := disk.ReadTemplate()
	if err != nil {
		printInternalErrorAndExit(err)
	}

	// !! TODO !! Hook into SIGINT etc and cancel this context if hit
	ctx := context.Background()

	callData, fatals := parse.ParseTemplate(ctx, template)
	if len(fatals) > 0 {
		for _, fatal := range fatals {
			fmt.Println(fatal)
		}

		os.Exit(1)
	}

	backendOutput, err := call.CallBackend(ctx, callData, "httpie")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Fprint(os.Stdout, backendOutput)
}
