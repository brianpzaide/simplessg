package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gosimple/slug"
	"github.com/urfave/cli/v3"
)

func handleNewBlogPost(ctx context.Context, c *cli.Command) error {
	title := c.String("title")
	if title == "" {
		return errors.New("Post title cannot be empty")
	}
	date := time.Now().Format("2006-01-02")
	content, err := renderNewPostMD(title, date)
	if err != nil {
		return err
	}

	generatedSlug := slug.Make(title)
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("could not get the current working directory in cmd: handleNewBlogPost")
	}
	postFileName := filepath.Join(wd, "posts", fmt.Sprintf("%s.md", generatedSlug))
	fmt.Printf("Creating post %s...\n", postFileName)

	err = os.WriteFile(postFileName, content, 0777)
	if err != nil {
		return err
	}
	assetsDir := filepath.Join(wd, "posts", "assets", generatedSlug)
	fmt.Printf("Creating assets directory %s...\n", assetsDir)
	err = os.MkdirAll(assetsDir, 0777)
	if err != nil {
		return err
	}
	fmt.Printf("Assets directory %s created successfully\n", assetsDir)

	return nil
}

func handleBuildStaticSite(ctx context.Context, c *cli.Command) error {

	return buildStaticSite()
}

func handleBuildAndServe(ctx context.Context, c *cli.Command) error {

	err := buildStaticSite()
	if err != nil {
		return err
	}

	// serve the static file
	return serve()
}
