package blogs

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	utils "github.com/johnietre/utils/go"
)

const initialHash = "0000000000000000000000000000000000000000000000000000000000000000"

func AddBlog(blog *Blog, content io.Reader) error {
	// TODO: add blog to database
	//f, err := os.Create(blogPath)
	f, err := os.Create("")
	if err != nil {
		return err
	}
	defer f.Close()
  blog.Edits = []Edit{
    {
      Description: "GENESIS",
      PrevHash: initialHash,
    },
  }
  blog.Edits[0].Hash, err = blog.Hash(io.TeeReader(content, f))
  if err != nil {
    return err
  }
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
	Id         uint64   `json:"id"`
	Title      string   `json:"title"`
	Authors    []string `json:"authors"`
	Categories []string `json:"categories"`
	// Second-precision
	Timestamp int64 `json:"timestamp"`
	// Seconds
	TzOffset int    `json:"tzOffset"`
	Edits    []Edit `json:"edits"`
}

func (b *Blog) Hash(content io.Reader) (string, error) {
  toWrite := utils.NewSlice[[]byte](nil)

  hash := sha256.New()
  toWrite.PushBack([]byte(b.Title))
  if len(b.Authors) != 0 {
    toWrite.PushBack([]byte(b.Authors[0]))
    for i := 1; i < len(b.Authors); i++ {
      toWrite.PushBack([]byte("\n"))
      toWrite.PushBack([]byte(b.Authors[i]))
    }
  }
  if len(b.Categories) != 0 {
    toWrite.PushBack([]byte(b.Categories[0]))
    for i := 1; i < len(b.Categories); i++ {
      toWrite.PushBack([]byte("\n"))
      toWrite.PushBack([]byte(b.Categories[i]))
    }
  }
  toWrite.PushBack(binary.BigEndian.AppendUint64(nil, uint64(b.Timestamp)))
  toWrite.PushBack(binary.BigEndian.AppendUint32(nil, uint32(b.TzOffset)))
  switch l := len(b.Edits); l {
  case 0:
    return "", fmt.Errorf("must have at least one edit")
  case 1:
    // Is initial/genesis
    toWrite.PushBack([]byte(b.Edits[0].PrevHash))
  default:
    toWrite.PushBack([]byte(b.Edits[l].Hash))
  }
  for _, b := range *toWrite.Ptr {
    if _, err := hash.Write(b); err != nil {
      return "", err
    }
  }
  if _, err  := io.Copy(hash, content); err != nil {
    return "", err
  }
  return hex.Dump(hash.Sum(nil)), nil
}

func parseBlogAuthors(s string) ([]string, bool) {
	l := len(s)
	if l < 4 {
		return nil, false
	}
	if s[0] != '|' || s[1] != '|' || s[l-1] != '|' || s[l-2] != '|' {
		return nil, false
	}
	s = s[2 : l-2]
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
	return "||" + strings.Join(authors, "||") + "||"
}

func parseBlogCategories(s string) ([]string, bool) {
	l := len(s)
	if l < 4 {
		return nil, false
	}
	if s[0] != '|' || s[1] != '|' || s[l-1] != '|' || s[l-2] != '|' {
		return nil, false
	}
	s = s[2 : l-2]
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
	return "||" + strings.Join(cats, "||") + "||"
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
  Id int64 `json:"id"`
	BlogId      int64  `json:"blogId"`
	Description string `json:"description"`
	// Second-precision
	Timestamp int64 `json:"timestamp"`
	// Seconds
	TzOffset int    `json:"tzOffset"`
	PrevHash string `json:"prevHash"`
	Hash     string `json:"hash"`
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
	Edit       Edit
	ContentUrl string
}

func NewBlogPostPageData() BlogPostPageData {
	return BlogPostPageData{
		Blog:       Blog{},
		Edit:       Edit{},
		ContentUrl: "",
	}
}
