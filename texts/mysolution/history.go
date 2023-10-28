package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

// Regex pattern to look for a match which starts with ';go ' and subsequent word can be anything between letters a - z (can have more than one letter)
var cmdRE = regexp.MustCompile(`;go ([a-z]+)`)

// cmdFreq returns the frequency of "go" subcommand usage in ZSH history
func cmdFreq(fileName string) (map[string]int, error) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var freqs = make(map[string]int)
	s := bufio.NewScanner(file)
	for s.Scan() {
		matches := cmdRE.FindStringSubmatch(s.Text())
		if len(matches) == 0 {
			continue
		}
		// we need only the first leftmost match
		cmd := matches[1]
		freqs[cmd]++
	}

	if err := s.Err(); err != nil {
		return nil, err
	}

	return freqs, nil
}

func main() {
	freqs, err := cmdFreq("zsh_history")
	if err != nil {
		log.Fatal(err)
	}

	for cmd, count := range freqs {
		fmt.Printf("%s -> %d\n", cmd, count)
	}
}
