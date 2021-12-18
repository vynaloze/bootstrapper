provider "github" {
  owner = var.repo_owner
  token = var.repo_password
}

resource "github_actions_organization_permissions" "this" {
  allowed_actions      = "selected"
  enabled_repositories = "selected"

  allowed_actions_config {
    github_owned_allowed = true
    verified_allowed     = true
  }

  enabled_repositories_config {
    repository_ids = [for k, v in module.github_repos : v.repository.repo_id]
  }
}

module "github_repos" {
  for_each = merge(
    var.tf_infra_repos,
    var.misc_repos,
  )

  source  = "mineiros-io/repository/github"
  version = "~> 0.11.0"

  name           = each.key
  visibility     = "public" #FIXME
  auto_init      = true
  default_branch = each.value.default_branch

  allow_rebase_merge     = true
  allow_merge_commit     = !each.value.strict
  delete_branch_on_merge = true

  branch_protections_v3 = [
    {
      branch = each.value.default_branch
      required_status_checks = {
        strict   = each.value.strict
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
