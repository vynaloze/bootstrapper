resource "github_actions_organization_permissions" "this" {
  allowed_actions      = "local_only"
  enabled_repositories = "selected"

  enabled_repositories_config {
    repository_ids = [for k, v in module.repo_tf_infra: v.repository.repo_id]
  }
}

module "repo_tf_infra" {
  for_each = toset([
    {{- range $_, $value := .Repos }}
    "{{ $value }}",
    {{- end }}
  ])

  source  = "mineiros-io/repository/github"
  version = "~> 0.11.0"

  name           = "tf-infra-shared"
  visibility     = "public" #FIXME
  auto_init      = true
  default_branch = "{{ .DefaultBranch }}"

  allow_rebase_merge     = true
  allow_merge_commit     = {{ not .Strict }}
  delete_branch_on_merge = true

  branch_protections_v3 = [
    {
      branch         = "{{ .DefaultBranch }}"
      required_status_checks = {
        strict   = {{ .Strict }}
        contexts = ["ci/TODO"] #TODO
      }
      required_pull_request_reviews = {
        dismiss_stale_reviews           = true
        require_code_owner_reviews      = true
        required_approving_review_count = 1
      }
    }
  ]
  # restrictions = {} TODO
}
