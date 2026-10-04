## Simplessg

simplessg is a simple static site generator

#### Usage
simplessg is a CLI with three commands:

1) `new --title <title>` — Creates a new blog post. This command creates `<title_slug>.md` file and a corresponding directory named `<title_slug>` under assets/ folder, for assets used by this new blog post.
2) `build`— Builds the static site in the `dist/` directory.
3) `dev`— Builds the static site and starts a development server listening on port 4000.

#### Directory structure
simplessg expects the following directory structure:

```bash
.
├── assets
│   ├── hello-world
│   ├── post-1
│   ├── post-2
│   │   └── img1.png
│   ├── post-3
│   │   ├── img1.avif
│   │   ├── img2.avif
│   │   └── img3.avif
│   └── post-4
│       ├── img1.png
│       └── img2.png
├── hello-world.md
├── post-1.md
├── post-2.md
├── post-3.md
└── post-4.md
```

#### Using with Github Actions
If your blog is maintained as a GitHub repository, you can use the following GitHub Actions workflow to automatically build and deploy the static site to GitHub Pages whenever changes are pushed to the main branch.

```yaml
name: Build and Deploy

on:
  push:
    branches:
      - main

permissions:
  contents: read
  pages: write
  id-token: write

jobs:
  build:
    runs-on: ubuntu-latest

    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Build site
        run: |
          docker run --rm -v "$PWD:/simplessg/posts" brianpzaide/simplessg

      - name: Upload Pages artifact
        uses: actions/upload-pages-artifact@v4
        with:
          path: ./dist

  deploy:
    runs-on: ubuntu-latest
    needs: build

    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}

    steps:
      - name: Deploy to GitHub Pages
        id: deployment
        uses: actions/deploy-pages@v4
```

#### Docker

`simplessg` is also available as a Docker image

```bash
docker run --rm -v "$PWD:/simplessg/posts" brianpzaide/simplessg
```
You can also use the provided Dockerfile to build your own image and host it on a container registry.


### Credit
Ideas are taken from [tanmoysrt's static site generator](https://github.com/tanmoysrt/tanmoysrt.dev)

