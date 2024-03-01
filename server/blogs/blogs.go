package blogs

import (
	"database/sql"

	//utils "github.com/johnietre/utils/go"
	_ "github.com/mattn/go-sqlite3"
)

var (
  //BlogsData = utils.NewAValue[BlogData](BlogData{})
  blogsDb *sql.DB
)

func InitBlogs(blogsDir, dbPath string) error {
  if err := openBlogsDB(dbPath); err != nil {
    return err
  }
  return LoadBlogData()
}

func openBlogsDB(dbPath string) (err error) {
  blogsDb, err = sql.Open("sqlite3", dbPath)
  return
}

func LoadBlogData() error {
  /*
  blogRows, err := appsDb.Query(`SELECT * FROM blogs`)
  if err != nil {
    return err
  }
  defer blogRows.Close()
  data := BlogData{}
  for blogRows.Next() {
    blog := Blog{}
    if err := blogRows.Scan(
      &blog.Id, &blog.Title, &blog.Timestamp, &blog.TzOffset,
    ); err != nil {
      return err
    }
    data.Blogs = append(data.Blogs, blog)
  }

  catRows, err := appsDb.Query(`SELECT * FROM categories`)
  if err != nil {
    return err
  }
  defer catRows.Close()
  for catRows.Next() {
    cat := ""
    if err := catRows.Scan(&cat); err != nil {
      return err
    }
    data.Categories = append(data.Categories, cat)
  }

  BlogsData.Store(data)
  */
  return nil
}

func blogNotFound() {
}

type BlogsPageData struct {
  Categories []string
  Blogs []Blog
}

func NewBlogsPageData() BlogsPageData {
  return BlogsPageData{
    Categories: []string{"One", "Two", "Three"},
    Blogs: []Blog{
      {Id: 1, Title: "Yes"},
      {Id: 2, Title: "No"},
      {Id: 3, Title: "Maybe"},
    },
  }
}

type Blog struct {
  Id uint64 `json:"id"`
  Title string `json:"title"`
  Timestamp int64 `json:"timestamp"`
  TzOffset int `json:"tzOffset"`
}

type BlogIndex struct {
  Categories []string
}
