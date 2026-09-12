package main

import (
	"context"
	"errors"
	"fmt"
	"os"
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
	postFileName := fmt.Sprintf("%s.md", generatedSlug)
	fmt.Printf("Creating post %s...\n", postFileName)

	err = os.WriteFile(postFileName, content, 0777)
	if err != nil {
		return err
	}
	assetsDir := fmt.Sprintf("assets/%s", generatedSlug)
	fmt.Printf("Creating assets directory %s...\n", assetsDir)
	err = os.MkdirAll(assetsDir, 0777)
	if err != nil {
		return err
	}
	fmt.Printf("Assets directory %s created successfully\n", assetsDir)

	return nil
}

func handleBuildStaticSite(ctx context.Context, c *cli.Command) error {
	// create the output folder fresh
	err := recreateOutputDir()
	if err != nil {
		return err
	}

	// build the static site

	return nil
}

func handleBuildAndServe(ctx context.Context, c *cli.Command) error {
	// create the output folder fresh
	err := recreateOutputDir()
	if err != nil {
		return err
	}

	// build the static site

	// serve the static file
	return serve()
}
