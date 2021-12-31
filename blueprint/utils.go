package blueprint

import (
	"bootstrapper/actor/git"
	"bootstrapper/template"
	"fmt"
)

func templatesToGitFiles(basePath string, templates []Template) ([]git.File, error) {
	gitFiles := make([]git.File, 0)
	for _, file := range templates {
		filename := fmt.Sprintf("templates/%s/%s", basePath, file.SourceFile)
		var pipelineFile []byte
		var err error
		if file.Data == nil {
			pipelineFile, err = template.Raw(filename)
		} else {
			pipelineFile, err = template.Parse(filename, file.Data)
		}
		if err != nil {
			return nil, fmt.Errorf("error fetching template: %w", err)
		}
		gitFiles = append(gitFiles, git.File{Filename: file.TargetFile, Content: string(pipelineFile)})
	}
	return gitFiles, nil
}
