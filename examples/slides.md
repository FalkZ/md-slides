---
# You can uncomment this to create your own theme.
# theme:
#   base: ../themes/liip.yaml
#   light:
#     root: bg-white
#     h1: text-blue-600 font-bold
#   dark:
#     root: bg-slate-900
#     h1: text-sky-400 font-bold
---

# md-slides

Present from your terminal, write in markdown.

---

## How it works

Your presentation is a single markdown file. Separate slides with `---` on its own line.

```markdown
# Slide 1

Some content

---

# Slide 2

More content
```

---

## Text formatting

Markdown formatting works as you'd expect.

You can write **bold**, _italic_, and ~~strikethrough~~ text. Use `inline code` for technical terms.

---

## Lists

Unordered lists with nesting:

- Markdown features
  - Text formatting
  - Code blocks
  - Tables
- Theming via frontmatter
- Light and dark mode

Ordered lists:

1. Write your slides
2. Run `md-slides slides.md`
3. Present

---

## Code blocks

Fenced code blocks support syntax highlighting for 30+ languages.

```go
package main

import "fmt"

func main() {
	// present from your terminal
	fmt.Println("Hello, audience!")
}
```

---

## Blockquotes

> Simplicity is the ultimate sophistication.

> Blockquotes work across
>
> multiple paragraphs.

---

## Links

Inline links: [md-slides on GitHub](https://github.com/FalkZ/md-slides)

Autolinks: <https://github.com>

---

## Tables

| Feature      | Supported |
| ------------ | --------- |
| Bold, italic | yes       |
| Code blocks  | yes       |
| Tables       | yes       |
| Images       | yes       |
| Theming      | yes       |

---

## Images

Standard image:

![Make a sandwich](https://imgs.xkcd.com/comics/sandwich.png)

_Image from: <https://xkcd.com/149/>_

---

## Images with custom height

Control the height in terminal rows with `| N` in the alt text:

```markdown
![Alt text | 10](./image.png)
```

![Make a sandwich | 8](./xkcd-sandwich.png)

_Image from: <https://xkcd.com/149/>_

If an image can't load, the alt text is shown instead:

![This alt text is shown as fallback](./nonexistent.png)

---

## Theming

Style your slides with a YAML frontmatter block and Tailwind classes:

```yaml
---
theme:
  light:
    h1: text-blue-600 font-bold
    code:
      syntax:
        keyword: text-purple-600
  dark:
    h1: text-sky-400 font-bold
---
```

Press `t` to toggle between light and dark mode.

---

## Keyboard controls

| Key                   | Action         |
| --------------------- | -------------- |
| Right / l / n / Space | Next slide     |
| Left / h / p          | Previous slide |
| g / Home              | First slide    |
| G / End               | Last slide     |
| t                     | Toggle theme   |
| q / Ctrl+c            | Quit           |

---

## Theme inheritance

Use `base` to inherit from a base theme (URL or local path):

```yaml
---
theme:
  base: https://example.com/theme.yaml
  light:
    h1: text-red-600
---
```

### Merge Order

Merge order: default / base → local overrides.

---

# Go present something.

---

# Reference

## Markdown

- **Bold**: `**bold**`
- _Italic_: `_italic_`
- ~~Strikethrough~~: `~~strikethrough~~`
- `Inline code`: `` `code` ``
- [Links](https://example.com): `[text](url)`
- Autolinks: `<https://example.com>`
- Images: `![alt](path)`
- Images with height: `![alt | rows](path)`

## Code Blocks

```yaml
---
theme:
  light:
    h1: text-blue-600 font-bold
    code:
      syntax:
        keyword: text-purple-600
  dark:
    h1: text-sky-400 font-bold
---
```

## Lists

- Unordered item
  - Nested item

1. Ordered item
2. Another item

## Blockquotes

> Quote text

## Tables

| Header | Header |
| ------ | ------ |
| Cell   | Cell   |

## Images

![Make a sandwich | 8](./xkcd-sandwich.png)

_Image from: <https://xkcd.com/149/>_

If an image can't load, the alt text is shown instead:

![This alt text is shown as fallback](./nonexistent.png)
