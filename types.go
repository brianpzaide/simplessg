package main

import (
	"html/template"
	"time"
)

type Post struct {
	BlogTitle     string        `yaml:"-"`
	Title         string        `yaml:"title"`
	Description   string        `yaml:"description"`
	Date          string        `yaml:"date"`
	DateObj       *time.Time    `yaml:"-"`
	FormattedDate string        `yaml:"-"`
	Content       template.HTML `yaml:"-"`
	Slug          string        `yaml:"-"`
	FilePath      string        `yaml:"-"`
	Route         string        `yaml:"-"`
}

type PostList struct {
	BlogTitle string
	Posts     []Post
}

type RSSChannel struct {
	Title         string    `xml:"title"`
	Link          string    `xml:"link"`
	Description   string    `xml:"description"`
	LastBuildDate string    `xml:"lastBuildDate"`
	AtomLink      AtomLink  `xml:"atom:link"`
	Items         []RSSItem `xml:"item"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
	GUID        string `xml:"guid"`
}

type AtomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}

type envelope map[string]interface{}
