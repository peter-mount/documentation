package util

import (
  "github.com/peter-mount/go-kernel/util"
  "log"
  "sort"
  "strings"
)

type Notes struct {
  Notes  []*Note
  lookup map[string]*Note
}

type Note struct {
  Key   int
  Value string
}

func NewNotes() *Notes {
  return &Notes{lookup: make(map[string]*Note)}
}

func (n *Notes) Get(s string) *Note {
  return n.lookup[s]
}

func (n *Notes) GetId(i int) *Note {
  if i < 1 || i > len(n.Notes) {
    return nil
  }
  return n.Notes[i-1]
}

func (n *Notes) Add(s string) *Note {
  if s == "" {
    return nil
  }

  if note, exists := n.lookup[s]; exists {
    return note
  }

  note := &Note{Value: s}
  n.Notes = append(n.Notes, note)
  n.lookup[s] = note
  return note
}

func (n *Notes) add(e interface{}) *Note {
  if s, ok := e.(string); ok {
    return n.Add(s)
  } else {
    log.Println(e)
  }
  return nil
}

func (n *Note) Compare(b *Note) bool {
  if strings.ToLower(n.Value) == strings.ToLower(b.Value) {
    log.Printf("*** %q", n.Value)
  }
  return strings.ToLower(n.Value) < strings.ToLower(b.Value)
}

func (n *Notes) Normalise() {
  sort.SliceStable(n.Notes, func(i, j int) bool {
    return n.Notes[i].Compare(n.Notes[j])
  })

  // Now change key to match location in the slice
  for i, v := range n.Notes {
    v.Key = i + 1
  }
}

func (n *Notes) DecodePageNotes(v interface{}) {
  if v == nil {
    return
  }

  _ = util.ForEachInterface(v, func(e interface{}) error {
    _ = n.add(e)
    return nil
  })

  return
}

// Merge merges the notes in b into this instance
func (n *Notes) Merge(b *Notes) {
  if b == nil {
    return
  }

  // Add to this entry then replace the instance in b so they share the same instance
  for i, e := range b.Notes {
    k := e.Value
    c := n.Add(k)
    b.Notes[i] = c
    b.lookup[k] = c
  }
}
