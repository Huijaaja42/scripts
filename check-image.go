package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type image struct {
	path   string
	errors []string
}

func parse(s *bufio.Scanner) (o []image, err error) {
	for s.Scan() {
		t := s.Text()
		ts := strings.Split(t, " ")
		if ts[0] == "Checking" {
			o = append(o, image{path: strings.Join(ts[2:], " ")})
		} else if len(o) > 0 {
			o[len(o)-1].errors = append(o[len(o)-1].errors, t)
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return o, nil
}

func filter(i []image) (o []image) {
	for _, v := range i {
		if len(v.errors) == 0 {
			continue
		}
		s := strings.Split(v.path, "/")
		if s[len(s)-1][0] == '.' {
			continue
		}
		if strings.Contains(v.path, "$RECYCLE.BIN") {
			continue
		}
		e := strings.Join(v.errors, " ")
		if strings.Contains(e, "No such file or directory") {
			continue
		}
		if strings.Contains(e, "unable to decode APP fields") {
			continue
		}
		if strings.Contains(e, "overread 8") {
			continue
		}
		o = append(o, v)
	}
	return
}

func main() {
	var f *os.File
	if len(os.Args) > 1 {
		var err error
		f, err = os.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
	} else {
		f = os.Stdin
	}
	o, err := parse(bufio.NewScanner(f))
	if err != nil {
		log.Fatal(err)
	}
	fo := filter(o)
	fmt.Printf("Found %d images with errors\n\n", len(fo))
	for _, v := range fo {
		fmt.Println(v.path)
		for _, e := range v.errors {
			fmt.Println(e)
		}
		fmt.Println("")
	}
}
