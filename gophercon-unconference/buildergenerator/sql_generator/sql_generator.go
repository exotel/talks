// +build OMIT

package main

import (
	"bufio"
	"go/parser"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	err = filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if !strings.HasSuffix(path, ".sql") {
			return nil
		}

		return parser.ParseFile(path, func(raw []string) {
			invocation, err := parseCommand(trimPath(path), raw)
			if err != nil {
				log.Printf("Could not parse command found in %v, %v", path, raw)
				log.Fatal(err)
			}
			invocations = append(invocations, invocation)
		})
	})
}

func (p *Parser) Parse(reader io.Reader, callback Handler) error {
	var currentCommand []string

	endCommand := func() {
		if currentCommand != nil {
			callback(currentCommand)
		}
		currentCommand = nil
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, p.Comment) {
			endCommand()
			continue
		}

		line = strings.TrimSpace(strings.TrimPrefix(line, p.Comment))

		if currentCommand != nil {
			if line == "" {
				endCommand()
				continue
			}

			currentCommand = append(currentCommand, line)
			continue
		}

		if strings.HasPrefix(line, p.Command) {
			line := strings.TrimSpace(strings.TrimPrefix(line, p.Command))
			currentCommand = []string{line}
		}
	}

	endCommand()
	return scanner.Err()
}
