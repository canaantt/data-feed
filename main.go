package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/karrick/godirwalk"
)

type Dirent struct {
	name     string
	modeType os.FileMode
}

func rsl(fn string, n int) (string, error) {
	if n < 1 {
		return "", fmt.Errorf("invalid request: line %d", n)
	}
	f, err := os.Open(fn)
	if err != nil {
		return "", err
	}
	defer f.Close()
	bf := bufio.NewReader(f)
	var line string
	for lnum := 0; lnum < n; lnum++ {
		line, err = bf.ReadString('\n')
		if err == io.EOF {
			switch lnum {
			case 0:
				return "", errors.New("no lines in file")
			case 1:
				return "", errors.New("only 1 line")
			default:
				return "", fmt.Errorf("only %d lines", lnum)
			}
		}
		if err != nil {
			return "", err
		}
	}
	if line == "" {
		return "", fmt.Errorf("line %d empty", n)
	}
	return line, nil
}

func main() {
	dirname := "."
	if len(os.Args) > 1 {
		dirname = os.Args[1]
	}
	err := godirwalk.Walk(dirname, &godirwalk.Options{
		// Unsorted: true, // set true for faster yet non-deterministic enumeration (see godoc)
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			// fmt.Printf("%s %s\n", de.ModeType(), osPathname)
			if de.ModeType()&os.ModeDir != 0 {
				// fmt.Printf(" %s is directory. => %s |  \n", de.ModeType()&os.ModeDir, osPathname)
			} else {
				// match, _ := regexp.MatchString("p([a-z]+)ch", osPathname)
				// fmt.Println(match)
				r, _ := regexp.Compile("(clinical[w+\\_+]patient)")
				if r.MatchString(osPathname) {
					header, err := rsl(osPathname, 1)
					if err == nil {
						fmt.Println(osPathname)
						fmt.Println(header)
					} else {
						fmt.Println(err)
					}

				}

			}
			return nil
		},
		ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
			// Your program may want to log the error somehow.
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)

			// For the purposes of this example, a simple SkipNode will suffice,
			// although in reality perhaps additional logic might be called for.
			return godirwalk.SkipNode
		},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
