# Naming

- ALWAYS use full words e.g. `cellHeight` over `cellH`

# Widgets `./widgets/*.go`

Each markdown element has a corresponding vaxis widget in `./widgets/`.

Every widget has:

- A **struct** containing all information for this element. NO ast types — they are fully abstracted away at this point.
- A `new*` function that converts an ast type to the struct.
- A `Draw` function that renders the widget.
- A `Xml` function that 1:1 represents what will be rendered.

The `./widgets/` folder contains ONLY widgets, nothing else.

# Theming `./theming`

## Frontmatter `./theming/frontmatter.go`

Theming uses frontmatter and a subset of tailwind classes.

```md
---
theme:
  light:
    root: bg-white
    p: text-red-300 font-italic
    h1: font-bold
    code:
      block: bg-gray-200
      inline: bg-gray-200
      syntax:
        keyword: text-gray-500

  dark: ...
---

# Start of Slides
```

## Supported Tailwind Classes `./theming/tailwind-parser.go`

- `bg-red-200`
- `text-red-200`
- `font-bold`
- `font-italic`
- font-dim, decoration-single, decoration-solid, decoration-double, decoration-wavy, decoration-dotted, decoration-dashed

## Default Theme `./theming/default.yaml`

The default theme is minimal and adds very little styling.

When providing a partial theme without a base:

```md
---
theme:
  light:
    root: bg-white
---
```

The default theme is used as base.

If NO frontmatter theme key is present, use the `one.yaml` theme.

## Light / Dark Mode

Mode is selected via system preference. Can be changed with `t` keybinding.

## Theme Schema `./theming/schema.go`

The theme schema is validated using a JSON schema definition. Invalid themes show warnings in the footer.

## Base Theme

Base themes are loaded via the `base` keyword. Supports URLs and local paths.

```md
---
theme:
  base: https://example.com/theme.yaml
---
```

## Theme Override

Override specific keys while keeping a base theme.

```md
---
theme:
  base: ../local/theme.yaml
  light:
    root: bg-white
---
```

Merging happens at key level — when a key is set, ALL base values for that key are ignored.

# URL Loading `./fetch.go`

Themes and images can be provided by URL.

- A `resolveUrl` function resolves a URL to a temp file path.
- Resources are fetched and stored in the OS temp folder.
- Temp files are deleted on TUI close.
- ALWAYS refetch on TUI restart.
- Handle failed fetches gracefully.
