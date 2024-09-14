# gtag
A simple [yattag](https://github.com/leforestier/yattag) and [gomponents](https://github.com/maragudk/gomponents) inspired html templating library for Go.
```go
html := gtag.Doc(); {
    head := html.Head(); {
        head.Tag("title").Text("hello world")
        head.Tag("style").Asis("body { background-color: black; color: white}")
    }
    body := html.Body(); {
        body.Tag("marquee").Text("hello world!")
    }
}
fmt.Println(html.String())
```
The above example features braces being used to structure the relationship between parent and child. you're not obliged to do that. I aimed to design something similar to yattag's use of python's `with` statement and found that this combination of declaring a handler for the tag + brances gets the closeset to the real thing. there are some pitfalls to this method. be cautious. declaring the handles inside their scope also does the job but you do you.
```go
{
    html := gtag.Doc(); 
    {
        head := html.Head(); 
        head.Tag("title").Text("hello world")
        head.Tag("style").Asis("body { background-color: black; color: white}")
    }
    {
        body := html.Body(); 
        body.Tag("marquee").Text("hello world!")
    }
}
fmt.Println(html.String())
```

the main idea was to have a framework that could levarage Go's static typing while also being composable, statement-frieldly and didn't require additional build steps (gomponents and [templ](https://github.com/a-h/templ) fail at two of those things). I believe gtag can excel at those things at the cost of being incompatible with HTML's standard syntax which at this day and age means AI tools can't help you. 

There might be other similar libraries out there working with simiar concepts. I haven't done my research properly.

## Documentation
Read the source code for a documentation. The entire thing is 120 lines only. keep an eye out for what references are being returned from the functions. use empty string to declare boolean attributes.

## Basic Examples

statement-friendliness:
```go
html := gtag.Doc(); {
    html.Head().Tag("title").Text("fizzbuz")
    body := html.Body(); {
        div := body.Div(); {
            for i := range 31 {
                i3, i5 := i % 3 == 0, i % 5 == 0
                if i3 && i5 {
                    div.P().Text("fizzbuzz")
                } else if i5 {
                    div.P().Text("buzz")
                } else if i3 {
                    div.P().Text("fizz")
                } else {
                    div.P().Text(strconv.Itoa(i))
                }
            }
        }
    }
}
html.Render(os.Stdout)
```

composition:
```go
func card(name string, age int) *gtag.Tag {
    div := gtag.New("div"); {
        div.P().Text("Name: ").Tag("span").Text(name)
        div.P().Text("Age: ").Tag("span").Text(strconv.Itoa(age))
    }
    return div
} 

func main() {
    html := gtag.Doc(); {
        head := html.Head().Tag("title").Text("hello world")

        body := html.Body(); {
            body.Append(card("John", 34))
            body.Append(card("Jane", 25))
            body.Append(card("Daniel", 22))
        }
    }
    fmt.Println(html.Render())
}

```

wip...
