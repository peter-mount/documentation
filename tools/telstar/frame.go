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
  FrameType        string        `json:"frame-type,omitempty"`
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

func NewFrame(pid PageId) *Frame {
  parent := pid.PageNo / 10
  childId := pid.PageNo * 10
  return &Frame{
    PID:        pid,
    Visible:    true,
    HeaderText: "[G]area51.dev[Y]",
    FrameType:  "information",
    RoutingTable: []int{
      childId,     // 0 to first entry?
      childId + 1, // 1 to sub-page 1
      childId + 2, // 2 to sub-page 2
      childId + 3, // 3 to sub-page 3
      childId + 4, // 4 to sub-page 4
      childId + 5, // 5 to sub-page 5
      childId + 6, // 6 to sub-page 6
      childId + 7, // 7 to sub-page 7
      childId + 8, // 8 to sub-page 8
      childId + 9, // 9 to sub-page 9
      parent,      // Not sure which key this is for
    },
  }
}

func (f *Frame) SetTitle(t string) *Frame {
  f.Title.Data = t
  f.Title.Type = "markup"
  return f
}

func (f *Frame) SetContent(t string) *Frame {
  f.Content.Data = t
  f.Content.Type = "markup"
  return f
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
