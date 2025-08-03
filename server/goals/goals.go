package goals

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"strings"
	"unicode"

	utils "github.com/johnietre/utils/go"
)

var (
	goalsDb *sql.DB
)

func ParseMarkdown(s string) (*goalItem, error) {
	return ParseMarkdownReader(strings.NewReader(s))
}

func ParseMarkdownReader(r io.Reader) (*goalItem, error) {
	br := bufio.NewReader(r)
	items := []*goalItem{newRootGoalItem()}
	depth := 0
	for {
		line, err := br.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		indent := 0
		for _, r := range line {
			if r != '\t' {
				break
			}
			indent++
		}
		if false {
			line = strings.TrimRightFunc(line, unicode.IsSpace)
		}
		line = strings.TrimSpace(line)
		parts := splitSpace(line)

		listType := giTypeNone
		if len(parts) != 0 {
			if parts[0] == "-" {
				listType = giTypeUListItem
			}
			if listType != giTypeNone {
				parts = parts[1:]
			}
		}

		checkbox := (*bool)(nil)
		if len(parts) != 0 {
			switch parts[0] {
			case "[ ]":
				checkbox = utils.NewT(false)
				parts = parts[1:]
			case "[x]", "[X]":
				checkbox = utils.NewT(true)
				parts = parts[1:]
			}
		}

		headerCount, validHeader := 0, false
		if len(parts) != 0 {
			for _, r := range line {
				if r != '#' {
					if r == ' ' {
						validHeader = true
					}
				}
				headerCount++
				if headerCount > 6 {
					break
				}
			}
			if validHeader {
				parts = parts[1:]
			}
		}

		l := len(items)
		if listType != giTypeNone {
			l++
		}

		gi := &goalItem{
			Id:        0,
			What:      strings.Join(parts, " "),
			Parent:    &items[l-1].Id,
			Completed: checkbox,
		}
		items[l-1].children = append(items[l-1].children, gi)
	}

	return items[0], nil
}

type goalItem struct {
	Id        int64
	What      string
	Type      giType
	Parent    *int64
	Completed *bool
	Hidden    bool

	children []*goalItem
}

func newRootGoalItem() *goalItem {
	return &goalItem{
		Id: 0,
	}
}

// Always returns a goalItem with a Type of giTypeNone and Id of 0 as root.
// Any items without parents are the children of this root item.
func consolidateGoalItems(gis []*goalItem) *goalItem {
	idMap := make(map[int64]*goalItem)
	for _, gi := range gis {
		idMap[gi.Id] = gi
	}
	orphans := []*goalItem{}
	if gi := idMap[0]; gi != nil {
		orphans = append(orphans, gi)
	}
	idMap[0] = &goalItem{}
	for _, gi := range gis {
		parent := idMap[utils.ValOr(gi.Parent, 0)]
		if parent == nil {
			orphans = append(orphans, gi)
			continue
		}
		parent.children = append(parent.children, gi)
	}
	return idMap[0]
}

func (gi *goalItem) generateHtml(showHidden bool) string {
	if !showHidden && gi.Hidden {
		return ""
	}
	ret, tag := "", ""
	switch gi.Type {
	case giTypeHeader1, giTypeHeader2, giTypeHeader3,
		giTypeHeader4, giTypeHeader5, giTypeHeader6:
		tag = fmt.Sprintf("h%d", gi.Type-giTypeHeader1+1)
		ret = `<` + tag + `>`
	case giTypeUListItem, giTypeOListItem:
		ret += `<li>`
		tag = `li`
	case giTypeCheckbox:
		if utils.ValOr(gi.Completed, false) {
			ret += `<input type="checkbox" checked disabled>`
		} else {
			ret += `<input type="checkbox" disabled>`
		}
	}
	ret += gi.What
	listCloseTag := ""
	for _, child := range gi.children {
		if child.Type == giTypeOListItem {
			if listCloseTag != `</ol>` {
				ret += listCloseTag
				ret += `<ol>`
				listCloseTag = `</ol>`
			}
		} else if child.Type == giTypeUListItem {
			if listCloseTag != `</ul>` {
				ret += listCloseTag
				ret += `<ul>`
				listCloseTag = `</ul>`
			}
		} else {
			ret += listCloseTag
			listCloseTag = ""
		}
		ret += child.generateHtml(showHidden)
	}
	ret += listCloseTag
	if tag != "" {
		ret += `</` + tag + `>`
	}
	return ret
}

type giType int

const (
	giTypeNone      giType = 0
	giTypeUListItem giType = 1
	giTypeOListItem giType = 2
	giTypeCheckbox  giType = 3
	giTypeText      giType = 4
	giTypeHeader1   giType = 5
	giTypeHeader2   giType = 6
	giTypeHeader3   giType = 7
	giTypeHeader4   giType = 8
	giTypeHeader5   giType = 9
	giTypeHeader6   giType = 10
)

func splitSpace(s string) []string {
	return utils.FilterSliceInPlace(
		strings.Split(s, " "),
		func(s string) bool { return s != "" },
	)
}
