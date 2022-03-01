package telstar

import (
  "encoding/json"
  "io/ioutil"
  "net/http"
)

type Frame struct {
  PID              PageId        `json:"pid"`
  Visible          bool          `json:"visible"`
  HeaderText       string        `json:"header-text"`
  CacheId          string        `json:"cache-id"`
  Content          Content       `json:"content"`
  Title            Content       `json:"title"`
  RoutingTable     []int         `json:"routing-table"`
  Cursor           bool          `json:"cursor"`
  Connection       Connection    `json:"connection"`
  AuthorId         string        `json:"author-id"`
  ResponseFields   []interface{} `json:"response-fields"`
  NavMsgForeColour string        `json:"navmsg-forecolour"`
  NavMsgBackColour string        `json:"navmsg-backcolour"`
  NavMsgHighlight  string        `json:"navmsg-highlight"`
}

type PageId struct {
  PageNo  int    `json:"page-no"`
  FrameId string `json:"frame-id"`
}

type Connection struct {
  Address string `json:"address"`
  Mode    string `json:"mode"`
  Port    int    `json:"port"`
}

type Content struct {
  Data string `json:"data"`
  Type string `json:"type"`
}

func (s *Service) GetFrame(pid PageId) (*Frame, error) {
  s.Printf("GetFrame %d%s", pid.PageNo, pid.FrameId)

  var frame *Frame

  if err := s.Call("GET",
    nil,
    func(resp *http.Response) error {

      body, err := ioutil.ReadAll(resp.Body)
      if err != nil {
        return err
      }
      frame = &Frame{PID: pid}
      return json.Unmarshal(body, frame)
    },
    "frame/%d%s", pid.PageNo, pid.FrameId); err != nil {
    return nil, err
  }

  return frame, nil
}
