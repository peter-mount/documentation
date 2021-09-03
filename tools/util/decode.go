package util

import (
  "strconv"
)

// ForEachInterface will invoke a function for every entry in v if v is a slice
func ForEachInterface(v interface{}, f func(interface{})) {
  if a, ok := v.([]interface{}); ok {
    for _, e := range a {
      f(e)
    }
  }
}

// IfMap will invoke a function if v is a map
func IfMap(v interface{}, f func(map[interface{}]interface{})) {
  if m, ok := v.(map[interface{}]interface{}); ok {
    f(m)
  }
}

// IfMapEntry invokes a function if a map contains an entry
func IfMapEntry(m map[interface{}]interface{}, n interface{}, f func(interface{})) {
  if v, ok := m[n]; ok {
    f(v)
  }
}

func DecodeString(v interface{}, def string) string {
  var r string

  if v != nil {
    if s, ok := v.(string); ok {
      r = s
    } else if i, ok := v.(int); ok {
      r = strconv.Itoa(i)
    }
  }

  if r == "" {
    r = def
  }

  return r
}

func DecodeInt(v interface{}, def int) (int, bool) {
  r := def

  if v != nil {
    if s, ok := v.(string); ok {
      i, err := strconv.Atoi(s)
      if err != nil {
        return 0, false
      }
      r = i
    } else if i, ok := v.(int); ok {
      r = i
    }
  }

  return r, true
}

func DecodeBool(v interface{}) (bool, bool) {

  if v != nil {
    if s, ok := v.(string); ok {
      b, err := strconv.ParseBool(s)
      if err != nil {
        return false, false
      }
      return b, true
    } else if i, ok := v.(int); ok {
      return i != 0, true
    }

    b, ok := v.(bool)
    return b, ok
  }

  return false, true
}

func IfMapEntryString(m map[interface{}]interface{}, n string) string {
  var s string
  IfMapEntry(m, n, func(i interface{}) {
    s = DecodeString(i, "")
  })
  return s
}

func IfMapEntryBool(m map[interface{}]interface{}, n string) bool {
  var s bool
  IfMapEntry(m, n, func(i interface{}) {
    s, _ = DecodeBool(i)
  })
  return s
}
