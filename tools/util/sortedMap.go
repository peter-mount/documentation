package util

type SortedMap map[string]interface{}

func NewSortedMap() *SortedMap {
  m := make(SortedMap)
  return &m
}

func (m *SortedMap) AddAll(s map[string]interface{}) *SortedMap {
  for k, v := range s {
    (*m)[k] = v
  }
  return m
}

func (m *SortedMap) DecodeMap(s map[interface{}]interface{}) *SortedMap {
  for k, v := range s {
    ks := DecodeString(k, "")
    if ks != "" {
      (*m)[ks] = v
    }
  }
  return m
}

func (m *SortedMap) Decode(v interface{}) *SortedMap {
  _ = IfMap(v, func(sm map[interface{}]interface{}) error {
    m.DecodeMap(sm)
    return nil
  })
  return m
}

func (m *SortedMap) Keys() StringSlice {
  var a StringSlice
  for k, _ := range *m {
    a = append(a, k)
  }
  return a
}

func (m *SortedMap) ForEach(f func(string, interface{}) error) error {
  return m.Keys().
    Sort().
      ForEach(func(k string) error {
        return f(k, (*m)[k])
      })
}
