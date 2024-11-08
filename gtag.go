package gtag

import (
	"fmt"
	"io"
	re "regexp"
	"strings"
)

type Node interface { 
    Render(io.Writer) error
}
type Tag struct {
    name string
    void, root bool
    children []Node
    attributes map[string]string
}
type Literal struct {
    raw string 
}
type Text struct {
    text string
}

func (l *Literal) Render(w io.Writer) error { 
    _, err := w.Write([]byte(l.raw))
    return err
}

var escape = re.MustCompile(`[<>&"']`)
var lookup = map[string]string { "<": "&lt;", ">": "&gt;", "&": "%amp;", "\"": "&quot;", "'": "&#39;", }
func replace(a string) string {
    return escape.ReplaceAllStringFunc(a, func(s string) string { return lookup[s] })
}
func (t *Text) Render(w io.Writer) error {
    _, err := w.Write([]byte(replace(t.text)))
    return err
}

func (t *Tag) Render(w io.Writer) error {

    if t.root {
        _, _ = fmt.Fprintf(w,  "<!DOCTYPE html>\n")
    }
    _, _ = fmt.Fprintf(w,  "<%s", t.name)

    if len(t.attributes) > 0 {
        for key := range t.attributes {
            _, _ = fmt.Fprintf(w, " %s=\"%s\"", key, replace(t.attributes[key]))
        }
    }
    if t.void { 
        _, _ = fmt.Fprintf(w, "/")
    }

    _, _ = fmt.Fprintf(w, ">")

    if t.void {
        return nil
    }
    for _, child := range t.children {
        err := child.Render(w)
        if err != nil {
            return err
        }
    }
    _, _ = fmt.Fprintf(w, "</%s>", t.name)
    return nil
}

func (t *Tag) String() string {
    b := new(strings.Builder)
    t.Render(b)
    return b.String()
}

func (t *Tag) Append(comp ...Node) *Tag {
    t.children = append(t.children, comp...)
    return t
}

func (t *Tag) Tag(name string) *Tag {
    n := New(name)
    t.Append(n)
    return n
}

func (t *Tag) VoidTag(name string) *Tag {
    n := NewVoid(name)
    t.Append(n)
    return n
}

func (t *Tag) Asis(raw string)      *Tag {  return t.Append(&Literal{raw: raw}) }
func (t *Tag) Text(content string)  *Tag {  return t.Append(&Text{text: content}) }

func (t *Tag) SetAttr(  key, value  string) *Tag {  t.attributes[key] = value;                       return t }
func (t *Tag) Class(    classes  ...string) *Tag {  t.SetAttr("class",  strings.Join(classes, " ")); return t }
func (t *Tag) Style(    style       string) *Tag {  t.SetAttr("style",  style);                      return t }
func (t *Tag) Id(       id          string) *Tag {  t.SetAttr("id",     id);                         return t }
func (t *Tag) Href(     link        string) *Tag {  t.SetAttr("href",   link);                       return t }

func newTag(name string, void bool) *Tag {
    n := new(Tag)
    n.attributes = make(map[string]string)
    n.name = name
    n.void = void
    return n
}
func Doc() *Tag {
    root := new(Tag).Tag("html") 
    root.root = true
    return root
}
func Div()                  *Tag { return newTag("div", false) }
func New(name string)       *Tag { return newTag(name, false) }
func NewVoid(name string)   *Tag { return newTag(name, true) }

func (t *Tag) Head()    *Tag { return t.Tag("head") }
func (t *Tag) Body()    *Tag { return t.Tag("body") }
func (t *Tag) P()       *Tag { return t.Tag("p") }
func (t *Tag) Div()     *Tag { return t.Tag("div") }
