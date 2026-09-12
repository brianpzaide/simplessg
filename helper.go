package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"os"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/gosimple/slug"
	"gopkg.in/yaml.v3"
)

// the following are used for creating a new blog post

const NewPostMD = `
---
title: "{{ .Title }}"
description: "Write your description here"
date: "{{ .Date}}"
---

Write your post here
`

var HomepageHTML, NewPostHTML string
var postMDTemplate, homeHTMLTemplate, postHTMLTemplate *template.Template

func initTemplates() error {
	homec, err := os.ReadFile("templates/home.html")
	if err != nil {
		fmt.Println("error in reading homepage template", err)
		return err
	}
	HomepageHTML = string(homec)
	homeHTMLTemplate = template.Must(template.New("homehtml").Parse(HomepageHTML))

	postc, err := os.ReadFile("templates/post.html")
	if err != nil {
		fmt.Println("error in reading post template", err)
		return err
	}
	NewPostHTML = string(postc)
	postHTMLTemplate = template.Must(template.New("posthtml").Parse(NewPostHTML))

	postMDTemplate = template.Must(template.New("postmd").Parse(NewPostMD))

	return nil
}

func renderNewPostMD(title, date string) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := postMDTemplate.Execute(buf, struct {
		Title string
		Date  string
	}{
		Title: title,
		Date:  date,
	})

	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// the following are used for building a static site
func recreateOutputDir() error {
	if err := os.RemoveAll(OutputDir); err != nil {
		return fmt.Errorf("remove %q: %w", OutputDir, err)
	}

	if err := os.MkdirAll(OutputDir, 0755); err != nil {
		return fmt.Errorf("create %q: %w", OutputDir, err)
	}

	return nil
}

func parseMetadataAndContent(content string) (*Post, error) {
	content = strings.TrimSpace(content)
	frontMatterNotFound := errors.New("no front matter found")

	if !strings.HasPrefix(content, "---") {
		return nil, frontMatterNotFound
	}
	content = content[3:]

	metadata, mdContent, found := strings.Cut(content, "---")
	if !found {
		return nil, frontMatterNotFound
	}
	post := Post{}
	err := yaml.Unmarshal([]byte(metadata), &post)
	if err != nil {
		return nil, errors.New("error parsing YAML")
	}
	// format date
	parsedDate, err := time.Parse("2006-01-02", post.Date)
	if err != nil {
		return nil, fmt.Errorf("error parsing date: %v", err)
	}
	post.DateObj = &parsedDate
	post.FormattedDate = parsedDate.Format("Jan 2, 2006")
	post.Slug = slug.Make(post.Title)
	post.Route = fmt.Sprintf("%s.html", post.Slug)

	// convert content to HTML
	post.Content = template.HTML(markdownToHTML(strings.TrimSpace(mdContent)))

	return &post, nil
}

func markdownToHTML(md string) string {
	// parse markdown
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(md))

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	htmlContent := string(markdown.Render(doc, renderer))
	return htmlContent
}

func renderNewPostHTML(post *Post) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := postHTMLTemplate.Execute(buf, post)

	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func renderHomePageHTML(posts PostList) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := homeHTMLTemplate.Execute(buf, posts)

	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
