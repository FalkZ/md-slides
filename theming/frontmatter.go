package theming

import "bytes"

func ExtractFrontmatter(raw []byte) (frontmatter []byte, body []byte) {
	trimmed := bytes.TrimLeft(raw, " \t\r\n")
	if !bytes.HasPrefix(trimmed, []byte("---")) {
		return nil, raw
	}
	start := bytes.Index(raw, []byte("---"))
	afterFirst := start + 3
	for afterFirst < len(raw) && raw[afterFirst] != '\n' {
		afterFirst++
	}
	if afterFirst < len(raw) {
		afterFirst++
	}
	rest := raw[afterFirst:]
	end := bytes.Index(rest, []byte("\n---"))
	if end == -1 {
		return nil, raw
	}
	frontmatter = rest[:end]
	bodyStart := afterFirst + end + 4
	for bodyStart < len(raw) && raw[bodyStart] != '\n' {
		bodyStart++
	}
	if bodyStart < len(raw) {
		bodyStart++
	}
	return frontmatter, raw[bodyStart:]
}
