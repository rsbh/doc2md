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

type S struct {
	Text  string
	Image ImageObject
}

type TagContent map[string]S

func all(vs []TagContent, f func(TagContent) bool) bool {
	for _, v := range vs {
		if !f(v) {
			return false
		}
	}
	return true
}

func getTagContent(p *docs.Paragraph, tag string, ios map[string]docs.InlineObject) []TagContent {
	var tagContent []TagContent
	for _, e := range p.Elements {

		tr := e.TextRun
		if e.InlineObjectElement != nil {
			i := getImage(ios, e.InlineObjectElement.InlineObjectId)
			x := TagContent{"img": {"", i}}
			tagContent = append(tagContent, x)
		} else if tr != nil && tr.Content != "\n" {
			// headingID := p.ParagraphStyle.HeadingId
			text := getText(e)
			x := TagContent{tag: {text, ImageObject{}}}
			tagContent = append(tagContent, x)
		}
	}

	if all(tagContent, func(m TagContent) bool { return m[tag].Text != "" }) {
		var a []string
		for _, tc := range tagContent {
			a = append(a, tc[tag].Text)
		}
		s := strings.Join(a, " ")
		s = strings.ReplaceAll(s, " .", ".")
		s = strings.ReplaceAll(s, " ,", ",")
		return []TagContent{{tag: {s, ImageObject{}}}}
	} else {
		return tagContent
	}
}

func getParagraph(p *docs.Paragraph, ios map[string]docs.InlineObject) []TagContent {
	var tc []TagContent
	t := p.ParagraphStyle.NamedStyleType
	tag := tags[t]
	if tag != "" {
		tc = getTagContent(p, tag, ios)
	}
	return tc
}

type ImageObject struct {
	Source      string
	Title       string
	Description string
}

func getImage(ios map[string]docs.InlineObject, objectID string) ImageObject {
	var image ImageObject
	eo := ios[objectID].InlineObjectProperties.EmbeddedObject

	if eo != nil && eo.ImageProperties != nil {
		image = ImageObject{eo.ImageProperties.ContentUri, eo.Title, eo.Description}
	}
	return image
}

// DocToJSON Convert docs api response to json
func DocToJSON(doc *docs.Document) []TagContent {
	b := doc.Body
	ios := doc.InlineObjects
	var content []TagContent
	for _, s := range b.Content {
		if s.Paragraph != nil {
			c := getParagraph(s.Paragraph, ios)
			content = append(content, c...)
		}
	}
	return content
}
