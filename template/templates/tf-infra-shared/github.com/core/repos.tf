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
    var.tf_module_repos,
    var.misc_repos,
  )

  source  = "mineiros-io/repository/github"
  version = "~> 0.11.0"

  name               = each.key
  visibility         = "public" #FIXME
  auto_init          = true
  default_branch     = each.value.default_branch
  archive_on_destroy = false #FIXME?

  allow_rebase_merge     = true
  allow_merge_commit     = !each.value.strict
  delete_branch_on_merge = true

  branch_protections_v3 = [
    {
      branch = each.value.default_branch
      required_status_checks = try(each.value.build_checks != null && length(each.value.build_checks) > 0 ? {
        strict   = each.value.strict
        contexts = each.value.build_checks
      } : tomap(false), {}) # workaround: see https://github.com/hashicorp/terraform/issues/22405#issuecomment-591917758
      #      required_pull_request_reviews = { # FIXME
      #        dismiss_stale_reviews           = true
      #        require_code_owner_reviews      = true
      #        required_approving_review_count = 1
      #      }
    }
  ]
  # restrictions = {} TODO
}
