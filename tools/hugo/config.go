package hugo

import (
  "fmt"
  "gopkg.in/yaml.v2"
  "io/ioutil"
  "path/filepath"
)

// Config file for doctool
type Config struct {
  Webserver  WebserverConfig `yaml:"webserver"`                            // Webserver
  PDF        PDF             `yaml:"pdf"`                                  // Common PDF template, can be overridden per book
  configFile *string         `yaml:"-" kernel:"flag,c,Config file to use"` // Config file name
}

type WebserverConfig struct {
  Address string `yaml:"address"` // Webserver bind address, default "127.0.0.1"
  Port    int    `yaml:"port"`    // Port defaults to 8080
}

// PDF configuration
type PDF struct {
  PrintBackground     bool      `yaml:"printBackground"`     // Print background graphics
  Margin              PDFMargin `yaml:"margin"`              // Page margins
  Width               float64   `yaml:"width"`               // Page width in inches
  Height              float64   `yaml:"height"`              // Page height in inches
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

func (c *Config) Start() error {
  // Verify then load the config file
  if *c.configFile == "" {
    *c.configFile = "config.yaml"
  }

  filename, err := filepath.Abs(*c.configFile)
  if err != nil {
    return err
  }

  in, err := ioutil.ReadFile(filename)
  if err != nil {
    return err
  }

  err = yaml.Unmarshal(in, c)
  if err != nil {
    return err
  }

  // Default webserver config
  if c.Webserver.Address == "" {
    c.Webserver.Address = "127.0.0.1"
  }
  if c.Webserver.Port == 0 {
    c.Webserver.Port = 8080
  }

  return nil
}

func (c *Config) WebPath(f string, a ...interface{}) string {
  return fmt.Sprintf("http://%s:%d/"+f, append([]interface{}{c.Webserver.Address, c.Webserver.Port}, a...)...)
}
