package main

import (
	"bytes"
	"encoding/xml"
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

	// copying the assets folder into the output dir
	err := os.CopyFS(OutputDir, os.DirFS("assets"))
	if err != nil {
		return fmt.Errorf("copy assets: %w", err)
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

func buildStaticSite() error {
	err := recreateOutputDir()
	if err != nil {
		fmt.Println("error in function recreateOutputDir", err)
		return err
	}

	/* walk the folder fetching all the '.md' files. For each of these files:
	-> run 'parseMetadataAndContent'
	-> run 'renderNewPostHTML'
	-> write the html file for this post content
	*/
	posts, err := getAllPosts()
	if err != nil {
		return err
	}
	for _, p := range posts {
		err = genHTMLFileForPost(p)
		if err != nil {
			return err
		}
	}

	// generating HTML home page
	err = genHTMLFileForHome(posts)
	if err != nil {
		return err
	}

	// generating RSS feed

	return nil
}

func getAllPosts() ([]*Post, error) {
	fis, err := os.ReadDir(".")
	if err != nil {
		fmt.Println("could not the post directory", err)
		return nil, err
	}
	posts := make([]*Post, 0)
	for _, fi := range fis {
		if !fi.IsDir() && strings.HasSuffix(fi.Name(), ".md") {
			mdContent, err := os.ReadFile(fi.Name())
			if err != nil {
				fmt.Printf("could not the markdown file '%s'\n", fi.Name())
				fmt.Println(err)
				return nil, err
			}
			p, err := parseMetadataAndContent(string(mdContent))
			if err != nil {
				fmt.Printf("could not parse metada and content for '%s'\n", fi.Name())
				fmt.Println(err)
				return nil, err
			}
			posts = append(posts, p)
		}
	}
	return posts, nil
}

func genHTMLFileForPost(post *Post) error {
	postHTMLContent, err := renderNewPostHTML(post)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s/%s.html", OutputDir, post.Slug)

	err = os.WriteFile(fileName, postHTMLContent, 0755)
	if err != nil {
		fmt.Printf("could not write blog post html file '%s'\n", fileName)
		return err
	}

	return nil
}

func genHTMLFileForHome(posts []*Post) error {
	homePage, err := renderHomePageHTML(PostList{Posts: posts, BlogTitle: BlogTitle})
	if err != nil {
		fmt.Println("could not template.Execute home page", err)
		return err
	}

	err = os.WriteFile("index.html", homePage, 0755)
	if err != nil {
		fmt.Println("could not write blog post html", err)
		return err
	}

	return nil
}

func genXMLFileForRSS() error {

	rssStruct := RSS{
		XMLName: xml.Name{Space: "", Local: "rss"},
		XMLNS:   "http://www.w3.org/2005/Atom",
		Version: "2.0",
		Channel: RSSChannel{
			Title:         BlogTitle,
			Description:   BlogDescription,
			LastBuildDate: time.Now().Format("Fri, 30 Jan 2026"),
			AtomLink: AtomLink{
				Href: "/rss.xml",
				Rel:  "self",
				Type: "application/rss+xml",
			},
			Items: []RSSItem{},
		},
	}

	rssFeed, err := xml.Marshal(rssStruct)
	if err != nil {
		fmt.Println("could not marshal to xml encoding", err)
		w.WriteHeader(500)
		return
	}

	return nil
}
