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
// NOTE: In Aurora MySQL 5.7, there may be no space between "Time" and "ID"
var reMySQL57 = regexp.MustCompile(`(?s)^([^\sZ]+Z)\s*(\d+)\s+([^\t]+)\t(.*)`)

// https://github.com/mysql/mysql-server/blob/5.6/sql/log.cc#L1676
// https://github.com/mysql/mysql-server/blob/5.7/sql/log.cc#L696
var reIgnore = regexp.MustCompile(
	`^(?:` + strings.Join([]string{
		`\S+, Version: \S+ (.+). started with:`,
		`Tcp port: \d+  Unix socket: \S+`,
		`Time                 Id Command    Argument`,
	}, "|") + `)$`)

var ReadLineBufSize = 4096

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
	arg := strings.TrimRight(argBldr.String(), "\n")
	block.Argument = arg
	cb(block)
}

func Parse(file io.Reader, cb func(block *Block)) error {
	reader := bufio.NewReader(file)

	var block *Block
	var argBldr *strings.Builder
	var prevTm string

	for {
		rawLine, err := readLine(reader)

		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		line := string(rawLine)

		if reIgnore.MatchString(line) {
			continue
		}

		line += "\n"

		if m := reMySQL56.FindStringSubmatch(line); m != nil {
			tm := m[1]

			if tm == "\t" {
				tm = prevTm
			}

			prevTm = tm

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

	return nil
}

func readLine(reader *bufio.Reader) ([]byte, error) {
	buf := make([]byte, 0, ReadLineBufSize)
	var err error

	for {
		line, isPrefix, e := reader.ReadLine()
		err = e

		if len(line) > 0 {
			buf = append(buf, line...)
		}

		if !isPrefix || err != nil {
			break
		}
	}

	return buf, err
}
