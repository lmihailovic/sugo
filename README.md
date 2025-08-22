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
│   └── lorem.md
└── hobbies
    └── films
        └── index.md
```

Given the example above, path `example.com/blog` will use the `section.html`
template, and `example.com/blog/lorem.html` will use `single.html`.

## To do
- [x] Dynamic front matter entries (ditch the current temporary struct solution)
- [ ] Specify custom templates in front matter (so content's template doesn't
necessarily depend on the type (subdir) to which the content belongs to,
e.g `blog/`, `hobbies/`, `project/`...)
- [x] Implement `index.md` templates which represent the index of the given subdir
- [ ] Implement special `home` templates for the website's main index page
- [ ] Implement nested templates (for the document head, header, footer...)