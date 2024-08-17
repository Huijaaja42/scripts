package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type image struct {
	path   string
	errors []string
}

type filter struct {
	paths       []string
	errors      []string
	hiddenDirs  bool
	hiddenFiles bool
}

func parseInput(s *bufio.Scanner) (out []image, err error) {
	for s.Scan() {
		t := s.Text()
		if len(t) == 0 {
			continue
		}

		ts := strings.Split(t, " ")
		if ts[0] == "Checking" {
			out = append(out, image{path: strings.Join(ts[2:], " ")})
		} else if len(out) > 0 {
			out[len(out)-1].errors = append(out[len(out)-1].errors, t)
		}
	}

	if err := s.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func filterImages(input []image, filter *filter) (out []image) {
	for _, v := range input {
		if len(v.errors) == 0 {
			continue
		}

		if !filter.hiddenDirs {
			if strings.Contains(v.path, "/.") {
				continue
			}
		}

		if !filter.hiddenFiles {
			s := strings.Split(v.path, "/")
			if s[len(s)-1][0] == '.' {
				continue
			}
		}

		match := false
		for _, f := range filter.paths {
			if strings.Contains(v.path, f) {
				match = true
				break
			}
		}
		if match {
			continue
		}

		e := strings.Join(v.errors, " ")
		for _, f := range filter.errors {
			if strings.Contains(e, f) {
				match = true
				break
			}
		}
		if match {
			continue
		}

		out = append(out, v)
	}
	return
}

func main() {
	verbose := flag.Bool("v", false, "verbose output")
	all := flag.Bool("a", false, "include hidden files/directories")
	flag.Parse()

	var file *os.File
	if flag.Arg(0) == "" {
		file = os.Stdin
	} else {
		var err error
		file, err = os.Open(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
	}

	input, err := parseInput(bufio.NewScanner(file))
	if err != nil {
		log.Fatal(err)
	}

	filter := &filter{
		paths: []string{
			"$RECYCLE.BIN",
		},
		errors: []string{
			"No such file or directory",
			"unable to decode APP fields",
			"overread 8",
		},
		hiddenDirs:  *all,
		hiddenFiles: *all,
	}

	out := filterImages(input, filter)
	if *verbose {
		fmt.Printf("Found %d images with errors\n", len(out))
		for _, v := range out {
			fmt.Println("")
			fmt.Println(v.path)
			for _, e := range v.errors {
				fmt.Println(e)
			}
		}
	} else {
		for _, v := range out {
			fmt.Println(v.path)
		}
	}
}
