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
	return fmt.Sprintf("![%v](%v \"%v\")", i.Title, i.Source, i.Title)
}

func getTableTag(t Table) string {
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

func getCodeTag(c CodeBlock) string {
	return "```" + c.Lang + "\n" + strings.Join(c.Content, "\n") + "```"
}

func convertOrderedList(list []string) string {
	var str string
	for i, li := range list {
		str = str + fmt.Sprintf("\n %v. %v", i+1, li)
	}
	return str
}

func convertUnorderedList(list []string) string {
	var str string
	for _, li := range list {
		str = str + fmt.Sprintf("\n -  %v", li)
	}
	return str
}

//JSONToMD convert json to markdown
func JSONToMD(json []Tag) string {
	var content []string
	for _, j := range json {
		key := j.Name
		i, ok := headings[key]
		if ok {
			s := getHeader(j.Content.Text, i)
			content = append(content, s, "\n")
		} else if key == "p" {
			s := j.Content.Text
			content = append(content, s, "\n")
		} else if key == "img" {
			i := getImageTag(j.Content.Image)
			content = append(content, i, "\n")
		} else if key == "table" {
			t := getTableTag(j.Content.Table)
			content = append(content, t, "\n")
		} else if key == "code" {
			c := getCodeTag(j.Content.CodeBlock)
			content = append(content, c, "\n")
		} else if key == "ol" {
			list := convertOrderedList(j.Content.List)
			content = append(content, list, "\n")
		} else if key == "ul" {
			list := convertUnorderedList(j.Content.List)
			content = append(content, list, "\n")
		} else {

		}
	}
	return strings.Join(content, "")
}
