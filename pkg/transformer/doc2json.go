package transformer

import (
	"fmt"
	"io/ioutil"
	"path"
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

func replaceTags(text string) string {
	re1 := regexp2.MustCompile(`<(?![<br/>])`, 0)
	re2 := regexp2.MustCompile(`>`, 0)
	isMatch1, _ := re1.MatchString(text)
	isMatch2, _ := re2.MatchString(text)
	if isMatch1 && isMatch2 {
		text, _ = re1.Replace(text, "&lt;", -1, -1)
		text, _ = re2.Replace(text, "&gt;", -1, -1)
	}
	return text
}

func getText(e *docs.ParagraphElement, ignoreLineBreak bool, isHeader bool) string {
	text := e.TextRun.Content
	text = replaceTags(text)

	isEmptyString := len(text) == 0

	if isEmptyString {
		return ""
	}

	text = cleanText(text, ignoreLineBreak)

	if e.TextRun.TextStyle.Italic {
		text = fmt.Sprintf("_%v_", text)
	}

	if e.TextRun.TextStyle.Bold && !isHeader {
		text = fmt.Sprintf("**%v**", text)
	}

	if e.TextRun.TextStyle.Strikethrough {
		text = fmt.Sprintf("~~%v~~", text)
	}

	if e.TextRun.TextStyle.Link != nil {
		text = fmt.Sprintf("[%v](%v)", text, e.TextRun.TextStyle.Link.Url)
	}
	if isHeader {
		text = text + "\n"
	}
	return text
}

type TagContent struct {
	Text  string
	Image ImageObject
	Table
	CodeBlock
	List []string
}

type Tag struct {
	Name    string
	Content TagContent
}

func all(vs []Tag, f func(Tag) bool) bool {
	for _, v := range vs {
		if !f(v) {
			return false
		}
	}
	return true
}

func getTagContent(p *docs.Paragraph, tag string, imageFolder string, ios map[string]docs.InlineObject) []Tag {
	var tagContent []Tag
	isHeader, _ := regexp.MatchString("^h[1-6]$", tag)
	for _, e := range p.Elements {

		tr := e.TextRun
		if e.InlineObjectElement != nil {
			i := getImage(ios, imageFolder, e.InlineObjectElement.InlineObjectId)
			x := Tag{"img", TagContent{Image: i}}
			tagContent = append(tagContent, x)
		} else if tr != nil && tr.Content != "\n" {
			// headingID := p.ParagraphStyle.HeadingId
			text := getText(e, true, isHeader)
			x := Tag{tag, TagContent{Text: text}}
			tagContent = append(tagContent, x)
		}
	}

	if all(tagContent, func(m Tag) bool { return m.Content.Text != "" }) {
		var a []string
		for _, tc := range tagContent {
			a = append(a, tc.Content.Text)
		}
		s := joinStrings(a)
		return []Tag{{tag, TagContent{Text: s}}}
	}
	return tagContent
}

func getListType(lists map[string]docs.List, listID string) string {
	tag := "ol"
	if val, ok := lists[listID]; ok {
		nestingLevels := val.ListProperties.NestingLevels
		if len(nestingLevels) > 0 {
			glyphType := nestingLevels[0].GlyphType
			if glyphType == "" {
				tag = "ul"
			}
		}
	}
	return tag
}

func getBulletContents(ios map[string]docs.InlineObject, e *docs.ParagraphElement, imageFolder string) TagContent {
	var s TagContent
	if e.InlineObjectElement != nil {
		i := getImage(ios, imageFolder, e.InlineObjectElement.InlineObjectId)
		t := getImageTag(i)
		s = TagContent{Text: t, Image: i}
	} else {
		t := getText(e, true, false)
		s = TagContent{Text: t}
	}
	return s
}

func getNestedListIndent(level int, listTag string) string {
	indentType := "-"
	if listTag == "ol" {
		indentType = "1."
	}
	indent := strings.Repeat("  ", level)
	return fmt.Sprintf("%v%v ", indent, indentType)
}

func getParagraph(p *docs.Paragraph, imageFolder string, ios map[string]docs.InlineObject, lists map[string]docs.List, prev *docs.Paragraph, contents *[]Tag) []Tag {
	var tc []Tag
	if p.Bullet != nil {
		listID := p.Bullet.ListId
		var prevID string

		if prev.Bullet != nil {
			prevID = prev.Bullet.ListId
		}

		listTag := getListType(lists, listID)
		var bulletContents []string
		for _, e := range p.Elements {
			s := getBulletContents(ios, e, imageFolder)
			if s.Text != "" {
				bulletContents = append(bulletContents, s.Text)
			}
		}
		bc := joinStrings(bulletContents)
		if listID == prevID {
			c := *contents
			last := c[len(c)-1].Content.List
			nestingLevel := p.Bullet.NestingLevel
			lastIndex := len(last) - 1

			if nestingLevel != 0 && lastIndex > 0 {
				indent := getNestedListIndent(lastIndex, listTag)
				last[lastIndex] += fmt.Sprintf("\n%v %v", indent, bc)
			} else {
				last = append(last, bc)
				c[len(c)-1] = Tag{listTag, TagContent{List: last}}
				*contents = c
			}
		} else {
			tc = append(tc, Tag{listTag, TagContent{List: []string{bc}}})
		}
	} else {
		t := p.ParagraphStyle.NamedStyleType
		tag := tags[t]
		if tag != "" {
			tc = getTagContent(p, tag, imageFolder, ios)
		}
	}
	return tc
}

type ImageObject struct {
	Source      string
	Title       string
	Description string
}

func getImage(ios map[string]docs.InlineObject, imageFolder string, objectID string) ImageObject {
	var image ImageObject
	eo := ios[objectID].InlineObjectProperties.EmbeddedObject

	if eo != nil && eo.ImageProperties != nil {
		src := eo.ImageProperties.ContentUri
		name, content := ReplaceImage(src)
		imgPath := path.Join(imageFolder, name)
		ioutil.WriteFile(imgPath, content, 0644)
		imgLink := path.Join("images", name)

		image = ImageObject{imgLink, eo.Title, eo.Description}
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

func joinStrings(sa []string) string {
	s := strings.Join(sa, " ")
	s = strings.ReplaceAll(s, " .", ".")
	s = strings.ReplaceAll(s, " ,", ",")
	return s
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
	s := joinStrings(sa)
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
	}
	for _, i := range content {
		a := getTextFromParagraph(i.Paragraph, false)
		sa = append(sa, a)
	}
	s := strings.Join(sa, " ")
	return s
}

type CodeBlock struct {
	Lang    string
	Content []string
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
	codeStr := strings.Join(codeArr, "")
	code := strings.Split(codeStr, "\u000b")
	return CodeBlock{"sh", code}
}

func getTable(t *docs.Table, supportCodeBlock bool) Tag {
	if supportCodeBlock && t.Rows == 1 && t.Columns == 1 {
		cell := t.TableRows[0].TableCells[0]
		cb := getCodeBlock(cell)
		return Tag{"code", TagContent{CodeBlock: cb}}
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
		return Tag{"table", TagContent{Table: Table{header, rows}}}
	}
}

func checkInToc(headingID string, toc []*TocHeading) (bool, string) {
	var isInToc bool
	var title string
	if headingID == "" {
		isInToc = false
	} else {
		for _, c := range toc {
			for _, j := range c.Items {
				if j.HeadingID == headingID {
					isInToc = true
					title = j.Text
					break
				}
			}
		}
	}
	return isInToc, title
}

type Page struct {
	Title    string
	Contents []Tag
}

// DocToJSON Convert docs api response to json
func DocToJSON(doc *docs.Document, imageFolder string, supportCodeBlock bool, breakPages bool) ([]Page, []*TocHeading) {
	b := doc.Body
	var toc []*TocHeading
	ios := doc.InlineObjects
	lists := doc.Lists
	var content []Tag
	var pages []Page
	var prevTitle string
	for i, s := range b.Content {
		if s.TableOfContents != nil {
			toc = getToc(s.TableOfContents)
		} else if s.Paragraph != nil {
			if breakPages {
				headingID := s.Paragraph.ParagraphStyle.HeadingId
				headingTag := s.Paragraph.ParagraphStyle.NamedStyleType

				isInToc, title := checkInToc(headingID, toc)
				if (isInToc && headingTag == "HEADING_2") || headingTag == "HEADING_1" {
					if prevTitle != "" {
						page := Page{prevTitle, content}
						pages = append(pages, page)
					}
					content = []Tag{}
					prevTitle = title
				}
			}

			prev := b.Content[i-1].Paragraph
			c := getParagraph(s.Paragraph, imageFolder, ios, lists, prev, &content)
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
