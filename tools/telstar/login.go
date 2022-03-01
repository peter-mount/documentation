package telstar

import (
  "errors"
  "strconv"
)

type login struct {
  UN int    `json:"un"`
  PW string `json:"pw"`
}

func (s *Service) loginTelstar() error {
  s.Printf("Login")

  un, err := strconv.Atoi(s.uri.User.Username())
  if err != nil {
    return err
  }
  pw, ok := s.uri.User.Password()
  if !ok {
    return errors.New("no secret provided")
  }
  payload := &login{UN: un, PW: pw}

  return s.Call("PUT", payload, nil, "login")
}
