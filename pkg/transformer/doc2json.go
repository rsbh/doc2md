package transformer

import (
	"fmt"
	"strings"

	"google.golang.org/api/docs/v1"
)

var tags = map[string]string{
	"HEADING_1":   "h1",
	"HEADING_2":   "h2",
	"HEADING_3":   "h3",
	"HEADING_4":   "h4",
	"HEADING_5":   "h5",
	"NORMAL_TEXT": "p",
	"SUBTITLE":    "blockquote",
}

func getText(e *docs.ParagraphElement) string {
	text := e.TextRun.Content
	isEmptyString := len(text) == 0
	if e.TextRun.TextStyle.Italic && !isEmptyString {
		text = fmt.Sprintf("_%v_", text)
	}

	if e.TextRun.TextStyle.Bold && !isEmptyString {
		text = fmt.Sprintf("**%v**", text)
	}

	if e.TextRun.TextStyle.Strikethrough && !isEmptyString {
		text = fmt.Sprintf("~~%v~~", text)
	}

	if e.TextRun.TextStyle.Link != nil && !isEmptyString {
		text = fmt.Sprintf("[%v](%v)", text, e.TextRun.TextStyle.Link.Url)
	}
	return text
}

type M map[string]string

func all(vs []M, f func(M) bool) bool {
	for _, v := range vs {
		if !f(v) {
			return false
		}
	}
	return true
}

func getTagContent(p *docs.Paragraph, tag string) []M {
	var tagContent []M
	for _, e := range p.Elements {

		tr := e.TextRun
		if tr != nil {
			// headingID := p.ParagraphStyle.HeadingId
			text := getText(e)
			x := M{tag: text}
			tagContent = append(tagContent, x)
		}
	}

	if all(tagContent, func(m M) bool { return m[tag] != "" }) {
		var a []string
		for _, tc := range tagContent {
			a = append(a, tc[tag])
		}
		s := strings.Join(a, " ")
		s = strings.ReplaceAll(s, " .", ".")
		s = strings.ReplaceAll(s, " ,", ",")
		return []M{M{tag: s}}
	} else {
		return tagContent
	}
}

func getParagraph(p *docs.Paragraph) []M {
	var tc []M
	t := p.ParagraphStyle.NamedStyleType
	tag := tags[t]
	if tag != "" {
		tc = getTagContent(p, tag)
	}
	return tc
}

// DocToJSON Convert docs api response to json
func DocToJSON(b *docs.Body) []M {
	var content []M
	for _, s := range b.Content {
		if s.Paragraph != nil {
			c := getParagraph(s.Paragraph)
			content = append(content, c...)
		}
	}
	return content
}
