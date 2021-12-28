package template

import (
	"bytes"
	"embed"
	"fmt"
	"strings"
	"text/template"
)

func Parse(tpl string, data interface{}) ([]byte, error) {
	content, err := templates.ReadFile(tpl)
	if err != nil {
		return nil, fmt.Errorf("cannot read template %s: %w", tpl, err)
	}
	t, err := template.New("").Delims("[[", "]]").Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("cannot parse template %s: %w", tpl, err)
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	if err != nil {
		return nil, fmt.Errorf("cannot execute template %s: %w", tpl, err)
	}
	return buf.Bytes(), nil
}

func Raw(tpl string) ([]byte, error) {
	content, err := templates.ReadFile(tpl)
	if err != nil {
		return nil, fmt.Errorf("cannot read template %s: %w", tpl, err)
	}
	return content, nil
}

func RawAll(root string) (map[string][]byte, error) {
	dirs := make([]string, 0)
	files := make(map[string][]byte)
	err := getDirContents(root, root, dirs, files)
	return files, err
}

func getDirContents(root string, path string, dirs []string, files map[string][]byte) error {
	dirEntries, err := templates.ReadDir(path)
	if err != nil {
		return err
	}
	for _, f := range dirEntries {
		newPath := fmt.Sprintf("%s/%s", path, f.Name())

		if f.IsDir() {
			if !isAlreadyProcessed(newPath, dirs) {
				dirs = append(dirs, newPath)
				err = getDirContents(root, newPath, dirs, files)
				if err != nil {
					return err
				}
			}
		} else {
			content, err := templates.ReadFile(newPath)
			if err != nil {
				return fmt.Errorf("cannot read template %s: %w", newPath, err)
			}
			files[strings.TrimPrefix(newPath, root+"/")] = content
		}
	}
	return nil
}

func isAlreadyProcessed(file string, processedDirs []string) bool {
	for i := 0; i < len(processedDirs); i++ {
		if processedDirs[i] != file {
			continue
		}
		return true
	}
	return false
}

//go:embed templates templates/cicd/* templates/cicd/pipeline_templates/*
var templates embed.FS
