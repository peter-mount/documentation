package util

type StringSlice []string

func (s StringSlice) ForEach(f func(string) error) error {
  for _, b := range s {
    err := f(b)
    if err != nil {
      return err
    }
  }
  return nil
}
