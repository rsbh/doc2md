package transformer

import (
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

//JSONToMD convert json to markdown
func JSONToMD(json []M) string {
	var content []string
	for _, j := range json {
		keys := make([]string, 0, len(j))
		for k := range j {
			keys = append(keys, k)
		}
		key := keys[0]
		i, ok := headings[key]
		if ok {
			s := getHeader(j[key], i)
			content = append(content, s, "\n")
		} else if key == "p" {
			s := j[key]
			content = append(content, s, "\n")
		}
	}
	return strings.Join(content, "")
}
