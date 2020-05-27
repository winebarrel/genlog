package genlog

import (
	"bufio"
	"io"
	"regexp"
	"strings"
)

// https://github.com/mysql/mysql-server/blob/5.6/sql/log.cc#L1897
var reMySQL56 = regexp.MustCompile(`(?s)^(\d{6}\s+\d{1,2}:\d{2}:\d{2}|\t)\s+(\d+)\s+([^\t]+)\t(.*)`)

// https://github.com/mysql/mysql-server/blob/5.7/sql/log.cc#L783
var reMySQL57 = regexp.MustCompile(`(?s)^(\S+)\s+(\d+)\s+([^\t]+)\t(.*)`)

// https://github.com/mysql/mysql-server/blob/5.6/sql/log.cc#L1676
// https://github.com/mysql/mysql-server/blob/5.7/sql/log.cc#L696
var reIgnore = regexp.MustCompile(
	`^(?:` + strings.Join([]string{
		`\S+, Version: \S+ (.+). started with:`,
		`Tcp port: \d+  Unix socket: \S+`,
		`Time                 Id Command    Argument`,
	}, "|") + `)$`)

type Block struct {
	Time     string
	Id       string
	Command  string
	Argument string
}

func newBlock(tm, id, cmd, arg string) (*Block, *strings.Builder) {
	block := &Block{
		Time:    tm,
		Id:      id,
		Command: cmd,
	}

	argBldr := &strings.Builder{}
	argBldr.WriteString(arg)

	return block, argBldr
}

func callBack(block *Block, argBldr *strings.Builder, cb func(block *Block)) {
	block.Argument = argBldr.String()
	cb(block)
}

func Parse(file io.Reader, cb func(block *Block)) error {
	scanner := bufio.NewScanner(file)

	var block *Block
	var argBldr *strings.Builder
	var prevTm string

	for scanner.Scan() {
		line := scanner.Text()

		if reIgnore.MatchString(line) {
			continue
		}

		line += "\n"

		if m := reMySQL56.FindStringSubmatch(line); m != nil {
			tm := m[1]

			if tm == "\t" {
				tm = prevTm
			}

			if tm == "" {
				continue
			}

			if block != nil {
				callBack(block, argBldr, cb)
			}

			block, argBldr = newBlock(tm, m[2], m[3], m[4])
		} else if m := reMySQL57.FindStringSubmatch(line); m != nil {
			if block != nil {
				callBack(block, argBldr, cb)
			}

			block, argBldr = newBlock(m[1], m[2], m[3], m[4])
		} else if block != nil {
			argBldr.WriteString(line)
		}
	}

	if block != nil {
		callBack(block, argBldr, cb)
	}

	return scanner.Err()
}
