//go:build tools

package main

import (
	"sort"
)

// swag v2 wraps every request body schema in a redundant
// oneOf: [{type: object}, {$ref: ...}]; clients need the bare $ref.
func fixRequestBody(node map[string]any) {
	variants, ok := node["oneOf"].([]any)
	if !ok || len(variants) != 2 {
		return
	}

	first, ok := variants[0].(map[string]any)
	if !ok || len(first) != 1 || first["type"] != "object" {
		return
	}

	second, ok := variants[1].(map[string]any)
	if !ok {
		return
	}

	ref, ok := second["$ref"].(string)
	if !ok {
		return
	}

	delete(node, "oneOf")
	node["$ref"] = ref
}

// "@Security []" is emitted as [{"": [""]}] instead of an empty list.
func fixEmptySecurity(node map[string]any) {
	security, ok := node["security"].([]any)
	if !ok {
		return
	}

	cleaned := make([]any, 0, len(security))

	for _, entry := range security {
		requirement, ok := entry.(map[string]any)
		if !ok {
			cleaned = append(cleaned, entry)

			continue
		}

		if _, blank := requirement[""]; blank && len(requirement) == 1 {
			cleaned = append(cleaned, map[string]any{})

			continue
		}

		cleaned = append(cleaned, requirement)
	}

	node["security"] = cleaned
}

// swag emits "type: file" under x-www-form-urlencoded for formData file
// params; the real handler reads a multipart field named photo.
func fixMultipart(node map[string]any) {
	content, ok := node["content"].(map[string]any)
	if !ok {
		return
	}

	if _, multipart := content["multipart/form-data"]; !multipart {
		return
	}

	urlencoded, ok := content["application/x-www-form-urlencoded"].(map[string]any)
	if !ok {
		return
	}

	schema, ok := urlencoded["schema"].(map[string]any)
	if !ok || schema["type"] != "file" {
		return
	}

	field, ok := schema["title"].(string)
	if !ok || field == "" {
		field = "file"
	}

	delete(content, "application/x-www-form-urlencoded")

	content["multipart/form-data"] = map[string]any{
		"schema": map[string]any{
			"type":     "object",
			"required": []any{field},
			"properties": map[string]any{
				field: map[string]any{
					"type":        "string",
					"format":      "binary",
					"description": "Файл изображения — jpeg, png или webp.",
				},
			},
		},
	}
}

func walk(node any, visit func(map[string]any)) {
	switch typed := node.(type) {
	case map[string]any:
		visit(typed)

		keys := make([]string, 0, len(typed))
		for key := range typed {
			keys = append(keys, key)
		}

		sort.Strings(keys)

		for _, key := range keys {
			walk(typed[key], visit)
		}
	case []any:
		for _, item := range typed {
			walk(item, visit)
		}
	}
}

func merge(base, overlay any) any {
	overlayMap, ok := overlay.(map[string]any)
	if !ok {
		return overlay
	}

	baseMap, ok := base.(map[string]any)
	if !ok {
		return overlay
	}

	for key, value := range overlayMap {
		if existing, found := baseMap[key]; found {
			baseMap[key] = merge(existing, value)

			continue
		}

		baseMap[key] = value
	}

	return baseMap
}
