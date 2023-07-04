package telstar

import (
	"context"
	"errors"
	"net/http/cookiejar"
	"strconv"
)

type login struct {
	UN int    `json:"un"`
	PW string `json:"pw"`
}

func (s *Service) loginTelstar(_ context.Context) error {
	// If not enabled or already got cookies then don't try to login
	if *s.enabled == "" || s.cookies != nil {
		return nil
	}

	s.Printf("Login")

	// Setup cookie jar, parse parameter url then login to the remote service
	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}
	s.cookies = jar

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
