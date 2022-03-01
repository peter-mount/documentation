package telstar

import (
  "bytes"
  "encoding/json"
  "fmt"
  "github.com/peter-mount/go-kernel/util/task"
  "log"
  "net/http"
  "net/http/cookiejar"
  "net/url"
)

// Service implements a bridge which updates a telstar server with generated documentation
type Service struct {
  enabled *string    `kernel:"flag,telstar,Update telstar instance"`
  worker  task.Queue `kernel:"worker"`
  cookies http.CookieJar
  uri     *url.URL
}

func (s *Service) Start() error {
  if *s.enabled != "" {
    // Setup cookie jar, parse parameter url then login to the remote service
    jar, err := cookiejar.New(nil)
    if err != nil {
      return err
    }
    s.cookies = jar

    u, err := url.Parse(*s.enabled)
    if err != nil {
      return err
    }
    s.uri = u
    /*
       s.worker.AddPriorityTask(-1, func(ctx context.Context) error {
         frame, err := s.GetFrame(PageId{PageNo: 9, FrameId: "a"})
         if err != nil {
           s.Printf("error %v", err)
           return err
         }
         s.Printf("Frame: %v", frame)
         os.Exit(0)
         return nil
       })
    */

    //walk.NewPathWalker().
    //  Walk()

    return s.loginTelstar()
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
