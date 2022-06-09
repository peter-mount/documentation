package m68k

import (
  "context"
  "github.com/peter-mount/documentation/tools"
  "github.com/peter-mount/documentation/tools/generator"
  "github.com/peter-mount/go-kernel/util/task"
  "strconv"
  "strings"
)

func delayOpTask(t task.Task) task.Task {
  return func(ctx context.Context) error {
    task.GetQueue(ctx).AddPriorityTask(tools.PriorityIndices,
      task.Of(t).
        WithContext(ctx, generator.BookKey))
    return nil
  }
}

func decodeOpcode(s string) int64 {
  s = strings.ReplaceAll(s, "nn", "00")
  for len(s) < 8 {
    s = s + "00"
  }
  if i, err1 := strconv.ParseInt(s, 16, 64); err1 == nil {
    return i
  }
  return -1
}
