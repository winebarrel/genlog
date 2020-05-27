package main

import (
	"fmt"
	"genlog"
	"io"
	"log"
	"os"

	jsoniter "github.com/json-iterator/go"
)

func init() {
	log.SetFlags(0)
}

func main() {
	file := parseArgs()
	defer file.Close()

	err := genlog.Parse(file, func(block *genlog.Block) {
		line, err := jsoniter.MarshalToString(block)

		if err != nil {
			panic(err)
		}

		fmt.Println(line)
	})

	if err != nil {
		log.Fatal(err)
	}
}

func parseArgs() io.ReadCloser {
	fmt.Println(os.Args)

	if len(os.Args) > 2 {
		log.Fatalf("usage: %s GENERAL_LOG", os.Args[0])
	}

	if len(os.Args) == 1 {
		return os.Stdin
	}

	file, err := os.OpenFile(os.Args[1], os.O_RDONLY, 0)

	if err != nil {
		log.Fatal(err)
	}

	return file
}
