# sugo

A minimal and dead-simple static site generator.

## Overview

Sugo is a static site generator written in Go, designed to be simple and offer
a _plug-and-play_ experience. 

It was initially imagined as a personal tool, but later it was decided to
make it a thesis project and be made open-source and available to the public.

## Installation

Download the binary from the
[Releases](https://github.com/lmihailovic/sugo/releases) page, if there is one.  
If there is no binary, you must compile from source.

### Compiling from source

```
git clone https://github.com/lmihailovic/sugo
cd sugo
go build
```

## Arguments

```
-d      run dev server
-i string
        path to website directory (default ".")
-o string
        path for generated static web files (default "website")
```

## Getting started

Create a root directory where you wish to contain the resources for your website.  
This directory further needs `content/` and `templates/` directories, with
`static/` as an optional one, should you need it. See the
[section on static files](#static-files) for its use cases.

The root directory can be explicitly set using the `-i` flag, see
[the arguments section](#arguments) for other flags.

After running sugo, a new directory, `website/`, is created by default inside
of the root directory you created yourself. 

### Writing content files

The content files use markdown as their file format.

Front matter is written at the very top of the file, and is delimited by "+++".

Front matter format is **JSON**. Sugo will refuse to build the website if any
page has badly formatted JSON.

Example content file:
```
+++
"Title": "Content file formatting",
"Description": "Showcases how what can be done in a content file",
"Date": "14-9-2025"
+++
This here text represents the content itself. Everything above this is
front matter.

# What can I use here?

Here I can use regular *markdown*, and **all** of it's standard features:
- bullet lists
- emphasising
- defining titles
- hyperlinks

I was extra nice, so I even allowed for tables :)
Big thanks to the [goldmark](https://github.com/yuin/goldmark) library for
markdown parsing!
```

### Templating of content files

Subdirectories inside of `content/` are treated like different sections. 
Each section requires its own template, found in the `templates/` directory. 
Templates can be either `section.gohtml`, or `single.gohtml`. 

Index files for sections use the `section.gohtml` template, while all other
files inside of a section use the `single.gohtml` template.

**Content and template files must have the same relative path in order to work.**

This means that the file `content/blog/lorem.md` will need a
`templates/blog/single.gohtml` template to render properly, and a
`content/hobbies/films/index.md` file would need
`templates/hobbies/films/section.gohtml` as the template.

Here is an example of a typical `content/` directory, and the accompanying
`templates/` directory:
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

```
templates
├── blog
│   ├── section.gohtml
│   └── single.gohtml
└── hobbies
    ├── section.gohtml
    ├── single.gohtml
    └── films
        ├── section.gohtml
        └── single.gohtml 
```

### Static files

Some files just need to be copied over, without changing
(images, styles, fonts...). For that need, just place them in the `static/`
directory, and they will be copied over to the generated website, nothing
changed.

## Quick tips

### Navigation links menu

To create a navigational list of links for a specific page, you might use
the following code
```
{{ range $link, $data := GetChildPages "blog" false }}
    <a href="{{ $link }}"> {{ $data.Title }} </a>
{{ end }}
```
This generates relative links for each `.html` file inside of `blog/`.

## To do

- [ ] Command to generate an example website
- [ ] Make a live server functionality for live previews. Possibly replace the existing `-d` functionality with this one.
- [ ] Get a chef gopher as a mascot
