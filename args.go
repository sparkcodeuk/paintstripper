package main

import (
	"flag"
	"fmt"
	"os"
)

type Args struct {
	Help            bool
	Version         bool
	Quiet           bool
	Force           bool
	Color           bool
	WriteStripped   string
	WriteUnstripped string
	InputPath       string
}

func fatalError(message string) {
	fmt.Println("Error:", message)
	fmt.Printf("\n---\n\n")
	flag.Usage()
	os.Exit(1)
}

func ParseArgs() Args {
	args := Args{}

	flag.BoolVar(&args.Help, "help", false, "Display this help message")
	flag.BoolVar(&args.Version, "version", false, "Print version and exit")
	flag.BoolVar(&args.Quiet, "quiet", false, "Print no output")
	flag.BoolVar(&args.Force, "force", false, "Force file overwrite")
	flag.BoolVar(&args.Color, "color", false, "Strip shell color codes")
	flag.StringVar(&args.WriteStripped, "write-stripped", "", "Write stripped content to a file (print out unstripped content)")
	flag.StringVar(&args.WriteUnstripped, "write-unstripped", "", "Write unstripped content to a file (print out stripped content)")

	// Custom usage message
	flag.Usage = func() {
		PrintVersion()
		fmt.Println(`
Usage:-
  paintstripper <args> [input file]`)
		fmt.Println()
		flag.PrintDefaults()
		fmt.Println(`
Command examples:-

Strip color codes from a given input file and print output to terminal:
  paintstripper -color <input file>

... the same but using the pipe operator instead of an input file as the source:
  some_colorful_command.sh | paintstripper -color

Strip colors & print and also write a copy of the original, unstripped content to a file:
  some_colorful_command.sh | paintstripper -color -write-unstripped colorful_output.log

... the same but printing the colorful output and writing the stripped output to a file:
  some_colorful_command.sh | paintstripper -color -write-stripped stripped_output.log

Print out no output and write the stripped output to a file:
  some_colorful_command.sh | paintstripper -color -quiet -write-stripped stripped_output.log

... the same (print no output), but write both the stripped & unstripped output to separate files:
  some_colorful_command.sh | paintstripper -color -quiet -write-stripped stripped_output.log -write-unstripped colorful_output.log`)
		fmt.Println()
	}

	flag.Parse()

	if args.Help {
		flag.Usage()
		os.Exit(0)
	}

	if args.Version {
		PrintVersion()
		os.Exit(0)
	}

	// Parse positional arguments
	switch flag.NArg() {
	case 0:
		// Will read from stdin
	case 1:
		args.InputPath = flag.Arg(0)
	default:
		fatalError("this command only takes a maximum of one input file")
	}

	// Parse any conflicting argument combinations
	if !args.Color {
		// (other stripping options will be available later; one processing option must be selected)
		fatalError("you must specify at least one processing option (e.g., --color)")
	}

	return args
}
