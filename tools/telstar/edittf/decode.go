package edittf

import (
  "strings"
)

// Decode an EditIF and returns the raw content
func Decode(s string) []string {
  return DecodeFrame(s, 1, 22, 0, 39, true)
}

func DecodeFrame(s string, rowBegin, rowEnd, columnBegin, columnEnd int, trimEnds bool) []string {
  data := DecodeUrl(s)

  var r []string

  for row := rowBegin; row <= rowEnd; row++ {
    rd := data[row*40 : (row+1)*40]
    rd = rd[columnBegin : columnEnd+1]
    if trimEnds {
      rd = strings.TrimRight(rd, " ")
    }

    r = append(r, rd)
  }

  return r
}

const (
  alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
)

func DecodeUrl(s string) string {
  data := s

  // Strip everything before the #, if it's there
  if p := strings.Index(data, "#"); p > -1 {
    data = data[p+1:]
  }

  // Strip everything before :
  if p := strings.Index(data, ":"); p > -1 {
    data = data[p+1:]
  }

  // check length is 1120 or 1167 characters long

  // []byte of 1000, initially 0
  decoded_data := make([]byte, 1000)

  for index, c := range data {
    findex := strings.Index(alphabet, string(c))
    for b := 0; b < 6; b++ {
      bit := findex & (1 << (5 - b))
      if bit > 0 {
        cbit := (index * 6) + b
        cpos := cbit % 7
        cloc := (cbit - cpos) / 7
        decoded_data[cloc] |= 1 << (6 - cpos)
      }
    }
  }

  return string(decoded_data)
}
