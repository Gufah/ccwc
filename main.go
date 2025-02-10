package main

import (
	"ccwc/internal/wc"
	"fmt"
	"os"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	opts, filenames := wc.ParseArgs(os.Args[1:])
	results, total := processInputs(opts, filenames)
	fmt.Println(results, total)
	printResults(results, total, opts, filenames)
}

func processInputs(opts wc.Options, filenames []string) ([]wc.Counts, wc.Counts) {
	var results []wc.Counts
	var total wc.Counts

	if len(filenames) == 0 {
		counts, _ := wc.Count(os.Stdin)
		results = append(results, counts)
		total = counts
	} else {
		for _, filename := range filenames {
			counts, err := wc.ProcessFile(filename)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "wc: %s: %v\n", filename, err)
				continue

			}
			results = append(results, counts)
			total.Lines += counts.Lines
			total.Chars += counts.Chars
			total.Bytes += counts.Bytes
			total.Words += counts.Words
		}
	}
	return results, total
}

func printResults(results []wc.Counts, total wc.Counts, opts wc.Options, filenames []string) {
	for i, counts := range results {
		filename := ""

		if len(results) >= 1 {
			if filenames == nil {
				filename = ""
			}
			filename = filenames[i]
			if filename == "-" {
				filename = ""
			}
			fmt.Println(wc.FormatOutput(counts, opts, filename))
		}

		if len(results) < 1 {
			fmt.Println(wc.FormatOutput(total, opts, "total"))
		}
	}
}
