package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/minodisk/thriftfmt/formatter"
)

type options struct {
	display bool
	report  bool
	list    bool
	write   bool
}

func main() {
	if err := _main(); err != nil {
		_, e := os.Stderr.WriteString(fmt.Sprintf("%s", err))
		if e != nil {
			// do nothing
		}
		os.Exit(2)
	}
}

func _main() error {
	_, files := parseFlag()
	for _, file := range files {
		if err := formatter.Format(file, os.Stdout); err != nil {
			return err
		}
	}
	return nil
}

func parseFlag() (options, []string) {
	display := flag.Bool("d", false, "display diffs instead of rewriting files")
	report := flag.Bool("e", false, "report all errors (not just the first 10 on different lines)")
	list := flag.Bool("l", false, "list files whose formatting differs from thriftfmt's")
	write := flag.Bool("w", false, "write result to (source) file instead of stdout")
	flag.Parse()
	files := flag.Args()
	return options{
		*display,
		*report,
		*list,
		*write,
	}, files
}
