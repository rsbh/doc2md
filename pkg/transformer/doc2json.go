package transformer

import (
	"fmt"

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

func getParagraph(p *docs.Paragraph) {
	t := p.ParagraphStyle.NamedStyleType
	tag := tags[t]

	if tag != "" {
		for _, e := range p.Elements {

			tr := e.TextRun
			if tr != nil {
				headingID := p.ParagraphStyle.HeadingId
				text := getText(e)
				fmt.Println(headingID, text)
			}
		}

	}
}

// DocToJSON Convert docs api response to json
func DocToJSON(b *docs.Body) {
	for _, s := range b.Content {
		if s.Paragraph != nil {
			getParagraph(s.Paragraph)
		}
	}
}
