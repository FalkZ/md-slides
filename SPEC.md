# Naming

- Always use full words e.g. cellHeight over cellH

# Widgets `./widgets/*.go`

- Every markdown element has a corresponding vaxis widget in `./widgets`
- Every widget has
  - a struct
    - containing all the information of this object
    - but no ast types (at this point they are fully abstracted away)
  - an `new*` function that converts an ast type to a struct
  - a Draw function that renders the widget as
  - a Xml function that 1:1 represents what will be rendered as a widget
  - the `./widgets/` folder only contains widgets nothing else

# Theming `./theming`

## Frontmatter `./theming/parse-frontmatter.go`

Theming is done via frontmatter and a subset of tailwind.

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

## Supported tailwind classes `./theming/tailwind-parser.go`

- bg-red-200
- text-red-200
- font-bold
- font-italic

## Default theme `./theming/default-theme.yaml`

There is a default theme that mimics the zed one dark & light themes closely.
Values that are not provided in the frontmatter fall back to the default.

## Light / dark mode

Selection of the mode is done via system preference. and can be changed with a keybinding.

## Theme schema `./theming/schema.go`

The whole theme schema is validated and gives warnings in the footer when theme is invalid.
A JSON schema definition is used to validate.
