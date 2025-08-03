package blogs

import (
	"path/filepath"
	"strings"
	"time"
)

func GenerateExamples(dir string) error {
	dbPath := filepath.Join(dir, "blogs.db")
	if err := InitBlogs(dir, dbPath); err != nil {
		return err
	}
	for i, blog := range exampleBlogs {
		err := AddBlog(&blog, strings.NewReader(exampleBlogContents[i]))
		if err != nil {
			return err
		}
	}
	for i, edit := range exampleBlogEdits {
		err := EditBlog(&edit, strings.NewReader(exampleBlogEditContents[i]))
		if err != nil {
			return err
		}
	}
	panic("unimplemented")
}

var (
	exampleBlogs = []Blog{
		{
			Title:      "Blog 1",
			Authors:    []string{"Author 1"},
			Categories: []string{"Category 1", "Category 2"},
			Timestamp:  time.Now().Add(-time.Hour * 24 * 10).Unix(),
			TzOffset:   -6 * 60 * 60,
		},
		{
			Title:      "Blog 2",
			Authors:    []string{"Author 2"},
			Categories: []string{"Category 2", "Category 3"},
			Timestamp:  time.Now().Add(-time.Hour * 24 * 5).Unix(),
			TzOffset:   -6 * 60 * 60,
		},
		{
			Title:      "Blog 3",
			Authors:    []string{"Author 1", "Author 2", "Author 3"},
			Categories: []string{"Category 1", "Category 3"},
			Timestamp:  time.Now().Add(-time.Hour * 12 * 5).Unix(),
			TzOffset:   -6 * 60 * 60,
		},
	}
	exampleBlogContents = []string{
		exampleBlog1Content,
		exampleBlog2Content,
		exampleBlog3Content,
	}
	exampleBlogEdits        = []Edit{}
	exampleBlogEditContents = []string{}
)

const exampleBlog1Content = `
`

const exampleBlog2Content = `
`

const exampleBlog3Content = `
`
