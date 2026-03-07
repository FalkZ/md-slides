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

- Headings (`# h1` through `###### h6`)
- Paragraphs
- Unordered and ordered lists
- **Bold**, _italic_, ~~strikethrough~~
- Fenced code blocks with syntax highlighting
- `Inline code`
- Blockquotes
- Tables (GFM-style)
- [Links](https://example.com) and auto-links
- Images

## Images

> For best performance use a terminal that supports Kitty. For example Ghostty.

Standard markdown image syntax:

```markdown
![Alt text](./image.png)
```

Optionally specify height in terminal rows using a pipe delimiter in the alt text:

```markdown
![Alt text | 12](./image.png)
```

If omitted, images default to 20 rows. Image paths are resolved relative to the markdown file's directory. Supports Kitty graphics, Sixel, and standard terminal graphics protocols.

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
| `hr`              | Horizontal rules    |
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

**Background:** `bg-{color}-{shade}`, `bg-[#RRGGBB]`

**Underline:** `underline`, `decoration-double`, `decoration-wavy`, `decoration-dotted`, `decoration-dashed`, `no-underline`

Colors: slate, gray, zinc, neutral, stone, red, orange, amber, yellow, lime, green, emerald, teal, cyan, sky, blue, indigo, violet, purple, fuchsia, pink, rose (shades: 50-950)

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

| Flag                | Description                           |
| ------------------- | ------------------------------------- |
| `--debug`           | Write debug XML output to disk        |
| `--schema`          | Print theme JSON schema and exit      |
| `--dump <page>`     | Dump rendered page as XML (1-indexed) |
| `--dump-raw <page>` | Dump raw terminal output of a page    |
