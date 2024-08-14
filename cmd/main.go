package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"strings"

	"github.com/corlys/blog-md/views"

	"github.com/a-h/templ"
	"github.com/adrg/frontmatter"
	"github.com/gin-gonic/gin"
	"github.com/yuin/goldmark"
)

type FrontMatter struct {
	Title string
}

func main() {
	server := gin.Default()
	server.Static("/dist", "./dist")
	server.GET("/", func(c *gin.Context) {
		var fileNames []string
		files, err := os.ReadDir("blogs/")
		for _, item := range files {
			fileNames = append(fileNames, item.Name())
		}
		if err != nil {
			log.Println(err)
			render(c, 500, views.ErrorMessage("Internal Server Error"))
			return
		}
		render(c, 200, views.Home(fileNames))
	})
	server.GET("/blogs/:slug", func(c *gin.Context) {
		slug := c.Param("slug")
		matter, markdownContent, err := markdownReader(slug)
		if err != nil {
			render(c, 500, views.ErrorMessage("Item Not Found"))
			return
		}
		var buf bytes.Buffer
		if err := goldmark.Convert([]byte(markdownContent), &buf); err != nil {
			panic(err)
		}
		render(c, 200, views.Blog(templ.Raw(buf.String()), matter.Title))
	})
	server.Run(":3000")
}

func render(c *gin.Context, status int, template templ.Component) error {
	c.Status(status)
	return template.Render(c.Request.Context(), c.Writer)
}

func markdownReader(slug string) (FrontMatter, string, error) {
	var matter FrontMatter
	f, err := os.Open("blogs/" + slug + ".md")
	if err != nil {
		return FrontMatter{}, "", err
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		return FrontMatter{}, "", err
	}
	rest, err := frontmatter.Parse(strings.NewReader(string(b)), &matter)
	if err != nil {
		return FrontMatter{}, "", err
	}
	return matter, string(rest), nil
}
