package blueprint

import (
	"bootstrapper/actor/git"
	"bootstrapper/template"
	"fmt"
	"log"
	"time"
)

type SetupCloudEnvModuleOpts struct {
	EnvRepoOpts   git.Opts
	CloudProvider template.CloudProvider
	CICDTemplates []Template
}

func SetupCloudEnvModule(opts *SetupCloudEnvModuleOpts) error {
	log.Printf("setting up cloud env repo")

	localActor := git.NewLocal(&opts.EnvRepoOpts)
	remoteActor, err := git.NewRemote(&opts.EnvRepoOpts)
	if err != nil {
		return fmt.Errorf("cannot initialize remote Git actor: %w", err)
	}

	branch := fmt.Sprintf("%s/%d", opts.EnvRepoOpts.GetAuthorName(), time.Now().UnixMilli())

	log.Printf("preparing CICD pipelines templates")
	ciFiles, err := templatesToGitFiles("cicd/pipeline_templates", opts.CICDTemplates)
	if err != nil {
		return fmt.Errorf("error preparing CICD templates: %w", err)
	}
	err = localActor.CommitMany(branch, "chore: add CI/CD pipelines templates", ciFiles...)
	if err != nil {
		return fmt.Errorf("error committing ci files: %w", err)
	}

	log.Printf("preparing terraform files")
	tfFiles := make([]git.File, 0)
	allFiles, err := template.RawAll("templates/tf-env/" + string(opts.CloudProvider))
	if err != nil {
		return fmt.Errorf("error fetching templates: %w", err)
	}
	for filename, content := range allFiles {
		tfFiles = append(tfFiles, git.File{Filename: filename, Content: string(content)})
	}
	tfFiles = append(tfFiles, git.File{Filename: ".gitignore", Content: template.TerraformGitignore()})
	message := fmt.Sprintf("feat: add initial %s modules", opts.CloudProvider)
	err = localActor.CommitMany(branch, message, tfFiles...)
	if err != nil {
		return fmt.Errorf("error committing tf files: %w", err)
	}

	log.Printf("pushing changes to remote repository")
	err = localActor.Push()
	if err != nil {
		return fmt.Errorf("error pushing changes: %w", err)
	}
	err = remoteActor.RequestReview(&branch, &message)
	if err != nil {
		return fmt.Errorf("error creating PR: %w", err)
	}

	return nil
}
