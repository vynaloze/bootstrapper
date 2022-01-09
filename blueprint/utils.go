package blueprint

import (
	"bootstrapper/actor/git"
	"bootstrapper/template"
	"fmt"
)

func templatesToGitFiles(basePath string, templates []Template) ([]git.File, error) {
	gitFiles := make([]git.File, 0)
	for _, file := range templates {

		if file.SourceFile != "" && file.Source != "" {
			return nil, fmt.Errorf("%s: Source and SourceFile conflict with each other", file.TargetFile)
		}

		if file.SourceFile != "" {
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

		} else if file.Source != "" {
			pipelineFile := file.Source
			if file.Data != nil {
				parsed, err := template.ParseContent([]byte(file.Source), file.Data)
				if err != nil {
					return nil, fmt.Errorf("error fetching template: %w", err)
				}
				pipelineFile = string(parsed)
			}
			gitFiles = append(gitFiles, git.File{Filename: file.TargetFile, Content: pipelineFile})

		} else {
			return nil, fmt.Errorf("%s: eitehr Source or SourceFile required", file.TargetFile)
		}
	}
	return gitFiles, nil
}

func commitAndPush(localActor git.LocalActor, branch string, message string, gitFiles []git.File) error {
	err := localActor.CommitMany(branch, message, gitFiles...)
	if err != nil {
		return fmt.Errorf("error committing files: %w", err)
	}
	err = localActor.Push()
	if err != nil {
		return fmt.Errorf("error pushing changes: %w", err)
	}
	return nil
}
