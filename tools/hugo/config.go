package hugo

import (
  "flag"
  "fmt"
  "github.com/peter-mount/go-kernel"
  "gopkg.in/yaml.v2"
  "io/ioutil"
  "path/filepath"
)

// Config file for doctool
type Config struct {
  Webserver  WebserverConfig `yaml:"webserver"` // Webserver
  PDF        PDF             `yaml:"pdf"`       // Common PDF template, can be overridden per book
  Books      []*Book         `yaml:"books"`     // Array of book definitions
  configFile *string         `yaml:"-"`         // Config file name
}

type WebserverConfig struct {
  Address string `yaml:"address"` // Webserver bind address, default "127.0.0.1"
  Port    int    `yaml:"port"`    // Port defaults to 8080
}

// PDF configuration
type PDF struct {
  PrintBackground     bool      `yaml:"printBackground"`     // Print background graphics
  Margin              PDFMargin `yaml:"margin"`              // Page margins
  Landscape           bool      `yaml:"landscape"`           // Landscape or Portrait
  Header              string    `yaml:"header"`              // Header in html
  Footer              string    `yaml:"footer"`              // Footer in html
  DisableHeaderFooter bool      `yaml:"disableHeaderFooter"` // Disable header & footer generation
}

// PDFMargin in inches. Default is 1cm ~0.4 inches
type PDFMargin struct {
  Top    float64 `yaml:"top"`    // Top margin in inches
  Left   float64 `yaml:"left"`   // Left margin in inches
  Bottom float64 `yaml:"bottom"` // Bottom margin in inches
  Right  float64 `yaml:"right"`  // Right margin in inches
}

func (c *Config) Name() string {
  return "Config"
}

func (c *Config) Init(_ *kernel.Kernel) error {
  c.configFile = flag.String("c", "", "The config file to use")
  return nil
}

func (c *Config) Start() error {
  // Verify then load the config file
  if *c.configFile == "" {
    *c.configFile = "config.yaml"
  }

  if filename, err := filepath.Abs(*c.configFile); err != nil {
    return err
  } else if in, err := ioutil.ReadFile(filename); err != nil {
    return err
  } else if err := yaml.Unmarshal(in, c); err != nil {
    return err
  }

  // Default webserver config
  if c.Webserver.Address == "" {
    c.Webserver.Address = "127.0.0.1"
  }
  if c.Webserver.Port == 0 {
    c.Webserver.Port = 8080
  }

  // Ensure all books have PDF config by using the global version
  for _, book := range c.Books {
    if book.PDF == nil {
      book.PDF = &c.PDF
    }

    // Sets book modified time before any generated files
    _=book.Modified()
  }

  return nil
}

func (c *Config) Webpath(f string, a ...interface{}) string {
  return fmt.Sprintf("http://%s:%d/"+f, append([]interface{}{c.Webserver.Address, c.Webserver.Port}, a...)...)
}
