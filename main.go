package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
	"unicode"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

func main() {

	source := []byte(`
# hello world
	
This is an example
	
- list
- of
- items

` + "```mermaid" + `
some code
` + "```" + `

<details class="mermaid">
    <summary><img src="https://mermaid.ink/img/eyJjb2RlIjoiZ3JhcGggVERcbkFbQ2hyaXN0bWFzXSAtLT58R2V0IG1vbmV5fCBCKEdvIHNob3BwaW5nKVxuQiAtLT4gQ3tMZXQgbWUgdGhpbmt9XG5DIC0tPnxPbmV8IERbTGFwdG9wXVxuQyAtLT58VHdvfCBFW2lQaG9uZV1cbkMgLS0-fFRocmVlfCBGW2ZhOmZhLWNhciBDYXJdXG4iLCJtZXJtYWlkIjp7InRoZW1lIjoiZGVmYXVsdCJ9fQ" /></summary>
	<p>
` + "```mermaid" + `
some code
` + "```" + `
	</p>
</details>

> Good bye
`)
	reader := text.NewReader(source)

	// gm := goldmark.New()
	// root := gm.Parser().Parse(reader)
	root := goldmark.DefaultParser().Parse(reader)

	rendered := []*MermaidRendered{}
	fences := []*ast.FencedCodeBlock{}

	ast.Walk(root, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		switch tnode := n.(type) {
		case *ast.HTMLBlock:
			if !entering {
				return ast.WalkContinue, nil
			}
			lastLine := tnode.Lines().Len() - 1
			first := tnode.Lines().At(0).Start
			last := tnode.Lines().At(lastLine).Stop

			block := parseHTML(source[first:last])
			rendered = append(rendered, &MermaidRendered{
				Segment: text.NewSegment(first, last),
				Block:   block,
			})

		case *ast.FencedCodeBlock:
			if !entering {
				return ast.WalkContinue, nil
			}
			if string(tnode.Language(source)) == "mermaid" {
				fences = append(fences, tnode)
			}
		default:
			// fmt.Printf("%T : %s : %t\n", n, string(n.Text(source)), n.IsRaw())
		}
		return ast.WalkContinue, nil
	})

	fmt.Println("Found rendered mermaid blocks")
	for _, r := range rendered {
		fmt.Println(string(r.Marshal()))
	}

	fmt.Println("")
	fmt.Println("Found fenced mermaid blocks")
	for _, r := range fences {
		fmt.Print(r.Lines())
		lastLine := r.Lines().Len()
		first := r.Lines().At(0).Start
		last := r.Lines().At(lastLine - 1).Stop
		fmt.Println(string(source[first:last]))
	}
	// gm.Renderer().Render(os.Stdout, source, root)
}

