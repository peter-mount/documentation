package util

import (
	"bufio"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"path"
	"strings"
)

func LoadFrontMatter(d, n string) (map[string]interface{}, error) {
	f, err := os.Open(path.Join(d, n))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return ReadFrontMatter(f)
}

func ReadFrontMatter(r io.Reader) (map[string]interface{}, error) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		l := scanner.Text()
		if l == "---" {
			return readFrontMatter(scanner)
		}
	}

	// No frontmatter
	return nil, nil
}

func readFrontMatter(s *bufio.Scanner) (map[string]interface{}, error) {
	var a []string
	for s.Scan() {
		l := s.Text()
		if l == "---" {
			m := make(map[string]interface{})

			b := []byte(strings.Join(a, "\n"))
			err := yaml.Unmarshal(b, m)
			if err != nil {
				return nil, err
			}
			return m, nil
		}

		a = append(a, l)
	}

	// no --- terminator so ignore the page
	return nil, nil
}
