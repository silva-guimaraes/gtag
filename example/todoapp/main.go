package main

import (
	_ "embed"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/silva-guimaraes/gtag"
)

const (
	loremParagraph = "Lorem ipsum dolor sit amet, qui minim labore adipisicing minim sint cillum" +
		"sint consectetur cupidatat."

	h2Style = `display: block;
    font-size: 1.5em;
    margin-block-start: 0.83em;
    margin-block-end: 0.83em;
    margin-inline-start: 0px;
    margin-inline-end: 0px;
    font-weight: bold;`

	h3Style = `display: block;
    font-size: 1.17em;
    margin-block-start: 1em;
    margin-block-end: 1em;
    margin-inline-start: 0px;
    margin-inline-end: 0px;
    font-weight: bold;`
)

//go:embed style.css
var animation string

const (
	noAnimationMode int = 1 << iota
	editMode            = 1 << iota
)

func card(n *note, flags int) *gtag.Tag {

    id := n.id
	Id := fmt.Sprintf("#note%d", id)
	formId := fmt.Sprintf("#form%d", id)

	edit := (flags & editMode) > 0
	noAnimation := (flags & noAnimationMode) > 0

	card := gtag.Div().Id(Id[1:]).SetAttr("hx-target", "this")
	{
		if edit {
			d := card.Tag("form").
				Style("display: flex; flex-direction: row; justify-content: space-between").
				Id(formId).
				SetAttr("hx-post", fmt.Sprintf("/edit/%d", id))
			{
				d1 := d.Tag("div")
				{
					d1.VoidTag("input").SetAttr("value", n.title).SetAttr("name", "title").Style(h2Style)

					d1.VoidTag("input").SetAttr("value", n.summary).SetAttr("name", "summary").Style(h3Style)
				}
				buttons := d.Div()
				{
					buttons.Tag("button").
						Class("edit-delete").
						SetAttr("hx-post", fmt.Sprintf("/edit/%d", id)).
						SetAttr("hx-swap", "outerHTML").
                        SetAttr("tabindex", "-1").
						Asis("✅")
					buttons.Tag("button").
						Class("edit-delete").
						SetAttr("hx-get", fmt.Sprintf("/note/%d", id)).
						SetAttr("hx-swap", "outerHTML").
                        SetAttr("tabindex", "-1").
						Asis("❌")
					buttons.Tag("button").
						Class("edit-delete").
						SetAttr("hx-delete", fmt.Sprintf("/edit/%d", id)).
						SetAttr("hx-swap", "delete transition:true").
                        SetAttr("tabindex", "-1").
						Asis("🗑️")
				}
			}
			card.Tag("textarea").Text(n.body).SetAttr("form", formId).SetAttr("name", "body")
		} else {
			d := card.Div().Style("display: flex; flex-direction: row; justify-content: space-between")
			{
				d1 := d.Div()
				{
					d1.Tag("h2").Text(n.title).Style("cursor: pointer")
					d1.Tag("h3").Text(n.summary).Style("cursor: pointer")
				}
				buttons := d.Div()
				{
					buttons.Tag("button").
						Class("edit-delete").
						SetAttr("hx-get", fmt.Sprintf("/edit/%d", id)).
						SetAttr("hx-swap", "outerHTML").
                        SetAttr("tabindex", "-1").
						Asis("✏️")
					buttons.Tag("button").
						Class("edit-delete").
						SetAttr("hx-delete", fmt.Sprintf("/edit/%d", id)).
						SetAttr("hx-swap", "delete transition:true").
                        SetAttr("tabindex", "-1").
						Asis("🗑️")
				}
			}
			card.P().Text(n.body)
		}

	}
	if noAnimation {
		card.Class("box-shadow", "note", "no-animation")
	} else {
		card.Class("box-shadow", "note")
	}
	return card
}

func index(notes []*note) *gtag.Tag {

	html := gtag.Doc()
	{
		head := html.Head()
		{
			head.Tag("title").Text("foobar")
			head.Asis(`<script src="https://unpkg.com/htmx.org@2.0.2"></script>`)
			head.Asis(`<script src="https://unpkg.com/hyperscript.org@0.9.13"></script>`)
			head.Tag("style").Asis(animation)
		}
		html.Body().Style("margin: 100px 0 200px 0")
		{

			content := html.Div().Style("max-width: 700px; margin: auto;")
			{
				content.Tag("h1").Text("todo list:")

				content.Tag("button").
					Class("new-note", "box-shadow").
					SetAttr("hx-post", "/clicked").
					SetAttr("hx-swap", "afterbegin transition:true").
					SetAttr("hx-target", "#list").
					Text("new note")

				list := content.Tag("div").Id("list")
				{
					for _, n := range slices.Backward(notes) {
						list.Append(card(n, noAnimationMode))
					}
				}
			}
		}
	}
	return html
}

type note struct {
    id int
	title, summary, body string
	date                 time.Time
}

func main() {
	notes := []*note{
        {id: 0, title: "homework", summary: "homework due on monday", body: loremParagraph},
        {id: 1, title: "mom's birthday", summary: "have to buy mom a present", body: loremParagraph},
        {id: 2, title: "fix bike", summary: "chain is busted", body: loremParagraph},
	}
    idCounter := len(notes)

    findId := func(id int) (*note) {
        idx := slices.IndexFunc(notes, func(n *note) bool {
            return n.id == id
        })
        if idx < 0 {
            panic(idx)
        }
        return notes[idx]
    }

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		index(notes).Render(w)
	})
	http.HandleFunc("POST /clicked", func(w http.ResponseWriter, r *http.Request) {
        newNote := &note{
            id: idCounter,
            title:   "click to edit...",
            summary: "click to edit...",
            body:    "click to edit...",
        }
        n := newNote
		notes = append(notes, n)
		card(n, 0).Render(w)
        idCounter++
	})
	http.HandleFunc("GET /note/{id}", func(w http.ResponseWriter, r *http.Request) {
		_id := r.PathValue("id")
		if _id == "" {
			panic("id?")
		}
		id, err := strconv.Atoi(_id)
		if err != nil {
			panic(err)
		}
		card(findId(id), noAnimationMode).Render(w)
	})
	http.HandleFunc("/edit/{id}", func(w http.ResponseWriter, r *http.Request) {
		_id := r.PathValue("id")
		if _id == "" {
			panic("id?")
		}
		id, err := strconv.Atoi(_id)
		if err != nil {
			panic(err)
		}
		switch r.Method {
		case "DELETE":
            idx := slices.IndexFunc(notes, func(n *note) bool {
                return n.id == id
            })
            if idx < 0 {
                panic(idx)
            }
			notes = slices.Delete(notes, idx, idx+1)
			w.WriteHeader(200)
		case "GET":
			card(findId(id), editMode|noAnimationMode).Render(w)
		case "POST":
			title := r.FormValue("title")
			if title == "" {
				panic("title?")
			}
			summary := r.FormValue("summary")
			if summary == "" {
				panic("summary?")
			}
			body := r.FormValue("body")
			if body == "" {
				panic("body?")
			}
            n := findId(id)
			n.title = title
			n.summary = summary
			n.body = body
			card(findId(id), noAnimationMode).Render(w)
		default:
			w.WriteHeader(405)
		}
	})

	fmt.Println("listening on http://localhost:2425...")
	err := http.ListenAndServe(":2425", nil)
	if err != nil {
		panic(err)
	}
}
