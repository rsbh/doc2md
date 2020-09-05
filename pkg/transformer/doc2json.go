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

func getListType(lists map[string]docs.List, listID string) string {
	t := "ol"
	gt := lists[listID].ListProperties.NestingLevels[0].GlyphType
	if gt == "" {
		t = "ul"
	}
	return t
}

func getBulletContents(ios map[string]docs.InlineObject, e *docs.ParagraphElement) S {
	var s S
	if e.InlineObjectElement != nil {
		i := getImage(ios, e.InlineObjectElement.InlineObjectId)
		t := getImageTag(i)
		s = S{t, i}
	} else {
		t := getText(e)
		s = S{t, ImageObject{}}
	}
	return s
}

func getParagraph(p *docs.Paragraph, ios map[string]docs.InlineObject, lists map[string]docs.List) []TagContent {
	var tc []TagContent
	if p.Bullet != nil {
		listID := p.Bullet.ListId
		listTag := getListType(lists, listID)
		var bulletContents []string
		for _, e := range p.Elements {
			s := getBulletContents(ios, e)
			if s.Text != "" {
				bulletContents = append(bulletContents, s.Text)
			}
		}
		bc := strings.Join(bulletContents, " ")
		bc = strings.ReplaceAll(bc, " .", ".")
		bc = strings.ReplaceAll(bc, " ,", ",")
		tc = append(tc, TagContent{listTag: {bc, ImageObject{}}})
	} else {
		t := p.ParagraphStyle.NamedStyleType
		tag := tags[t]
		if tag != "" {
			tc = getTagContent(p, tag, ios)
		}
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

type TocHeading struct {
	HeadingID string       `json:"headingId"`
	Text      string       `json:"text"`
	Indent    float64      `json:"indent"`
	Items     []TocHeading `json:"items"`
}

func getToc(data *docs.TableOfContents) []*TocHeading {
	var toc []*TocHeading
	var cur *TocHeading
	for _, c := range data.Content {
		text := c.Paragraph.Elements[0].TextRun.Content
		headingID := c.Paragraph.Elements[0].TextRun.TextStyle.Link.HeadingId
		indent := c.Paragraph.ParagraphStyle.IndentStart.Magnitude
		if indent == 0 {
			t := TocHeading{headingID, text, indent, []TocHeading{}}
			toc = append(toc, &t)
			cur = &t
		} else {
			sub := TocHeading{headingID, text, indent, []TocHeading{}}
			cur.Items = append(cur.Items, sub)
		}
	}
	return toc
}

// DocToJSON Convert docs api response to json
func DocToJSON(doc *docs.Document) ([]TagContent, []ImageObject, []*TocHeading) {
	b := doc.Body
	var toc []*TocHeading
	ios := doc.InlineObjects
	lists := doc.Lists
	var content []TagContent
	var images []ImageObject
	for _, s := range b.Content {
		if s.TableOfContents != nil {
			toc = getToc(s.TableOfContents)
		} else if s.Paragraph != nil {
			c := getParagraph(s.Paragraph, ios, lists)
			content = append(content, c...)
		}
	}
	for _, c := range content {
		_, ok := c["img"]
		if ok {
			images = append(images, c["img"].Image)
		}
	}
	return content, images, toc
}
