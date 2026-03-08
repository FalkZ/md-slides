# md-slides

A lightweight terminal presentation tool. Written in Go as a simpler alternative to [slides](https://github.com/maaslalani/slides) with fewer dependencies, only essential features, and image support.

## Installation

```
go install github.com/FalkZ/md-slides@latest
```

## Usage

```
md-slides presentation.md
```

## Slide Separation

Slides are separated by a thematic break (`---`):

```markdown
# First Slide

Some content

---

# Second Slide

More content
```

## Markdown Support

- Headings
- Paragraphs
- Lists (unordered and ordered)
- **Bold**, _italic_, ~~strikethrough~~
- Code blocks (with syntax highlighting)
- Inline code
- Blockquotes
- Tables (GFM-style)
- [Links](https://example.com) and auto-links
- Images

## Images

Standard markdown image syntax:

```markdown
![Alt text](./image.png)
```

### Image size by rows

Optionally specify image height in rows using a pipe delimiter in the alt text:

```markdown
![Alt text | 12](./image.png)
```

If omitted, images default to 20 rows.

### Image support

For best results, use a terminal that supports the Kitty graphics protocol, such as Ghostty or Kitty. Sixel and other standard terminal graphics protocols are also supported.

Image paths are resolved relative to the markdown file's directory.

## Syntax Highlighting

Code blocks use Tree-sitter based syntax highlighting. Specify the language after the opening fence:

````markdown
```go
fmt.Println("Hello, world!")
```
````

Supports all languages available in [gotreesitter](https://github.com/odvcencio/gotreesitter) including Go, Rust, Python, JavaScript, TypeScript, C, C++, Java, Ruby, and many more.

## Theming

Themes are defined in YAML frontmatter using Tailwind-style utility classes. Separate styles for light and dark mode:

```yaml
---
theme:
  light:
    root: bg-[#fafafa]
    p: text-zinc-700
    h1: text-blue-500 font-bold decoration-double
    code:
      block: text-zinc-400
      inline: text-red-500
      syntax:
        keyword: text-purple-600 font-bold
        string: text-green-600
  dark:
    root: bg-[#282c34]
    p: text-zinc-400
    h1: text-sky-400 font-bold decoration-double
---
```

### Themeable Keys

| Key               | Description         |
| ----------------- | ------------------- |
| `root`            | Root background     |
| `p`               | Paragraph text      |
| `h1` - `h6`       | Headings            |
| `blockquote`      | Blockquotes         |
| `link`            | Links               |
| `img`             | Image fallback text |
| `list`            | Lists               |
| `status`          | Status bar          |
| `code.block`      | Code blocks         |
| `code.inline`     | Inline code         |
| `code.syntax.*`   | Syntax highlighting |
| `table.header`    | Table headers       |
| `table.cell`      | Table cells         |
| `table.separator` | Table separators    |

### Syntax Highlighting Keys

`keyword`, `string`, `comment`, `function`, `function_method`, `function_builtin`, `function_name`, `type`, `number`, `operator`, `variable`, `property`, `constant_builtin`, `escape`

### Supported Classes

**Font:** `font-bold`, `font-italic`, `font-dim`

**Text color:** `text-{color}-{shade}`, `text-white`, `text-black`, `text-[#RRGGBB]`

**Background:** `bg-{color}-{shade}`, `bg-white`, `bg-black`, `bg-[#RRGGBB]`

**Underline:** `decoration-single`, `decoration-solid`, `decoration-double`, `decoration-wavy`, `decoration-dotted`, `decoration-dashed`

Colors: slate, gray, zinc, neutral, stone, red, orange, amber, yellow, lime, green, emerald, teal, cyan, sky, blue, indigo, violet, purple, fuchsia, pink, rose (shades: 50-950)

### Theme Extension

Use the `extends` property to base your theme on an external theme file (local path or URL). Only the keys you specify are overridden; everything else is inherited from the extended theme:

```yaml
---
theme:
  extends: https://example.com/my-base-theme.yaml
  light:
    h1: text-pink-500 font-bold
  dark:
    h1: text-pink-400 font-bold
---
```

The merge order is: **default theme → extended theme → frontmatter overrides**. This means you can create reusable theme files and customize only what you need per presentation.

Partial themes work at the property level — setting `h1: text-pink-500` only changes the heading color while preserving any background, underline, or font style from the base theme.

## Default Theme

Based on [Zed](https://github.com/zed-industries/zed) One Dark / One Light. Automatically switches between light and dark mode based on your terminal's theme.

## Controls

| Key                           | Action                  |
| ----------------------------- | ----------------------- |
| `Right` / `l` / `n` / `Space` | Next slide              |
| `Left` / `h` / `p`            | Previous slide          |
| `g` / `Home`                  | First slide             |
| `G` / `End`                   | Last slide              |
| `t`                           | Toggle light/dark theme |
| `q` / `Ctrl+c`                | Quit                    |

## CLI Flags

| Flag                 | Description                           |
| -------------------- | ------------------------------------- |
| `--debug`            | Write debug XML output to disk        |
| `--schema`           | Print theme JSON schema and exit      |
| `--dump-xml <page>`  | Dump rendered page as XML (1-indexed) |
| `--dump-text <page>` | Dump raw terminal output of a page    |
