package main

const BlogTitle = "A hobbyist programmer"
const BlogDescription = "A hobbyist programmer"

// func main() {
// 	initTemplates()
// 	app := &cli.Command{
// 		Name:  "SimpleSSG",
// 		Usage: "A simple static site generator",
// 		Commands: []*cli.Command{
// 			{
// 				Name:  "new",
// 				Usage: "Creates a markdown file for the new blog post",
// 				Flags: []cli.Flag{
// 					&cli.StringFlag{
// 						Name:     "title",
// 						Usage:    "title of the new blog post",
// 						Required: true,
// 					},
// 				},
// 				Action: handleNewBlogPost,
// 			},

// 			{
// 				Name:   "build",
// 				Usage:  "Builds your static site",
// 				Action: handleBuildStaticSite,
// 			},

// 			{
// 				Name:   "dev",
// 				Usage:  "Builds and starts a local server to serve the static site locally ",
// 				Action: handleBuildAndServe,
// 			},
// 		},
// 	}

// 	err := app.Run(context.Background(), os.Args)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

func main() {
	initTemplates()
	// d := time.Now().Format("Fri, 30 Jan 2026")
	// content, err := renderNewPostMD("A fresh start", d)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// err = os.WriteFile("A-Fresh-Start.md", content, 0755)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// *************************************************

	// err := recreateOutputDir()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// err = initTemplates()
	// if err != nil {
	// 	return
	// }

	// content, err := os.ReadFile("example_post.md")
	// if err != nil {
	// 	fmt.Println("could not read markdown file", err)
	// 	return
	// }

	// post, err := parseMetadataAndContent(string(content))
	// if err != nil {
	// 	fmt.Println("could not fetch metadata and content from the markdown file", err)
	// 	return
	// }

	// // render one post
	// newPost, err := renderNewPostHTML(post)
	// if err != nil {
	// 	fmt.Println("could not template.Execute new post", err)
	// 	return
	// }

	// // write that post html to file
	// err = os.WriteFile(post.Slug+".html", newPost, 0755)
	// if err != nil {
	// 	fmt.Println("could not write blog post html", err)
	// 	return
	// }

	// // render homepage
	// homePage, err := renderHomePageHTML(PostList{Posts: []Post{*post}, BlogTitle: "A hobbyst programmer."})
	// if err != nil {
	// 	fmt.Println("could not template.Execute home page", err)
	// 	return
	// }

	// err = os.WriteFile("home.html", homePage, 0755)
	// if err != nil {
	// 	fmt.Println("could not write blog post html", err)
	// 	return
	// }

	// serve()

	// ****************************************************

	// getAllPosts()

	// err := os.CopyFS(OutputDir, os.DirFS("../tanmoysrt"))
	// if err != nil {
	// 	fmt.Printf("copy assets: %w", err)
	// }

}
