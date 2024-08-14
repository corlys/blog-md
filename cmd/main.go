package main

import (
	"bytes"
	"io"
	"os"

	"github.com/corlys/blog-md/views"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"github.com/yuin/goldmark"
)

func main() {
	server := gin.Default()
	server.GET("/blogs/:slug", func(c *gin.Context) {
		slug := c.Param("slug")
		markdownContent, err := markdownReader(slug)
		if err != nil {
			render(c, 500, views.NotFound())
			return
		}
		var buf bytes.Buffer
		if err := goldmark.Convert([]byte(markdownContent), &buf); err != nil {
			panic(err)
		}
		render(c, 200, views.Blog(templ.Raw(buf.String())))
	})
	server.Run(":3000")

}

func render(c *gin.Context, status int, template templ.Component) error {
	c.Status(status)
	return template.Render(c.Request.Context(), c.Writer)
}

func markdownReader(slug string) (string, error) {
	f, err := os.Open("blogs/" + slug + ".md")
	if err != nil {
		return "", err
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
