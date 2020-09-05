package transformer

import (
	"fmt"
	"strings"
)

var headings = map[string]int{
	"h1": 1,
	"h2": 2,
	"h3": 3,
	"h4": 4,
	"h5": 5,
	"h6": 6,
}

func getHeader(text string, repeat int) string {
	h := strings.Repeat("#", repeat)
	return h + " " + text
}

func getImageTag(i ImageObject) string {
	return fmt.Sprintf("![%v](%v %v)", i.Title, i.Source, i.Title)
}

func getList(i ImageObject) string {
	return fmt.Sprintf("![%v](%v %v)", i.Title, i.Source, i.Title)
}

func getTable(t Table) string {
	header := " | "
	spaces := " | "
	body := ""
	for _, h := range t.Header {
		header = header + h + " | "
		spaces = spaces + "----" + " | "
	}

	for _, r := range t.Rows {
		body = body + " | " + strings.Join(r, " | ") + " | \n"
	}
	return strings.Join([]string{header, spaces, body}, "\n")
}

//JSONToMD convert json to markdown
func JSONToMD(json []TagContent) string {
	var content []string
	for _, j := range json {
		keys := make([]string, 0, len(j))
		for k := range j {
			keys = append(keys, k)
		}
		key := keys[0]
		i, ok := headings[key]
		if ok {
			s := getHeader(j[key].Text, i)
			content = append(content, s, "\n")
		} else if key == "p" {
			s := j[key].Text
			content = append(content, s, "\n")
		} else if key == "img" {
			i := getImageTag(j[key].Image)
			content = append(content, i, "\n")
		} else if key == "table" {
			t := getTable(j[key].Table)
			content = append(content, t, "\n")
		} else {
			fmt.Println(key)
		}
	}
	return strings.Join(content, "")
}
