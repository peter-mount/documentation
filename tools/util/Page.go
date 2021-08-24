package util

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

func WriteReferenceFile(dir, name, title, desc string, weight int, handler func([]string) ([]string, error)) error {
	n := path.Join(dir, "reference", name, "_index.html")
	log.Printf("Generating %s", n)

	err := os.MkdirAll(path.Dir(n), 0755)
	f, err := os.Create(n)
	if err != nil {
		return err
	}
	defer f.Close()

	a := []string{"---",
		"# This file is generated.",
		"# To edit, change the files under content then run the generator.",
		"type: \"6502opcode\"",
		"title: \"" + title + "\"",
		"linkTitle: \"" + title + "\"",
		"weight: " + strconv.Itoa(weight),
		"description: \"" + desc + "\"",
		"notitle: \"don't show instructions header in table\"",
	}

	a, err = handler(a)
	if err != nil {
		return err
	}

	a = append(a, "---")
	_, err = f.Write([]byte(strings.Join(a, "\n")))
	if err != nil {
		return err
	}

	return nil
}

func WriteReferenceYaml(dir, name, title, desc string, weight int, val interface{}) error {
	return WriteReferenceFile(dir, name, title, desc, weight, func(a []string) ([]string, error) {
		c, err := yaml.Marshal(val)
		if err != nil {
			return nil, err
		}
		a = append(a, string(c[:]))

		return a, nil
	})
}
