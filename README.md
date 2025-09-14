# Sugo
A minimal and dead-simple static site generator.

## Overview
Sugo is a static site generator written in Go, designed to be simple and offer
a _plug-and-play_ experience. 

It was initially imagined to be a personal tool, but later
it was decided to make it a thesis project and be made
open-source and available to the public.

## Modus operandi
Each subdir inside of `content/` is treated like a different content type.

Index pages for content types utilise the `section.html` template, while
each content file of a certain type uses the `single.html` template.

Here is an example of a typical `content/` directory:
```
content
├── blog
│   ├── darkmode-difficulties.md
│   ├── lorem.md
│   └── index.md
└── hobbies
    ├── index.md
    └── films
        └── index.md
```

Given the example above, path `example.com/blog` will use the `section.html`
template, and `example.com/blog/lorem.html` will use `single.html`.

## Arguments
```
-d      run dev server
-i string
        path to website directory (default ".")
-o string
        path for generated static web files (default "website")
```

## Getting started
To create a navigational list of links for a specific page, you might use
the following code
```
{{ range $link, $data := GetChildPages "blog" false }}
    <a href="{{ $link }}"> {{ $data.Title }} </a>
{{ end }}
```
This generates relative links for each `.html` file inside of `blog/`.

## To do
- [x] Dynamic front matter entries (ditch the current temporary struct solution)
- [x] Specify custom templates in front matter (so content's template doesn't
necessarily depend on the type (subdir) to which the content belongs to,
e.g `blog/`, `hobbies/`, `project/`...)
- [x] Implement `index.md` templates which represent the index of the given subdir
- [x] Implement special `home` templates for the website's main index page -
implemented by placing an `index.md` in `content/` root and a `section.html` in
`template/` root
- [x] Allow for specification of custom content front matter delimiters
- [x] Add `static/` dir functionality for css, js and image files.
- [x] Implement nested templates (for the document head, header, footer...)
- [x] Add `.Path` property for pages to allow for nav element creations -
realised via function
- [x] Add ability to loop over pages in content sections inside of templates - 
realised via function
- [x] Function to loop pages over just one level of depth for a section
- [x] Pack all of page front matter in the value part of map in `GetChildPages`
- [ ] Command to generate an example website
- [ ] Get a chef gopher as a mascot