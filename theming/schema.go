package theming

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type ValidationWarning struct {
	Path    string
	Message string
}

func (w ValidationWarning) String() string {
	return fmt.Sprintf("%s: %s", w.Path, w.Message)
}

var validTopKeys = map[string]bool{
	"root": true, "p": true,
	"h1": true, "h2": true, "h3": true, "h4": true, "h5": true, "h6": true,
	"blockquote": true, "link": true, "link_url": true,
	"list": true, "status": true,
	"code": true, "table": true,
}

var validCodeKeys = map[string]bool{
	"block": true, "inline": true, "syntax": true,
}

var validSyntaxKeys = map[string]bool{
	"keyword": true, "string": true, "comment": true,
	"function": true, "function_method": true, "function_builtin": true, "function_name": true,
	"type": true, "number": true, "operator": true,
	"variable": true, "property": true, "constant_builtin": true, "escape": true,
}

var validTableKeys = map[string]bool{
	"header": true, "cell": true, "separator": true,
}

func Validate(frontmatterYAML []byte) []ValidationWarning {
	if len(frontmatterYAML) == 0 {
		return nil
	}
	var wrapper struct {
		Theme map[string]yaml.Node `yaml:"theme"`
	}
	if err := yaml.Unmarshal(frontmatterYAML, &wrapper); err != nil {
		return []ValidationWarning{{Path: "theme", Message: "invalid YAML: " + err.Error()}}
	}
	if wrapper.Theme == nil {
		return nil
	}

	var warnings []ValidationWarning
	for modeName, modeNode := range wrapper.Theme {
		if modeName == "base" {
			continue
		}
		if modeName != "light" && modeName != "dark" {
			warnings = append(warnings, ValidationWarning{
				Path:    "theme." + modeName,
				Message: "unknown mode (expected 'light' or 'dark')",
			})
			continue
		}
		warnings = append(warnings, validateModeNode(modeName, &modeNode)...)
	}
	return warnings
}

func validateModeNode(mode string, node *yaml.Node) []ValidationWarning {
	if node.Kind != yaml.MappingNode {
		return nil
	}
	var warnings []ValidationWarning
	for i := 0; i+1 < len(node.Content); i += 2 {
		key := node.Content[i].Value
		value := node.Content[i+1]
		path := "theme." + mode + "." + key

		if !validTopKeys[key] {
			warnings = append(warnings, ValidationWarning{Path: path, Message: "unknown theme key"})
			continue
		}

		switch key {
		case "code":
			warnings = append(warnings, validateCodeNode(path, value)...)
		case "table":
			warnings = append(warnings, validateMapNode(path, value, validTableKeys)...)
		default:
			if value.Kind == yaml.ScalarNode {
				warnings = append(warnings, validateClasses(path, value.Value)...)
			}
		}
	}
	return warnings
}

func validateCodeNode(path string, node *yaml.Node) []ValidationWarning {
	if node.Kind != yaml.MappingNode {
		return nil
	}
	var warnings []ValidationWarning
	for i := 0; i+1 < len(node.Content); i += 2 {
		key := node.Content[i].Value
		value := node.Content[i+1]
		subpath := path + "." + key

		if !validCodeKeys[key] {
			warnings = append(warnings, ValidationWarning{Path: subpath, Message: "unknown code theme key"})
			continue
		}

		if key == "syntax" {
			warnings = append(warnings, validateMapNode(subpath, value, validSyntaxKeys)...)
		} else if value.Kind == yaml.ScalarNode {
			warnings = append(warnings, validateClasses(subpath, value.Value)...)
		}
	}
	return warnings
}

func validateMapNode(path string, node *yaml.Node, validKeys map[string]bool) []ValidationWarning {
	if node.Kind != yaml.MappingNode {
		return nil
	}
	var warnings []ValidationWarning
	for i := 0; i+1 < len(node.Content); i += 2 {
		key := node.Content[i].Value
		value := node.Content[i+1]
		subpath := path + "." + key

		if !validKeys[key] {
			warnings = append(warnings, ValidationWarning{Path: subpath, Message: "unknown key"})
			continue
		}
		if value.Kind == yaml.ScalarNode {
			warnings = append(warnings, validateClasses(subpath, value.Value)...)
		}
	}
	return warnings
}

func validateClasses(path, classes string) []ValidationWarning {
	var warnings []ValidationWarning
	for _, class := range strings.Fields(classes) {
		if !ValidateClass(class) {
			warnings = append(warnings, ValidationWarning{
				Path:    path,
				Message: fmt.Sprintf("unknown class: %s", class),
			})
		}
	}
	return warnings
}
