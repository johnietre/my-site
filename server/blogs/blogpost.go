package blogs

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	utils "github.com/johnietre/utils/go"
)

func AddBlog(blog *Blog, content io.Reader) error {
  // TODO: add blog to database
  //f, err := os.Create(blogPath)
  f, err := os.Create("")
  if err != nil {
    return err
  }
  defer f.Close()
  io.TeeReader(content, f)
  // TODO: add edit
  return nil
}

func addBlog(blog *Blog, content io.Reader, attempts int) error {
  if attempts == 1000 {
    return fmt.Errorf("max attempts reached")
  }
  return fmt.Errorf("not implemented")
}

func EditBlog(edit *Edit, content io.Reader) error {
  return fmt.Errorf("TODO")
}

type Blog struct {
	Id        uint64 `json:"id"`
	Title     string `json:"title"`
  Authors []string `json:"authors"`
  Categories []string `json:"categories"`
  // Second-precision
	Timestamp int64  `json:"timestamp"`
  // Seconds
	TzOffset  int    `json:"tzOffset"`
  Edits []Edit `json:"edits"`
}

func parseBlogAuthors(s string) ([]string, bool) {
  l := len(s)
  if l < 4 {
    return nil, false
  }
  if s[0] != '|' || s[1] != '|' || s[l-1] != '|' || s[l-2] != '|' {
    return nil, false
  }
  s = s[2:l-2]
  return utils.MapSlice(
    strings.Split(s, "||"),
    func(s string) string {
      res, prevSlash := "", false
      for _, r := range s {
        if r == '\\' {
          if !prevSlash {
            prevSlash = true
            continue
          }
        }
        res += string(r)
        prevSlash = false
      }
      return res
    },
  ), true
}

func formatBlogAuthors(authors []string) string {
  authors = utils.MapSliceInPlace(
    authors,
    func(s string) string {
      res := ""
      for _, r := range s {
        if r == '\\' {
          res += `\\`
        } else if r == '|' {
          res += `\|`
        }
      }
      return res
    },
  )
  return "||"+strings.Join(authors, "||")+"||"
}

func parseBlogCategories(s string) ([]string, bool) {
  l := len(s)
  if l < 4 {
    return nil, false
  }
  if s[0] != '|' || s[1] != '|' || s[l-1] != '|' || s[l-2] != '|' {
    return nil, false
  }
  s = s[2:l-2]
  return utils.MapSlice(
    strings.Split(s, "||"),
    func(s string) string {
      res, prevSlash := "", false
      for _, r := range s {
        if r == '\\' {
          if !prevSlash {
            prevSlash = true
            continue
          }
        }
        res += string(r)
        prevSlash = false
      }
      return res
    },
  ), true
}

func formatBlogCategories(cats []string) string {
  cats = utils.MapSliceInPlace(
    cats,
    func(s string) string {
      res := ""
      for _, r := range s {
        if r == '\\' {
          res += `\\`
        } else if r == '|' {
          res += `\|`
        }
      }
      return res
    },
  )
  return "||"+strings.Join(cats, "||")+"||"
}

func (b Blog) FormatTimestamp() string {
  t := time.Unix(b.Timestamp, 0)
  offMin := (b.TzOffset / 60) % 60
  offHr := (b.TzOffset / 60 / 60) % 24
  tzName := fmt.Sprintf("UTC%+02d:%02d", offHr, offMin)
  t = t.In(time.FixedZone(tzName, b.TzOffset))
  return t.Format(time.RFC1123)
}

type Edit struct {
  BlogId int64 `json:"blogId"`
  Description string `json:"description"`
  // Second-precision
	Timestamp int64  `json:"timestamp"`
  // Seconds
	TzOffset  int    `json:"tzOffset"`
  PrevHash string `json:"prevHash"`
  Hash string `json:"hash"`
}

func (e Edit) FormatTimestamp() string {
  t := time.Unix(e.Timestamp, 0)
  offMin := (e.TzOffset / 60) % 60
  offHr := (e.TzOffset / 60 / 24) % 24
  tzName := fmt.Sprintf("UTC%+02d:%02d", offHr, offMin)
  t = t.In(time.FixedZone(tzName, e.TzOffset))
  return t.Format(time.RFC1123)
}

type BlogPostPageData struct {
  Blog Blog
  // The current edit
  Edit Edit
  ContentUrl string
}

func NewBlogPostPageData() BlogPostPageData {
  return BlogPostPageData{
    Blog: Blog{
    },
    Edit: Edit{
    },
    ContentUrl: "",
  }
}
