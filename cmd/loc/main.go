package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mdw-go/leaderboard"
)

var Version = "dev"

func main() {
	flags := flag.NewFlagSet(fmt.Sprintf("%s @ %s", filepath.Base(os.Args[0]), Version), flag.ExitOnError)
	flags.Usage = func() {
		_, _ = fmt.Fprintf(flags.Output(), "Usage of %s:\n", flags.Name())
		_, _ = fmt.Fprintln(flags.Output(), ""+
			"Renders a table of line counts per file type in the directory specified by the "+
			"first non-flag argument or, if no argument is provided, the current directory.",
		)
		flags.PrintDefaults()
	}
	_ = flags.Parse(os.Args[1:])

	root := "."

	if len(os.Args) > 1 {
		root = os.Args[1]
	}

	stats := make(leaderboard.Map[string, int])

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && strings.HasPrefix(d.Name(), ".") && len(d.Name()) > 1 {
			return filepath.SkipDir
		}
		if d.IsDir() {
			return nil
		}
		if d.Name() == ".DS_Store" {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if ext == "" {
			ext = "(no extension)"
		}
		if ext == d.Name() {
			return nil
		}
		lines, err := countLines(path)
		if err != nil {
			log.Printf("Error reading %s: %v\n", path, err)
			return nil
		}
		stats[ext] += lines
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking directory: %v\n", err)
	}

	// Print results
	fmt.Printf("%-15s %10s\n", "Extension", "Lines")
	fmt.Println(strings.Repeat("-", 26))
	for _, ext := range stats.TopN(len(stats)) {
		fmt.Printf("%-15s %10d\n", ext, stats[ext])
	}
}

// countLines returns the number of lines in a given file.
func countLines(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		count++
	}
	return count, scanner.Err()
}
