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

const cardStyle = "padding: 18px;"+ "margin: 10px;" + "border: 1px solid gray;" + "border-radius: 7px;"

const loremParagraph =  "Lorem ipsum dolor sit amet, qui minim labore adipisicing minim sint cillum sint consectetur cupidatat."

//go:embed style.css
var animation string
func card(id int, n note) *gtag.Tag {

    Id := fmt.Sprintf("#note%d", id)

    card := gtag.Div().Id(Id[1:]).Class("box-shadow", "add-note").Style(cardStyle); {
        d := card.Div().Style("display: flex; flex-direction: row; justify-content: space-between"); {
            d1 := d.Div(); {
                d1.Tag("h2").Text(n.title)
                d1.Tag("h3").Text(n.summary)
            }
            d.Tag("button").
                Style("color: black; height: 40px; width: 40px; font-size: x-large").
                SetAttr("hx-post", "/delete/" + Id[1:]).
                SetAttr("hx-swap", "delete transition:true").
                SetAttr("hx-target", Id).
                Asis("🗑")
        }
        card.P().Text(n.body)

    }
    return card
}

func index(notes []note) *gtag.Tag {
    const buttonStyle = 
        "padding: 10px;" +
        "margin: 10px;" +
        "border-radius: 6px;" +
        "font-size: medium;" +
        "font-family: sans-serif;" +
        "background-color: blue;" +
        "color: white; border-width: 0"

    html := gtag.Doc(); {
        head := html.Head(); {
            head.Tag("title").Text("foobar")
            head.Asis(`<script src="https://unpkg.com/htmx.org@2.0.2"></script>`)
            head.Tag("style").Asis(animation)
        }
        html.Body().Style("margin: 100px 0 200px 0"); {

            content := html.Div().Style("max-width: 700px; margin: auto;"); {
                content.Tag("h1").Text("todo list:")

                content.Tag("button").
                    Style(buttonStyle).
                    SetAttr("hx-post", "/clicked").
                    SetAttr("hx-swap", "afterbegin transition:true").
                    SetAttr("hx-target", "#list").
                    Text("new note")

                list := content.Tag("div").Id("list"); {
                    for i, n := range slices.Backward(notes) {
                        list.Append(card(i, n))
                    }
                }
            }
        }
    }
    return html
}
type note struct {
    title, summary, body string
    date time.Time
}

func main() {
    notes := []note {
        {title: "homework", summary: "homework due on monday", body: loremParagraph},
        {title: "mom's birthday", summary: "have to buy mom a present", body: loremParagraph},
        {title: "fix bike", summary: "chain is busted", body: loremParagraph},
    }
    newNote := note{
        title: "click to edit...",
        summary: "click to edit...",
        body:  "click to edit...",
    }

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        index(notes).Render(w)
    })
    http.HandleFunc("POST /clicked", func(w http.ResponseWriter, r *http.Request) {
        n := newNote
        notes = append(notes, n)
        card(len(notes) - 1, n).Render(w)
    })
    http.HandleFunc("POST /delete/{index}", func(w http.ResponseWriter, r *http.Request) {
        fmt.Println(r.Pattern)
        index := r.PostFormValue("index")
        if index == "" {
            panic("id?")
        }
        id, err := strconv.Atoi(index)
        if err != nil {
            panic(err)
        }
        fmt.Println("id:", id)
        notes = slices.Delete(notes, id, id+1)
        w.Write([]byte(""))
    })
    fmt.Println("listening on http://localhost:2425...")
    err := http.ListenAndServe(":2425", nil)
    if err != nil {
        panic(err)
    }
}
