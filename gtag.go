package gtag

import (
	"fmt"
	re "regexp"
	"strings"
)

type Node interface { Render() string }
type Tag struct {
    name string
    void, root bool
    children []Node
    attributes map[string]string
}
type Literal struct { raw string }
type Text struct { text string }

func (l *Literal) Render() string { return l.raw }

var escape = re.MustCompile(`[<>&"']`)
var lookup = map[string]string { "<": "&lt", ">": "&gt", "&": "%amp", "\"": "&quot", "'": "&#39", }
func replace(a string) string {
    return escape.ReplaceAllStringFunc(a, func(s string) string { return lookup[s] })
}

func (t *Text) Render() string {
    return replace(t.text)
}

func (t *Tag) Render() string {
    var ret string
    var renderedAttributes []string

    for key := range t.attributes {
        renderedAttributes = append(renderedAttributes,
            fmt.Sprintf("%s=\"%s\"", key, replace(t.attributes[key])))
    }
    var attributes string = strings.Join(renderedAttributes, " ")
    if len(attributes) > 0 {
        attributes = " " + attributes
    }
    var voidSlash string; if t.void { voidSlash = "/" }

    if t.root {
        ret += (&Literal{raw: "<!DOCTYPE html>"}).Render()
    }

    ret += fmt.Sprintf("<%s%s%s>", voidSlash, t.name, attributes)

    for _, child := range t.children {
        ret += child.Render()
    }

    if !t.void {
        ret += fmt.Sprintf("</%s>", t.name)
    }
    return ret
}

func (t *Tag) Append(comp ...Node) *Tag {
    t.children = append(t.children, comp...)
    return t
}

func (t *Tag) Tag(name string) *Tag {
    n := newTag(name, false)
    t.Append(n)
    return n
}

func (t *Tag) VoidTag(name string) *Tag {
    n := newTag(name, true)
    t.Append(n)
    return n
}

func (t *Tag) Asis(raw string) *Tag {       return t.Append(&Literal{raw: raw}) }
func (t *Tag) Text(content string) *Tag {   return t.Append(&Text{text: content}) }

func (t *Tag) SetAttr(key, value string) *Tag { t.attributes[key] = value;                      return t }
func (t *Tag) Class(classes ...string) *Tag {   t.SetAttr("class", strings.Join(classes, " ")); return t }
func (t *Tag) Style(style string) *Tag {        t.SetAttr("style", style);                      return t }
func (t *Tag) Id(id string) *Tag {              t.SetAttr("id", id);                            return t }

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
func Div() *Tag {
    return newTag("div", false)
}

func (t *Tag) Head()    *Tag { return t.Tag("head") }
func (t *Tag) Body()    *Tag { return t.Tag("body") }
func (t *Tag) P()       *Tag { return t.Tag("p") }
func (t *Tag) Div()     *Tag { return t.Tag("div") }
