package transformer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/dlclark/regexp2"
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

func cleanText(text string, ignoreLineBreak bool) string {
	s := text
	if ignoreLineBreak {
		s = strings.ReplaceAll(s, "\n", "<br/>")
	} else {
		s = strings.ReplaceAll(s, "\n", "")
	}
	return strings.TrimSpace(s)
}

func getText(e *docs.ParagraphElement, ignoreLineBreak bool, isHeader bool) string {
	text := e.TextRun.Content

	re1 := regexp2.MustCompile(`<(?![<br/>])`, 0)
	if isMatch, _ := re1.MatchString(text); isMatch {
		text, _ = re1.Replace(text, "&lt;", -1, -1)
	}

	re2 := regexp2.MustCompile(`/>`, 0)
	if isMatch, _ := re2.MatchString(text); isMatch {
		text, _ = re2.Replace(text, "&gt;", -1, -1)
	}

	isEmptyString := len(text) == 0
	text = cleanText(text, ignoreLineBreak)

	if e.TextRun.TextStyle.Italic && !isEmptyString {
		text = fmt.Sprintf("_%v_", text)
	}

	if e.TextRun.TextStyle.Bold && !isHeader && !isEmptyString {
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
	Table
	CodeBlock
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
	isHeader, _ := regexp.MatchString("^h[1-6]$", tag)
	fmt.Println(isHeader)
	for _, e := range p.Elements {

		tr := e.TextRun
		if e.InlineObjectElement != nil {
			i := getImage(ios, e.InlineObjectElement.InlineObjectId)
			x := TagContent{"img": {"", i, Table{}, CodeBlock{}}}
			tagContent = append(tagContent, x)
		} else if tr != nil && tr.Content != "\n" {
			// headingID := p.ParagraphStyle.HeadingId
			text := getText(e, true, isHeader)
			x := TagContent{tag: {text, ImageObject{}, Table{}, CodeBlock{}}}
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
		return []TagContent{{tag: {s, ImageObject{}, Table{}, CodeBlock{}}}}
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
		s = S{t, i, Table{}, CodeBlock{}}
	} else {
		t := getText(e, true, false)
		s = S{t, ImageObject{}, Table{}, CodeBlock{}}
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
		tc = append(tc, TagContent{listTag: {bc, ImageObject{}, Table{}, CodeBlock{}}})
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
	HeadingID string        `json:"headingId"`
	Text      string        `json:"text"`
	Indent    float64       `json:"indent"`
	Items     []*TocHeading `json:"items"`
}

func getToc(data *docs.TableOfContents) []*TocHeading {
	var toc []*TocHeading
	var cur *TocHeading
	for _, c := range data.Content {
		text := c.Paragraph.Elements[0].TextRun.Content
		headingID := c.Paragraph.Elements[0].TextRun.TextStyle.Link.HeadingId
		indent := c.Paragraph.ParagraphStyle.IndentStart.Magnitude
		if indent == 0 {
			t := TocHeading{headingID, text, indent, nil}
			toc = append(toc, &t)
			cur = &t
		} else {
			sub := TocHeading{headingID, text, indent, nil}
			cur.Items = append(cur.Items, &sub)
		}
	}
	return toc
}

func getTextFromParagraph(p *docs.Paragraph, ignoreLineBreak bool) string {
	var sa []string
	for _, e := range p.Elements {
		if e.TextRun != nil {
			a := getText(e, ignoreLineBreak, false)
			sa = append(sa, a)
		} else {
			sa = append(sa, "")
		}
	}
	s := strings.Join(sa, " ")
	s = strings.ReplaceAll(s, " .", ".")
	s = strings.ReplaceAll(s, " ,", ",")
	return s
}

type Table struct {
	Header []string
	Rows   [][]string
}

func getTableCellContent(content []*docs.StructuralElement) string {
	var sa []string
	if len(content) == 0 {
		return ""
	} else {
		for _, i := range content {
			a := getTextFromParagraph(i.Paragraph, false)
			sa = append(sa, a)
		}
		s := strings.Join(sa, " ")
		return s
	}
}

type CodeBlock struct {
	Lang    string
	Content string
}

func getCodeBlock(cell *docs.TableCell) CodeBlock {
	var codeArr []string
	for _, c := range cell.Content {
		if c != nil && c.Paragraph != nil {
			for _, e := range c.Paragraph.Elements {
				codeArr = append(codeArr, e.TextRun.Content)
			}
		}
	}
	code := strings.Join(codeArr, "")
	return CodeBlock{"sh", code}
}

func getTable(t *docs.Table, supportCodeBlock bool) TagContent {
	if supportCodeBlock && t.Rows == 1 && t.Columns == 1 {
		cell := t.TableRows[0].TableCells[0]
		cb := getCodeBlock(cell)
		return TagContent{"code": {"", ImageObject{}, Table{}, cb}}
	} else {
		thead, tbody := t.TableRows[0], t.TableRows[1:]
		var header []string
		var rows [][]string
		for _, t := range thead.TableCells {
			str := getTableCellContent(t.Content)
			header = append(header, str)
		}
		for _, b := range tbody {
			var temp []string
			for _, t := range b.TableCells {
				str := getTableCellContent(t.Content)
				temp = append(temp, str)
			}
			rows = append(rows, temp)
		}
		return TagContent{"table": {"", ImageObject{}, Table{header, rows}, CodeBlock{}}}
	}
}

func checkInToc(headingID string, toc []*TocHeading) (bool, string) {
	var isInToc bool
	var title string
	if headingID == "" {
		isInToc = false
	} else {
		for _, c := range toc {
			if c.HeadingID == headingID {
				isInToc = true
				title = c.Text
				break
			} else {
				isInToc, title = checkInToc(headingID, c.Items)
			}
		}
	}
	return isInToc, title
}

type Page struct {
	Title    string
	Contents []TagContent
}

// DocToJSON Convert docs api response to json
func DocToJSON(doc *docs.Document, supportCodeBlock bool, breakPages bool) ([]Page, []*TocHeading) {
	b := doc.Body
	var toc []*TocHeading
	ios := doc.InlineObjects
	lists := doc.Lists
	var content []TagContent
	var pages []Page
	var prevTitle string
	for _, s := range b.Content {
		if s.TableOfContents != nil {
			toc = getToc(s.TableOfContents)
		} else if s.Paragraph != nil {
			headingID := s.Paragraph.ParagraphStyle.HeadingId
			headingTag := s.Paragraph.ParagraphStyle.NamedStyleType
			isInToc, title := checkInToc(headingID, toc)
			if isInToc && breakPages && headingTag == "HEADING_2" {
				if prevTitle != "" {
					page := Page{prevTitle, content}
					pages = append(pages, page)
				}
				content = []TagContent{}
				prevTitle = title
			}
			c := getParagraph(s.Paragraph, ios, lists)
			content = append(content, c...)
		} else if s.Table != nil && len(s.Table.TableRows) > 0 {
			tc := getTable(s.Table, supportCodeBlock)
			content = append(content, tc)
		}
	}
	if len(pages) == 0 {
		page := Page{"index", content}
		pages = append(pages, page)
	}
	return pages, toc
}
