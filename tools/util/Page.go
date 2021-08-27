package util

import (
  "gopkg.in/yaml.v2"
  "io/ioutil"
  "log"
  "os"
  "path"
  "strconv"
  "strings"
  "time"
)

func WriteFile(fileName string, data []byte, fileTime time.Time) error {
  err := os.MkdirAll(path.Dir(fileName), 0755)
  if err != nil {
    return err
  }

  err = ioutil.WriteFile(fileName, data, 0644)
  if err != nil {
    return err
  }

  return os.Chtimes(fileName, fileTime, fileTime)
}

func WriteReferenceFile(dir, name, title, desc string, weight int, fileTime time.Time, handler func([]string) ([]string, error)) error {
  n := path.Join(dir, "reference", name, "_index.html")
  log.Printf("Generating %s", n)

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

  a, err := handler(a)
  if err != nil {
    return err
  }

  a = append(a, "---")

  return WriteFile(n, []byte(strings.Join(a, "\n")), fileTime)
}

func WriteReferenceYaml(dir, name, title, desc string, weight int, fileTime time.Time, val interface{}) error {
  return WriteReferenceFile(dir, name, title, desc, weight, fileTime, func(a []string) ([]string, error) {
    c, err := yaml.Marshal(val)
    if err != nil {
      return nil, err
    }
    a = append(a, string(c[:]))

    return a, nil
  })
}
