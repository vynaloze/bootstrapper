module "repo_tf_shared_infra" {
  source  = "mineiros-io/repository/github"
  version = "~> 0.11.0"

  name       = "tf-infra-shared"
  visibility = "private"
  auto_init  = false

  allow_rebase_merge     = true
  allow_merge_commit     = {{ not .Strict }}
  delete_branch_on_merge = true

  branch_protections_v3 = [
    {
      branch         = {{ .DefaultBranch }}
      enforce_admins = true
      required_status_checks = {
        strict   = {{ .Strict }}
        contexts = ["ci/TODO"] #TODO
      }
    }
  ]
  required_pull_request_reviews = {
    dismiss_stale_reviews           = true
    require_code_owner_reviews      = true
    required_approving_review_count = 1
  }
  # restrictions = {} TODO
}
