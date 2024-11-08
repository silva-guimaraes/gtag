package gtag_test

import (
	"testing"
    "strconv"
	"github.com/silva-guimaraes/gtag"
)

func stringEqual(tag *gtag.Tag, expected string, t *testing.T) {
    if r := tag.String(); r != expected {
        t.Fatalf("\n%s\n!=\n%s", expected, r)
    }
}

func TestDoc(t *testing.T) {
    expected := "<!DOCTYPE html>\n<html></html>"
    html := gtag.Doc()

    stringEqual(html, expected, t)
}

func TestComposition(t *testing.T) {
    expected := "<!DOCTYPE html>\n<html><head><title>hello world!</title></head>" +
    "<body><div>hello world!</div></body></html>"

    html := gtag.Doc()
    html.Head().Tag("title").Text("hello world!")
    html.Body().Div().Text("hello world!")

    stringEqual(html, expected, t)
}

func TestRawHtml(t *testing.T) {
    expected := "<!DOCTYPE html>\n<html><head><title>hello world!</title>" +
    "<script src=\"https://unpkg.com/htmx.org@2.0.3\"></script></head></html>"

    html := gtag.Doc()
    head := html.Head(); {
        head.Tag("title").Text("hello world!")
        head.Asis(`<script src="https://unpkg.com/htmx.org@2.0.3"></script>`)
    }

    stringEqual(html, expected, t)
} 

func TestEscapeHTML(t *testing.T) {

    expected := "<div>&lt;script&gt;alert(1)&lt;/script&gt;</div>"
    script := gtag.Div().Text("<script>alert(1)</script>")

    stringEqual(script, expected, t)
}

func TestAppend(t *testing.T) {
    comp := func(name string, age int) *gtag.Tag {
        div := gtag.New("div"); {
            div.P().Text("Name: ").Tag("span").Text(name)
            div.P().Text("Age: ").Tag("span").Text(strconv.Itoa(age))
        }
        return div
    } 
    expected := "<!DOCTYPE html>\n<html><head><title>hello world</title></head><body>" +
    "<div><p>Name: <span>John</span></p><p>Age: <span>34</span></p></div><div><p>Name: " +
    "<span>Jane</span></p><p>Age: <span>25</span></p></div><div><p>Name: <span>Daniel</span>" +
    "</p><p>Age: <span>22</span></p></div></body></html>"
    html := gtag.Doc(); {
        html.Head().Tag("title").Text("hello world")
        body := html.Body(); {
            body.Append(comp("John", 34))
            body.Append(comp("Jane", 25))
            body.Append(comp("Daniel", 22))
        }
    }

    stringEqual(html, expected, t)
}
