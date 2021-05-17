package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/pkg/errors"

	"github.com/jonaslu/ain/internal/assemble"
	"github.com/jonaslu/ain/internal/pkg/call"
	"github.com/jonaslu/ain/internal/pkg/disk"
)

func printInternalErrorAndExit(err error) {
	formattedError := fmt.Errorf("An error occurred: %v", err.Error())
	fmt.Fprintln(os.Stderr, formattedError.Error())
	os.Exit(1)
}

func main() {
	var leaveTmpFile, printCommand bool
	var envFile string

	flag.BoolVar(&leaveTmpFile, "l", false, "Leave any temp-files")
	flag.BoolVar(&printCommand, "p", false, "print command to the terminal (do not execute it")
	flag.StringVar(&envFile, "e", ".env", "Path to .env file")
	flag.Parse()

	if err := disk.ReadEnvFile(envFile, envFile != ".env"); err != nil {
		printInternalErrorAndExit(err)
	}

	localTemplateFileNames, err := disk.GetTemplateFilenames()
	if err != nil {
		printInternalErrorAndExit(err)
	}

	if len(localTemplateFileNames) == 0 {
		printInternalErrorAndExit(errors.New("Missing file name\nUsage ain <template.ain> or connect it to a pipe"))
	}

	// !! TODO !! Hook into SIGINT etc and cancel this context if hit
	// This needs to be set when running shells (from that point on)
	ctx := context.Background()

	callData, fatal, err := assemble.Assemble(ctx, localTemplateFileNames)
	if err != nil {
		printInternalErrorAndExit(err)
	}

	if fatal != "" {
		fmt.Println(fatal)
		os.Exit(1)
	}

	backendOutput, err := call.CallBackend(ctx, callData, leaveTmpFile, printCommand)
	if err != nil {
		fmt.Fprint(os.Stderr, err)

		var backendErr *call.BackedErr
		if errors.As(err, &backendErr) {
			os.Exit(backendErr.ExitCode)
		}
	}

	fmt.Fprint(os.Stdout, backendOutput)
}
