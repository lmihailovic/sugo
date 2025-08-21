# Sugo
A minimal and dead-simple static site generator.

## Overview
Sugo is a static site generator written in Go, designed to be simple and offer
a _plug-and-play_ experience. 

It was initially imagined to be a personal tool, but later
it was decided to make it a thesis project and be made
open-source and available to the public.

## To do
- [ ] Dynamic front matter entries (ditch the current temporary struct solution)
- [ ] Specify custom templates in front matter (so content's template doesn't
necessarily depend on the type (subdir) to which the content belongs to,
e.g `blog/`, `hobbies/`, `project/`...)
- [ ] Implement `_index.md` templates which represent the index of the given subdir