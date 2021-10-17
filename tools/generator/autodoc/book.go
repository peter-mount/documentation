package autodoc

import (
  "context"
  "github.com/peter-mount/documentation/tools/generator"
  "github.com/peter-mount/documentation/tools/util/autodoc"
  "strings"
  "time"
)

func buildHeaderFile(topic string, path string) autodoc.Handler {
  return func(ctx context.Context, builder autodoc.Builder) error {
    b := generator.GetBook(ctx)

    builder.Separator().
      Comment("%s for %s", topic, b.Title)

    if b.SubTitle != "" {
      builder.Comment(b.SubTitle)
    }
    if b.Author != "" {
      builder.Comment("Author: %s", b.Author)
    }
    if b.SubAuthor != "" {
      builder.Comment("SubAuthor: %s", b.SubAuthor)
    }

    builder.Comment("").
      Comment("URL: https://area51.dev/%s/", strings.Join(strings.Split(b.ContentPath(), "/")[1:], "/")).
      Comment("").
      Comment("Modified: %s", b.Modified().Format(time.RFC1123))

    if len(path) > 0 {
      builder.Comment("").
        Comment("Current version: https://area51.dev/%s", strings.Join(strings.Split(path, "/")[1:], "/"))
    }

    builder.Separator().
      Newline()

    return nil
  }
}
