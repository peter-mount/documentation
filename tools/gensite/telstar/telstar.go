package telstar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/peter-mount/documentation/tools/gensite/generator"
	"github.com/peter-mount/documentation/tools/gensite/telstar/frame"
	"github.com/peter-mount/go-kernel/v2/log"
	"github.com/peter-mount/go-kernel/v2/util/task"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Service implements a bridge which updates a telstar server with generated documentation
type Service struct {
	enabled   *string              `kernel:"flag,telstar,Update telstar instance"`
	worker    task.Queue           `kernel:"worker"`
	generator *generator.Generator `kernel:"inject"` // Generator
	cookies   http.CookieJar
	uri       *url.URL
}

func (s *Service) Start() error {
	s.generator.Register("telstar",
		task.Of().
			Then(s.loginTelstar))

	if *s.enabled != "" {
		u, err := url.Parse(*s.enabled)
		if err != nil {
			return err
		}
		s.uri = u
	}

	return nil
}

func (s *Service) Printf(f string, a ...interface{}) *Service {
	log.Printf("Telstar: "+f, a...)
	return s
}

type ResponseHandler func(resp *http.Response) error

func (s *Service) Call(method string, payload interface{}, responseHandler ResponseHandler, name string, args ...interface{}) error {
	var b []byte

	if payload != nil {
		var err error
		b, err = json.Marshal(payload)
		if err != nil {
			return err
		}
	}

	a := []interface{}{s.uri.Scheme, s.uri.Hostname(), s.uri.Port()}
	a = append(a, args...)
	uri := fmt.Sprintf("%s://%s:%s/"+name, a...)
	s.Printf("http %q", uri)

	req, err := http.NewRequest(method, uri, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	for _, c := range s.cookies.Cookies(req.URL) {
		s.Printf("add cookie %v", c)
		req.AddCookie(c)
	}

	//s.Printf("query %v", req)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	s.Printf("Status %d %q", resp.StatusCode, resp.Status)

	for _, c := range resp.Cookies() {
		s.Printf("set cookie %v", c)
	}
	s.cookies.SetCookies(req.URL, resp.Cookies())

	if responseHandler != nil {
		return responseHandler(resp)
	}

	return nil
}

func (s *Service) GetFrame(pid frame.PageId) (*frame.Frame, error) {
	s.Printf("GetFrame %d%s", pid.PageNo, pid.FrameId)

	var f *frame.Frame

	if err := s.Call("GET",
		nil,
		func(resp *http.Response) error {

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			f = &frame.Frame{PID: pid}
			return json.Unmarshal(body, f)
		},
		"frame/%d%s", pid.PageNo, pid.FrameId); err != nil {
		return nil, err
	}

	return f, nil
}
